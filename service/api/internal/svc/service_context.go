// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	HTTPClient *httpclient.Client
	DB         *gorm.DB // 数据库连接
}

func NewServiceContext(c config.Config, db *gorm.DB) *ServiceContext {
	// 创建通用 HTTP 客户端
	client := httpclient.NewClient(c.Kuaishou.BaseUrl, c.Kuaishou.Timeout)

	return &ServiceContext{
		Config:     c,
		HTTPClient: client,
		DB:         db,
	}
}
