package fz

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// FzSyncAdnDataHandler 同步ADN媒体数据
func FzSyncAdnDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求体
		var req types.FzSyncAdnDataReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 参数验证
		if req.MediaAdvId == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "媒体账户ID不能为空",
			})
			return
		}

		if req.ReportDate == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "报表日期不能为空",
			})
			return
		}

		l := fz.NewFzHourlyReportLogic(r.Context(), svcCtx)

		// 调用logic保存ADN数据
		err := l.SaveAdnData(&req)

		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "保存ADN数据失败: " + err.Error(),
			})
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"success": true,
			"message": "ADN数据保存成功",
		})
	}
}
