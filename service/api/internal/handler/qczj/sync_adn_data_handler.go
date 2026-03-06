package qczj

import (
	"encoding/json"
	"net/http"

	"media_report/service/api/internal/logic/qczj"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SyncAdnDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QczjSyncDataReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if req.ReportDate == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "report_date 不能为空",
			})
			return
		}

		l := qczj.NewSyncAdnDataLogic(r.Context(), svcCtx)
		if err := l.SaveData(&req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "保存数据失败: " + err.Error(),
			})
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"success": true,
			"message": "数据保存成功",
		})
	}
}
