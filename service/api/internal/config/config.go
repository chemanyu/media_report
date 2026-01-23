// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Kuaishou    KuaishouConfig
	OAuthConfig OAuthConfig
	MySQL       MySQLConfig    // 数据库配置
	Schedule    ScheduleConfig // 定时任务配置
}

type KuaishouConfig struct {
	BaseUrl      string
	Timeout      int
	AdvertiserId int64 // 广告主 ID
}

type OAuthConfig struct {
	AppId  int64  // 应用 ID
	Secret string // 应用密钥
}

type MySQLConfig struct {
	Host     string // 数据库地址
	Port     int    // 端口
	User     string // 用户名
	Password string // 密码
	Database string // 数据库名
	Charset  string // 字符集
}

type ScheduleConfig struct {
	ReportCron       string // 报表任务 cron 表达式
	TokenRefreshCron string // token 刷新 cron 表达式
}
