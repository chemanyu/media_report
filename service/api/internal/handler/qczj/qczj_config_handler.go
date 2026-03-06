package qczj

import (
	"encoding/json"
	"net/http"

	"media_report/service/api/internal/logic/qczj"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// QczjGetConfigHandler 查询 qczj_config 配置
func QczjGetConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := qczj.NewQczjConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// QczjUpdateConfigHandler 修改 qczj_config 配置
func QczjUpdateConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QczjUpdateConfigReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := qczj.NewQczjConfigLogic(r.Context(), svcCtx)
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

// QczjReportListHandler 查询 qczj_report_data 报表数据
func QczjReportListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := qczj.NewQczjConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetReportList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
