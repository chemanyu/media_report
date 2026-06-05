package fz

import (
	"encoding/json"
	"net/http"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// FzGetConfigHandler 查询 fz_config 系数配置
func FzGetConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fz.NewFzConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// FzUpdateConfigHandler 修改 fz_config 系数配置
func FzUpdateConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FzUpdateConfigReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fz.NewFzConfigLogic(r.Context(), svcCtx)
		if err := l.UpdateConfig(&req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "更新失败: " + err.Error(),
			})
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"success": true,
			"message": "更新成功",
		})
	}
}
