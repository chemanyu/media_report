package report

import (
	"net/http"

	logic "media_report/service/api/internal/logic/report"
	"media_report/service/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDhhCookieHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetDhhCookieLogic(r.Context(), svcCtx)
		resp, err := l.GetDhhCookie()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
