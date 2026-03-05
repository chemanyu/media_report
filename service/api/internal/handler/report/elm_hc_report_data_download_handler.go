package report

import (
	"fmt"
	"net/http"

	logic "media_report/service/api/internal/logic/report"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DownloadElmHcReportDataHandler 下载饿了么汇川报表数据（Excel）
func DownloadElmHcReportDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ElmHcReportDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDownloadElmHcReportDataLogic(r.Context(), svcCtx)
		excelData, filename, err := l.DownloadElmHcReportData(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Length", fmt.Sprint(len(excelData)))
		w.Write(excelData)
	}
}
