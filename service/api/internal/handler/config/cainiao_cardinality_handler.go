package config

import (
	"net/http"

	logic "media_report/service/api/internal/logic/config"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CainiaoCardinalityGetHandler 获取菜鸟基数配置
func CainiaoCardinalityGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCainiaoCardinalityGetLogic(r.Context(), svcCtx)
		resp, err := l.Get()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// CainiaoCardinalityUpdateHandler 更新菜鸟基数配置
func CainiaoCardinalityUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CainiaoCardinalityUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCainiaoCardinalityUpdateLogic(r.Context(), svcCtx)
		resp, err := l.Update(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
