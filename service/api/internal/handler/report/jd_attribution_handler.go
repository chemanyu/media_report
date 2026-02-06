package report

import (
	"net/http"

	logic "media_report/service/api/internal/logic/report"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetJDAttributionDataHandler 获取京东归因数据
func GetJDAttributionDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JDAttributionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetJDAttributionDataLogic(r.Context(), svcCtx)
		resp, err := l.GetJDAttributionData(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// ExportJDErrorCountsHandler 导出京东错误统计数据为Excel
func ExportJDErrorCountsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JDErrorExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewExportJDErrorCountsLogic(r.Context(), svcCtx)
		excelData, filename, err := l.ExportJDErrorCounts(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 设置响应头
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Length", string(len(excelData)))

		// 写入Excel数据
		w.Write(excelData)
	}
}
