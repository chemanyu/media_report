package zfb

import (
	"net/http"

	"media_report/service/api/internal/logic/zfb"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ZfbDownloadHandler ZFB独立下载接口处理器
// 该接口通过独立的 /zfb 反向代理访问，无需密码验证
// 与 guangyixinmedia 路径完全隔离
func ZfbDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ZfbDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := zfb.NewZfbDownloadLogic(r.Context(), svcCtx)
		filePath, filename, err := l.ZfbDownload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 设置响应头，让浏览器下载文件
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.Header().Set("Content-Transfer-Encoding", "binary")

		// 直接发送文件
		http.ServeFile(w, r, filePath)
	}
}
