package script

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

func Cron(config config.Config, db *gorm.DB) {
	// 检查是否配置了定时任务
	if config.Schedule.ReportCron == "" && config.Schedule.TokenRefreshCron == "" {
		logx.Info("未配置定时任务，跳过启动")
		return
	}

	// 创建 cron 调度器
	cronScheduler := cron.New()

	// 添加报表任务
	if config.Schedule.ReportCron != "" {
		_, err := cronScheduler.AddFunc(config.Schedule.ReportCron, func() {
			executeReportJob(db, config.Kuaishou, config.DingTalk)
		})
		if err != nil {
			log.Fatalf("添加报表定时任务失败: %v", err)
		}
		logx.Infof("快手报表定时任务已启动，Cron 表达式: %s", config.Schedule.ReportCron)
	}

	// 添加 token 刷新任务
	if config.Schedule.TokenRefreshCron != "" {
		_, err := cronScheduler.AddFunc(config.Schedule.TokenRefreshCron, func() {
			refreshAccessToken(db, config.Kuaishou, config.OAuthConfig)
		})
		if err != nil {
			log.Fatalf("添加 token 刷新定时任务失败: %v", err)
		}
		logx.Infof("Token 刷新定时任务已启动，Cron 表达式: %s", config.Schedule.TokenRefreshCron)
	}

	// 添加巨量报表任务
	if config.Schedule.JuliangReportCron != "" {
		_, err := cronScheduler.AddFunc(config.Schedule.JuliangReportCron, func() {
			executeJuliangReportJob(db, config.DingTalk, config.FileServer)
		})
		if err != nil {
			log.Fatalf("添加巨量报表定时任务失败: %v", err)
		}
		logx.Infof("巨量报表定时任务已启动，Cron 表达式: %s", config.Schedule.JuliangReportCron)
	}

	// 启动调度器
	cronScheduler.Start()

	// 立即刷新一次 token（可选）
	if config.Schedule.TokenRefreshCron != "" {
		logx.Info("立即刷新一次 token...")
		refreshAccessToken(db, config.Kuaishou, config.OAuthConfig)
	}
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
func executeReportJob(db *gorm.DB, ksConfig config.KuaishouConfig, dingTalk config.DingTalkConfig) {
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

	// 累加所有广告主的数据
	var totalDetail types.KsApiReportDetail
	var hasData bool

	// 循环查询所有广告主
	for _, advertiserId := range ksConfig.AdvertiserIds {
		// 构建请求参数
		req := map[string]interface{}{
			"start_date":           today,
			"end_date":             today,
			"advertiser_id":        advertiserId,
			"temporal_granularity": "DAILY",
		}

		// 调用快手 API
		var resp types.KsApiResponse
		err = client.Post(ctx, "/rest/openapi/v1/report/account_report", req, &resp)
		if err != nil {
			logx.Errorf("调用快手 API 失败 (advertiser_id=%d): %v", advertiserId, err)
			continue
		}

		// 检查响应
		if resp.Code != 0 {
			logx.Errorf("快手 API 返回错误 (advertiser_id=%d): code=%d, message=%s", advertiserId, resp.Code, resp.Message)
			continue
		}

		// 累加数据
		if len(resp.Data.Details) > 0 {
			hasData = true
			detail := resp.Data.Details[0]
			totalDetail.StatDate = detail.StatDate
			totalDetail.Charge += detail.Charge
			totalDetail.AdShow += detail.AdShow
			totalDetail.Bclick += detail.Bclick
			totalDetail.Activation += detail.Activation
			totalDetail.ConversionCost += detail.ConversionCost
			logx.Infof("成功获取广告主 %d 的数据: 消耗=%.2f, 曝光=%.0f, 点击=%d, 激活=%d",
				advertiserId, detail.Charge, detail.AdShow, detail.Bclick, detail.Activation)
		}
	}

	// 打印累加后的数据
	if !hasData {
		logx.Infof("今日暂无数据")
		return
	}

	// 计算平均转化率 和 转化成本
	totalDetail.Charge = totalDetail.Charge * 1.2
	totalDetail.ConversionRatio = float64(totalDetail.Activation) / float64(totalDetail.Bclick)
	totalDetail.ConversionCost = totalDetail.Charge / float64(totalDetail.Activation)

	currentHour := time.Now().Format("15") // 获取当前小时
	conversionRatio := fmt.Sprintf("%.2f%%", totalDetail.ConversionRatio*100)

	// 发送钉钉消息
	sendDingTalkNotification(ctx, dingTalk, totalDetail, conversionRatio, currentHour)

	logx.Infof("报表任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// sendDingTalkNotification 发送钉钉通知
func sendDingTalkNotification(ctx context.Context, dingConfig config.DingTalkConfig, detail types.KsApiReportDetail, conversionRatio, currentHour string) {
	if !dingConfig.Enabled || dingConfig.WebhookURL == "" {
		logx.Info("钉钉通知未启用，跳过发送")
		return
	}

	// 计算数据
	timeWithHour := fmt.Sprintf("%s %s时", detail.StatDate, currentHour)

	// 构建钉钉消息
	markdownText := fmt.Sprintf(
		"#### 美数时报  \n---\n"+
			"**时间**：%s  \n"+
			"**账户**：美数dsp  \n"+
			"**消耗金额**：%.2f  \n"+
			"**注册转化数**：%d  \n"+
			"**转化成本**：%.2f  \n"+
			"**曝光量**：%d  \n"+
			"**点击量**：%d  \n"+
			"**转化率**：%s",
		timeWithHour,
		detail.Charge,
		detail.Activation,
		detail.ConversionCost,
		int64(detail.AdShow),
		detail.Bclick,
		conversionRatio,
	)

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "美数时报",
			"text":  markdownText,
		},
	}

	// 创建 HTTP 客户端发送消息
	client := httpclient.NewClient("", 30)
	var result map[string]interface{}
	err := client.Post(ctx, dingConfig.WebhookURL, msg, &result)
	if err != nil {
		logx.Errorf("发送钉钉消息失败: %v", err)
		return
	}

	logx.Infof("钉钉消息发送成功: %v", result)
}
