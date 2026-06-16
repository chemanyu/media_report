package response

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// OkJsonCtx 与 httpx.OkJsonCtx 行为一致，但**显式设置 Content-Length**，
// 从而避免 net/http 在 body 超过内部缓冲(~2KB)时自动切换到
// Transfer-Encoding: chunked。
//
// 背景：线上前置 nginx 反代对 chunked 响应转发存在问题，导致较大的
// JSON 响应(如"查全部")在浏览器侧报 axios "Network Error"，而较小的
// 响应(分媒体/空结果)因带 Content-Length 可正常返回。强制 Content-Length
// 即可让所有响应走定长传输，绕开前置层对 chunked 的兼容问题。
func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	bs, err := jsonx.Marshal(v)
	if err != nil {
		// 与 go-zero 保持一致：marshal 失败返回 500
		logc.Errorf(ctx, "marshal json failed, error: %v", err)
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	// 显式 Content-Length，禁止 chunked
	header.Set("Content-Length", strconv.Itoa(len(bs)))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bs); err != nil {
		logc.Errorf(ctx, "write response failed, error: %v", err)
	}
}
