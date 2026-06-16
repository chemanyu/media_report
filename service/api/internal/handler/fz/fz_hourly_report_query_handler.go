package fz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/response"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// FzHourlyReportListHandler 获取飞猪小时报列表
func FzHourlyReportListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FzHourlyReportListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fz.NewFzHourlyReportLogic(r.Context(), svcCtx)
		resp, err := l.GetReportList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 使用自定义 OkJsonCtx：强制 Content-Length，避免大响应走 chunked
			// 被前置 nginx 截断（导致浏览器 Network Error）
			response.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
