package script

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

var configFile = flag.String("f", "etc/scheduler.yaml", "the config file")

// ScheduleConfig 定时任务配置
type ScheduleConfig struct {
	ReportCron       string // 报表任务 cron 表达式
	TokenRefreshCron string // token 刷新 cron 表达式
}

// Config 配置结构
type CronConfig struct {
	service.ServiceConf
	Schedule ScheduleConfig
}

func Cron(config config.Config, db *gorm.DB) {
	flag.Parse()

	// 加载配置
	var c CronConfig
	conf.MustLoad(*configFile, &c)

	// 设置日志
	logx.MustSetup(c.Log)
	defer logx.Close()

	// 创建 cron 调度器
	cronScheduler := cron.New()

	// 添加报表任务
	_, err := cronScheduler.AddFunc(c.Schedule.ReportCron, func() {
		executeReportJob(db, config.Kuaishou)
	})
	if err != nil {
		log.Fatalf("添加报表定时任务失败: %v", err)
	}

	// 添加 token 刷新任务
	_, err = cronScheduler.AddFunc(c.Schedule.TokenRefreshCron, func() {
		refreshAccessToken(db, config.Kuaishou, config.OAuthConfig)
	})
	if err != nil {
		log.Fatalf("添加 token 刷新定时任务失败: %v", err)
	}

	// 启动调度器
	cronScheduler.Start()
	defer cronScheduler.Stop()

	logx.Infof("快手报表定时任务已启动，Cron 表达式: %s", c.Schedule.ReportCron)
	logx.Infof("Token 刷新定时任务已启动，Cron 表达式: %s", c.Schedule.TokenRefreshCron)
	fmt.Println("按 Ctrl+C 退出...")

	// 立即刷新一次 token（可选）
	fmt.Println("立即刷新一次 token...")
	refreshAccessToken(db, config.Kuaishou, config.OAuthConfig)

	// 立即执行一次报表任务（可选）
	fmt.Println("立即执行一次报表任务...")
	executeReportJob(db, config.Kuaishou)

	// 保持程序运行
	select {}
}

// refreshAccessToken 刷新 access token
func refreshAccessToken(db *gorm.DB, ksConfig config.KuaishouConfig, oauthConfig config.OAuthConfig) {
	ctx := context.Background()
	logx.Infof("开始刷新 access token - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 从数据库获取当前的 refresh_token
	mediaToken, err := model.GetByMedia(db, "kuaishou")
	if err != nil {
		logx.Errorf("从数据库获取 token 失败: %v", err)
		return
	}

	// 创建 HTTP 客户端
	client := httpclient.NewClient(ksConfig.BaseUrl, ksConfig.Timeout)
	client.SetHeader("Content-Type", "application/json")

	// 构建刷新请求
	req := map[string]interface{}{
		"app_id":        oauthConfig.AppId,
		"secret":        oauthConfig.Secret,
		"refresh_token": mediaToken.RefreshToken,
	}

	// 调用刷新 token API
	var resp types.TokenRefreshResponse
	err = client.Post(ctx, "/rest/openapi/oauth2/authorize/refresh_token", req, &resp)
	if err != nil {
		logx.Errorf("调用刷新 token API 失败: %v", err)
		return
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("刷新 token 失败: code=%d, message=%s", resp.Code, resp.Message)
		return
	}

	// 更新数据库中的 token
	mediaToken.Token = resp.Data.AccessToken
	mediaToken.RefreshToken = resp.Data.RefreshToken
	err = db.Save(mediaToken).Error
	if err != nil {
		logx.Errorf("更新数据库 token 失败: %v", err)
		return
	}

	logx.Infof("Token 刷新成功，新 AccessToken: %s, 有效期: %d 秒", resp.Data.AccessToken, resp.Data.AccessTokenExpiresIn)
	logx.Infof("新 RefreshToken: %s, 有效期: %d 秒", resp.Data.RefreshToken, resp.Data.RefreshTokenExpiresIn)
}

// executeReportJob 执行报表任务
func executeReportJob(db *gorm.DB, ksConfig config.KuaishouConfig) {
	ctx := context.Background()
	logx.Infof("开始执行快手报表任务 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 从数据库获取当前有效的 access token
	accessToken, _, err := model.GetTokensByMedia(db, "kuaishou")
	if err != nil {
		logx.Errorf("从数据库获取 token 失败: %v", err)
		return
	}

	// 创建 HTTP 客户端
	client := httpclient.NewClient(ksConfig.BaseUrl, ksConfig.Timeout)
	client.SetHeader("Access-Token", accessToken)
	client.SetHeader("Content-Type", "application/json")

	// 获取当前日期
	today := time.Now().Format("2006-01-02")

	// 构建请求参数
	req := map[string]interface{}{
		"start_date":           today,
		"end_date":             today,
		"advertiser_id":        ksConfig.AdvertiserId,
		"temporal_granularity": "DAILY",
	}

	// 调用快手 API
	var resp types.KsApiResponse
	err = client.Post(ctx, "/rest/openapi/v1/report/account_report", req, &resp)
	if err != nil {
		logx.Errorf("调用快手 API 失败: %v", err)
		return
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("快手 API 返回错误: code=%d, message=%s", resp.Code, resp.Message)
		return
	}

	// 打印数据
	if len(resp.Data.Details) == 0 {
		logx.Infof("今日暂无数据")
		return
	}

	logx.Infof("成功获取 %d 条数据", len(resp.Data.Details))
	for i, detail := range resp.Data.Details {
		charge := detail.Charge * 1.5
		conversionCost := detail.ConversionCost * 1.5
		conversionRatio := fmt.Sprintf("%.2f%%", detail.ConversionRatio*100)

		logx.Infof("数据 %d: 时间=%s, 账户=美致dsp, 消耗=%.2f, 曝光=%d, 点击=%d, 注册转化数=%d, 转化成本=%.2f, 转化率=%s",
			i+1, detail.StatDate, charge, int64(detail.AdShow), detail.Bclick,
			detail.Activation, conversionCost, conversionRatio)
	}

	logx.Infof("报表任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}
