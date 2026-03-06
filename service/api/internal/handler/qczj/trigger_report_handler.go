package qczj

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"media_report/service/api/internal/logic/qczj"
	"media_report/service/api/internal/svc"
)

func TriggerQczjReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := qczj.NewTriggerQczjReportLogic(r.Context(), svcCtx)
		resp, err := l.TriggerReport()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
