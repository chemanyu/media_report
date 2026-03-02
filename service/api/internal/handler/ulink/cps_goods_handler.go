package ulink

import (
	"net/http"

	ulinklogic "media_report/service/api/internal/logic/ulink"
	"media_report/service/api/internal/svc"
)

// CpsGoodsBatchHandler CPS高佣商品库批量导出（选择商品库 → Excel 下载）
func CpsGoodsBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ulinklogic.NewCpsGoodsLogic(r.Context(), svcCtx)
		l.GetCpsGoods(w, r)
	}
}
