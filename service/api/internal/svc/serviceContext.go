// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	HTTPClient *httpclient.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建通用 HTTP 客户端
	client := httpclient.NewClient(c.Kuaishou.BaseUrl, c.Kuaishou.Timeout)
	// 设置快手 API 需要的请求头
	client.SetHeader("Access-Token", c.Kuaishou.AccessToken)
	client.SetHeader("Content-Type", "application/json")

	return &ServiceContext{
		Config:     c,
		HTTPClient: client,
	}
}
