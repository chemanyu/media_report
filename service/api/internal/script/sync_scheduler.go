package script

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/service/api/internal/config"
	"media_report/service/api/internal/syncrule"
)

const syncTokenHeader = "X-Sync-Token"

type syncDumpResp struct {
	Code  int              `json:"code"`
	Table string           `json:"table"`
	Count int              `json:"count"`
	Rows  []map[string]any `json:"rows"`
}

// RegisterSyncFromProd 在 cron 调度器中注册"从生产同步配置表"任务。
// 调用方在 Cron(...) 里把已有的 *cron.Cron 传进来即可。
func RegisterSyncFromProd(scheduler *cron.Cron, db *gorm.DB, cfg config.SyncFromProdConfig) {
	if !cfg.Enabled {
		logx.Info("SyncFromProd 未启用，跳过注册")
		return
	}
	if cfg.BaseURL == "" || cfg.Token == "" || cfg.Cron == "" || len(cfg.Tables) == 0 {
		logx.Error("SyncFromProd 启用但配置不完整，跳过注册")
		return
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	if _, err := scheduler.AddFunc(cfg.Cron, func() {
		runSyncOnce(db, cfg, timeout)
	}); err != nil {
		logx.Errorf("添加 SyncFromProd 定时任务失败: %v", err)
		return
	}
	logx.Infof("SyncFromProd 定时任务已启动，Cron: %s, 表: %v", cfg.Cron, cfg.Tables)

	// 启动后立即执行一次（避免等到下一个 cron 节点）
	go func() {
		time.Sleep(5 * time.Second)
		runSyncOnce(db, cfg, timeout)
	}()
}

func runSyncOnce(db *gorm.DB, cfg config.SyncFromProdConfig, timeoutSec int) {
	logx.Infof("[sync] 开始一轮同步, 共 %d 张表, BaseURL=%s, TokenLen=%d, BasicAuthUser=%q, Timeout=%ds",
		len(cfg.Tables), cfg.BaseURL, len(cfg.Token), cfg.BasicAuthUser, timeoutSec)
	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}

	successCount := 0
	for _, tbl := range cfg.Tables {
		if err := syncOneTable(client, db, cfg, tbl); err != nil {
			logx.Errorf("[sync] %s 同步失败: %v", tbl, err)
			continue
		}
		successCount++
	}
	logx.Infof("[sync] 本轮同步完成, 成功 %d / %d", successCount, len(cfg.Tables))
}

func syncOneTable(client *http.Client, db *gorm.DB, cfg config.SyncFromProdConfig, table string) error {
	logx.Infof("[sync] %s 开始抓取上游数据", table)
	rows, err := fetchDump(client, cfg, table)
	if err != nil {
		return err
	}
	logx.Infof("[sync] %s 上游返回 %d 行, 准备写入本地", table, len(rows))

	// 有部分同步规则的表（如 fz_hourly_report 只同步 media=huawei 的最近 30 天）：
	// 只删除+重灌规则范围内的行，本地自采的其他媒体、以及窗口外的历史数据都保持不动。
	if f, ok := syncrule.For(table); ok {
		return syncPartial(db, table, f, rows)
	}

	// 覆盖式：事务内 TRUNCATE + 批量 INSERT
	return db.Transaction(func(tx *gorm.DB) error {
		// 反引号包表名，与 dump 接口白名单形成双重防线
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)).Error; err != nil {
			return fmt.Errorf("truncate: %w", err)
		}
		if len(rows) == 0 {
			logx.Infof("[sync] %s 上游 0 行，已清空本地", table)
			return nil
		}
		if err := tx.Table(table).CreateInBatches(rows, 500).Error; err != nil {
			return fmt.Errorf("insert: %w", err)
		}
		logx.Infof("[sync] %s 已覆盖 %d 行", table, len(rows))
		return nil
	})
}

// syncPartial 按规则做"部分覆盖"：只删除+重灌规则范围内的行。
// 范围外的数据（其他媒体、日期窗口之外的历史）在从节点原地留存。
func syncPartial(db *gorm.DB, table string, f syncrule.Filter, rows []map[string]any) error {
	// now 只取一次：DELETE 的日期窗口和行过滤必须用同一个基准，
	// 否则跨零点时两者窗口错位，可能删掉刚写进去的行。
	now := time.Now()
	desc := f.Desc(now)

	kept := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if !f.Match(row, now) {
			continue
		}
		// 剔除自增主键等列，避免与本地自采数据的主键冲突
		for _, col := range f.DropColumns {
			delete(row, col)
		}
		kept = append(kept, row)
	}
	if skipped := len(rows) - len(kept); skipped > 0 {
		logx.Errorf("[sync] %s 上游返回了 %d 行不在同步范围(%s)的数据，已丢弃（两端规则可能不一致）",
			table, skipped, desc)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		q := tx.Table(table)
		for _, c := range f.Conds(now) {
			q = q.Where(c.SQL, c.Args...)
		}
		res := q.Delete(nil)
		if res.Error != nil {
			return fmt.Errorf("delete (%s): %w", desc, res.Error)
		}
		if len(kept) == 0 {
			logx.Infof("[sync] %s 上游范围内(%s) 0 行，已清空本地对应数据 %d 行", table, desc, res.RowsAffected)
			return nil
		}
		if err := tx.Table(table).CreateInBatches(kept, 500).Error; err != nil {
			return fmt.Errorf("insert: %w", err)
		}
		logx.Infof("[sync] %s 部分覆盖完成(%s)：删除 %d 行，写入 %d 行",
			table, desc, res.RowsAffected, len(kept))
		return nil
	})
}

func fetchDump(client *http.Client, cfg config.SyncFromProdConfig, table string) ([]map[string]any, error) {
	endpoint := strings.TrimRight(cfg.BaseURL, "/") + "/api/internal/sync/dump?table=" + url.QueryEscape(table)
	logx.Infof("[sync] %s GET %s", table, endpoint)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set(syncTokenHeader, cfg.Token)
	if cfg.BasicAuthUser != "" {
		req.SetBasicAuth(cfg.BasicAuthUser, cfg.BasicAuthPass)
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do (%s): %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	cost := time.Since(start)
	logx.Infof("[sync] %s HTTP %d, bodyLen=%d, cost=%s", table, resp.StatusCode, len(body), cost)
	if readErr != nil {
		return nil, fmt.Errorf("read body: %w", readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, truncate(body, 500))
	}

	var parsed syncDumpResp
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode: %w (body: %s)", err, truncate(body, 500))
	}
	if parsed.Code != 200 {
		return nil, fmt.Errorf("upstream code %d (body: %s)", parsed.Code, truncate(body, 500))
	}
	return parsed.Rows, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
