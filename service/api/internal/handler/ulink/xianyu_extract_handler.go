package ulink

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	ulinklogic "media_report/service/api/internal/logic/ulink"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// XianyuExtractHandler 闲鱼单链 deeplink 提取
func XianyuExtractHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.XianyuExtractReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := ulinklogic.NewXianyuExtractLogic(r.Context(), svcCtx)
		resp, err := l.ExtractSingle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// XianyuBatchHandler 闲鱼批量 deeplink 提取（文件上传 → Excel 下载）
func XianyuBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ulinklogic.NewXianyuExtractLogic(r.Context(), svcCtx)
		l.ExtractBatch(w, r)
	}
}
