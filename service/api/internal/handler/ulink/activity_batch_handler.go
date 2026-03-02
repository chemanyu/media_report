package ulink

import (
	"net/http"

	ulinklogic "media_report/service/api/internal/logic/ulink"
	"media_report/service/api/internal/svc"
)

// ActivityBatchHandler 淘宝客活动短链批量获取（文件上传 → Excel 下载）
func ActivityBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ulinklogic.NewActivityBatchLogic(r.Context(), svcCtx)
		l.GetActivityBatch(w, r)
	}
}
