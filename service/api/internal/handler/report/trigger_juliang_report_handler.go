package report

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	logic "media_report/service/api/internal/logic/report"
	"media_report/service/api/internal/svc"
)

func TriggerJuliangReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewTriggerJuliangReportLogic(r.Context(), svcCtx)
		resp, err := l.TriggerJuliangReport()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
