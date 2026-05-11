package sync

import (
	"net/http"

	"media_report/service/api/internal/svc"

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

		httpx.OkJsonCtx(r.Context(), w, map[string]any{
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

// dumpTable 用反引号包表名（白名单已限制 SQL 注入面），SELECT * 全表
func dumpTable(db *gorm.DB, table string) ([]map[string]any, error) {
	var rows []map[string]any
	if err := db.Table(table).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
