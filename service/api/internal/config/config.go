// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Kuaishou    KuaishouConfig
	OAuthConfig OAuthConfig
	MySQL       MySQLConfig      // 数据库配置
	Schedule    ScheduleConfig   // 定时任务配置
	DingTalk    DingTalkConfig   // 钉钉配置
	FileServer  FileServerConfig // 文件服务器配置
}

type KuaishouConfig struct {
	BaseUrl       string
	Timeout       int
	AdvertiserIds []int64 // 广告主 ID 列表
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
	ReportCron        string // 报表任务 cron 表达式
	TokenRefreshCron  string // token 刷新 cron 表达式
	JuliangReportCron string // 巨量报表任务 cron 表达式
}

type DingTalkConfig struct {
	WebhookURL         string // 钉钉机器人 webhook 地址
	JDReportWebhookURL string // 京东广义巨量数据 webhook 地址
	Enabled            bool   // 是否启用钉钉通知
}

type FileServerConfig struct {
	BaseURL string // 文件服务器基础URL，例如：http://localhost:8888
	Path    string // 文件存储路径，例如：./reports
}
