package config

import (
	"fmt"
	"net/http"
	"time"

	"media_report/service/api/internal/logic/config"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取服务费配置列表
func GetServiceFeeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = ip
		}
		logx.Infof("[ServiceFeeList] 请求开始: client=%s, url=%s, ua=%s", clientIP, r.URL.String(), r.UserAgent())
		defer func() {
			if err := recover(); err != nil {
				logx.Errorf("[ServiceFeeList] panic: %v", err)
				httpx.ErrorCtx(r.Context(), w, fmt.Errorf("服务器内部错误: %v", err))
			}
			logx.Infof("[ServiceFeeList] 请求结束: client=%s, url=%s, 耗时=%v", clientIP, r.URL.String(), time.Since(start))
		}()

		l := config.NewGetServiceFeeListLogic(r.Context(), svcCtx)
		resp, err := l.GetServiceFeeList()
		if err != nil {
			logx.Errorf("[ServiceFeeList] 业务错误: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// 创建服务费配置
func CreateServiceFeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateServiceFeeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewCreateServiceFeeLogic(r.Context(), svcCtx)
		resp, err := l.CreateServiceFee(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// 更新服务费配置
func UpdateServiceFeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateServiceFeeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewUpdateServiceFeeLogic(r.Context(), svcCtx)
		resp, err := l.UpdateServiceFee(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// 删除服务费配置
func DeleteServiceFeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteServiceFeeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewDeleteServiceFeeLogic(r.Context(), svcCtx)
		resp, err := l.DeleteServiceFee(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
