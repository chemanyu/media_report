// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Kuaishou    KuaishouConfig
	OAuthConfig OAuthConfig
	JuliangDLS  JuliangConfig    // 巨量DLS配置
	JuliangKH   JuliangConfig    // 巨量KH配置
	ADX         ADXConfig        // ADX 外部接口配置
	MySQL       MySQLConfig      // 数据库配置
	Schedule    ScheduleConfig   // 定时任务配置
	DingTalk    DingTalkConfig   // 钉钉配置
	FileServer  FileServerConfig // 文件服务器配置
	Tanx        TanxConfig       // 淘宝联盟配置
	OppoAPI     OppoAPIConfig    // OPPO广告API配置
	XiaomiAPI   XiaomiAPIConfig  // 小米广告API配置
	//HonorAPI    HonorAPIConfig   // 荣耀广告API配置
	Ulink        UlinkConfig        // 转链模块配置
	SyncFromProd SyncFromProdConfig // 从生产环境同步配置表（仅 34 启用）
	SyncDump     SyncDumpConfig     // 暴露给其他实例的只读 dump 接口（仅生产启用）
}

type KuaishouConfig struct {
	BaseUrl string
	Timeout int
	// AdvertiserIds 已移到数据库表 cainiao_advertiser 管理
}

type OAuthConfig struct {
	AppId  int64  // 应用 ID
	Secret string // 应用密钥
}

type JuliangConfig struct {
	BaseUrl string // API基础地址
	Timeout int    // 请求超时时间（秒）
	AppId   int64  // 应用 ID
	Secret  string // 应用密钥
}

type ADXConfig struct {
	BaseURL string // ADX 接口地址
	APIKey  string // X-API-KEY
	Secret  string // 用于签名的密钥
	Timeout int    // 请求超时时间（秒）
}

type MySQLConfig struct {
	Host            string // 数据库地址
	Port            int    // 端口
	User            string // 用户名
	Password        string // 密码
	Database        string // 数据库名
	Charset         string // 字符集
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	ConnMaxLifetime int    // 连接最大生命周期（秒）
	ConnMaxIdleTime int    // 连接最大空闲时间（秒）
	LogFile         string // SQL日志文件路径
	LogLevel        string // SQL日志级别: silent, error, warn, info
}

type ScheduleConfig struct {
	ReportCron            string // 报表任务 cron 表达式
	TokenRefreshCron      string // token 刷新 cron 表达式
	JuliangReportCron     string // 巨量报表时任务 cron 表达式
	JuliangDayReportCron  string // 巨量报表日报任务 cron 表达式
	HuichuanElmDailyCron  string // 汇川饿了么日报表任务 cron 表达式
	HuichuanElmHourlyCron string // 汇川饿了么小时报表任务 cron 表达式
	TanxCron              string // Tanx 数据抓取任务 cron 表达式
	FzHourCron            string // 飞猪时报更新外投媒体的数据
	//FzDayReportCron       string
	//FzDayCron             string
	QczjHourCron string // QCZJ 分时监控 cron 表达式（每小时触发）
}

type DingTalkConfig struct {
	WebhookURL         string // 钉钉机器人 webhook 地址
	JDReportWebhookURL string // 京东广义巨量数据 webhook 地址
	FzWebhookHourURL   string // 飞猪外投数据 webhook 地址
	FzWebhookDayURL    string
	QczjWebhookURL     string // QCZJ 分时监控钉钉 webhook 地址
	Enabled            bool   // 是否启用钉钉通知
}

type FileServerConfig struct {
	BaseURL string // 文件服务器基础URL，例如：http://localhost:8888
	Path    string // 文件存储路径，例如：./reports
}

// TanxConfig 淘宝联盟配置
type TanxConfig struct {
	AdSlots    []string   // 广告位列表
	Recipients []string   // 邮件接收人列表
	SMTP       SMTPConfig // 邮件服务器配置
}

// SMTPConfig 邮件服务器配置
type SMTPConfig struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 端口
	User     string // SMTP 用户名
	Password string // SMTP 密码
}

// OppoAPIConfig OPPO广告API配置
type OppoAPIConfig struct {
	OwnerId int    // 代理商ID
	ApiId   string // API ID
	ApiKey  string // API Key
}

// XiaomiAPIConfig 小米广告API配置
type XiaomiAPIConfig struct {
	SignId    string // 签名ID
	SecretKey string // 密钥
}

// HonorAPIConfig 荣耀广告API配置
type HonorAPIConfig struct {
	ClientID     string // OAuth2 Client ID
	ClientSecret string // OAuth2 Client Secret
}

// SyncFromProdConfig 从生产环境同步配置表（覆盖式）
// 仅在新部署的从节点启用，主生产保持 Enabled=false
type SyncFromProdConfig struct {
	Enabled       bool     // 是否启用
	BaseURL       string   // 生产入口，例如 https://rta.zhltech.net/guangyixinmedia
	Token         string   // 共享鉴权 Token（X-Sync-Token Header）
	BasicAuthUser string   // nginx HTTP Basic Auth 用户名（可选）
	BasicAuthPass string   // nginx HTTP Basic Auth 密码（可选）
	Cron          string   // 同步 cron 表达式
	Tables        []string // 要同步的表清单
	Timeout       int      // HTTP 超时（秒），默认 30
}

// SyncDumpConfig 只读 dump 接口配置
// 仅在主生产启用，从节点保持 Enabled=false
type SyncDumpConfig struct {
	Enabled bool     // 是否启用
	Token   string   // 共享鉴权 Token，与 SyncFromProd.Token 对齐
	Tables  []string // 允许 dump 的表白名单
}

// UlinkConfig 转链模块配置
type UlinkConfig struct {
	TaobaoAppKey      string // 淘宝客 AppKey
	TaobaoAppSecret   string // 淘宝客 AppSecret
	DefaultMaterialId string // 默认活动素材ID
	DefaultEventId    string // 默认CPA活动ID（5080256=福利购, 4297311=超级红包）
	PythonPath        string // python3 可执行文件路径（如 "python3"）
	ScriptDir         string // Python 脚本目录（如 "./scripts/ulink"）
	ChromeDriverPath  string // ChromeDriver 路径（传给 Python 脚本）
	TempDir           string // 临时文件目录（如 "../uploads"）
}
