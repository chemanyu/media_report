package fz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/svc"
)

// FzSyncOppoDataHandler 同步飞猪OPPO数据
func FzSyncOppoDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取日期参数（可选）
		reportDate := r.URL.Query().Get("date") // 格式: 20260211

		l := fz.NewFzHourlyReportLogic(r.Context(), svcCtx)

		var count int
		var err error

		if reportDate != "" {
			// 同步指定日期的数据
			count, err = l.SyncOppoData(reportDate)
		} else {
			// 默认同步今天的数据
			count, err = l.SyncTodayOppoData()
		}

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"success": true,
			"message": "同步成功",
			"count":   count,
		})
	}
}
