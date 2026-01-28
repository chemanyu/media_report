package download

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"media_report/service/api/internal/logic/download"
	"media_report/service/api/internal/svc"
)

func DownloadReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := download.NewDownloadReportLogic(r.Context(), svcCtx)
		err := l.DownloadReport(w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
