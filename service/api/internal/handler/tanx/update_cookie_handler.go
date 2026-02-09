package tanx

import (
	"net/http"

	"media_report/service/api/internal/logic/tanx"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UpdateCookieHandler 更新Cookie
func UpdateCookieHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TanxUpdateCookieReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tanx.NewUpdateCookieLogic(r.Context(), svcCtx)
		resp, err := l.UpdateCookie(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
