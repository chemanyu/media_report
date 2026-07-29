package sync

import (
	"net/http"
	"time"

	"media_report/service/api/internal/response"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/syncrule"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"
)

const tokenHeader = "X-Sync-Token"

// DumpHandler 返回指定白名单表的全部行（JSON 数组）。
// GET /api/internal/sync/dump?table=media_token
// Header: X-Sync-Token: <共享密钥>
func DumpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := svcCtx.Config.SyncDump

		if !cfg.Enabled {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]any{
				"code":    404,
				"message": "sync dump disabled",
			})
			return
		}

		if cfg.Token == "" || r.Header.Get(tokenHeader) != cfg.Token {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "invalid sync token",
			})
			return
		}

		table := r.URL.Query().Get("table")
		if table == "" || !inWhitelist(table, cfg.Tables) {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]any{
				"code":    400,
				"message": "table not allowed: " + table,
			})
			return
		}

		rows, err := dumpTable(svcCtx.DB, table)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]any{
				"code":    500,
				"message": "dump failed: " + err.Error(),
			})
			return
		}

		// 列表类大响应：走显式 Content-Length，避免前置 nginx 截断 chunked
		response.OkJsonCtx(r.Context(), w, map[string]any{
			"code":  200,
			"table": table,
			"count": len(rows),
			"rows":  rows,
		})
	}
}

func inWhitelist(name string, list []string) bool {
	for _, t := range list {
		if t == name {
			return true
		}
	}
	return false
}

// dumpTable 用反引号包表名（白名单已限制 SQL 注入面），SELECT * 全表。
// 若该表在 syncrule 里配了部分同步规则（如 fz_hourly_report 只同步 media=huawei
// 且只取最近 30 天），则只导出规则范围内的行。
func dumpTable(db *gorm.DB, table string) ([]map[string]any, error) {
	var rows []map[string]any
	q := db.Table(table)
	if f, ok := syncrule.For(table); ok {
		for _, c := range f.Conds(time.Now()) {
			q = q.Where(c.SQL, c.Args...)
		}
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
