package fz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// FzAdvertiserListHandler 获取飞猪媒体账户列表
func FzAdvertiserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fz.NewFzAdvertiserLogic(r.Context(), svcCtx)
		resp, err := l.GetList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// FzAdvertiserAddHandler 添加飞猪媒体账户
func FzAdvertiserAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FzAdvertiserAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fz.NewFzAdvertiserLogic(r.Context(), svcCtx)
		resp, err := l.Add(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// FzAdvertiserUpdateHandler 更新飞猪媒体账户
func FzAdvertiserUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FzAdvertiserUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fz.NewFzAdvertiserLogic(r.Context(), svcCtx)
		resp, err := l.Update(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// FzAdvertiserDeleteHandler 删除飞猪媒体账户
func FzAdvertiserDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FzAdvertiserDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fz.NewFzAdvertiserLogic(r.Context(), svcCtx)
		resp, err := l.Delete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
