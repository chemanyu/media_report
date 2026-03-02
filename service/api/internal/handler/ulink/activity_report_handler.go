package ulink

import (
	"net/http"

	ulinklogic "media_report/service/api/internal/logic/ulink"
	"media_report/service/api/internal/svc"
)

// ActivityReportHandler 淘宝客活动报表批量查询（Excel 上传 → Excel 下载）
func ActivityReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ulinklogic.NewActivityReportLogic(r.Context(), svcCtx)
		l.GetActivityReport(w, r)
	}
}
