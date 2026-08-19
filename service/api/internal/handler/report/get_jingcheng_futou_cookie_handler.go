package report

import (
	"net/http"

	logic "media_report/service/api/internal/logic/report"
	"media_report/service/api/internal/response"
	"media_report/service/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetJingchengFutouCookieHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetJingchengFutouCookieLogic(r.Context(), svcCtx)
		resp, err := l.GetJingchengFutouCookie()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// Cookie 串通常远超 2KB，必须显式带 Content-Length，否则线上 nginx 会截断
			response.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
