package script

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

// executeJuliangReportJob 执行巨量报表任务
func executeJuliangReportJob(db *gorm.DB, dingTalk config.DingTalkConfig) {
	ctx := context.Background()
	logx.Infof("开始执行巨量报表任务 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 从数据库获取 cookie 和 csrf token
	mediaToken, err := model.GetByMedia(db, "juliang_pachong")
	if err != nil {
		logx.Errorf("从数据库获取巨量token失败: %v", err)
		return
	}

	cookie := mediaToken.Token
	csrfToken := mediaToken.RefreshToken

	if cookie == "" || csrfToken == "" {
		logx.Error("巨量 Cookie 或 CSRF Token 为空，无法执行任务")
		return
	}

	// 计算时间范围（今天的开始和结束时间戳）
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	startTime := strconv.FormatInt(startOfDay.Unix(), 10)
	endTime := strconv.FormatInt(endOfDay.Unix(), 10)

	// 创建 HTTP 客户端
	client := httpclient.NewClient("https://business.oceanengine.com", 30)
	client.SetHeader("Cookie", cookie)
	client.SetHeader("x-csrftoken", csrfToken)
	client.SetHeader("Content-Type", "application/json")

	// 累加统计数据
	var totalCost float64     // 总消耗
	var totalShowCnt int64    // 总曝光
	var totalClickCnt int64   // 总点击
	var totalConvertCnt int64 // 总转化
	var totalAccounts int     // 总账户数

	// 分页查询
	page := 1
	limit := 100
	hasMore := true

	for hasMore {
		// 构建请求参数
		req := map[string]interface{}{
			"start_time":   startTime,
			"end_time":     endTime,
			"offset":       page,
			"limit":        limit,
			"order_type":   1,
			"account_type": 0,
			"cascade_metrics": []string{
				"advertiser_name",
				"advertiser_id",
				"advertiser_status",
				"advertiser_remark",
				"advertiser_agent_name",
				"advertiser_agent_id",
				"advertiser_followed",
			},
			"fields": []string{
				"stat_cost",
				"stat_cash_cost",
				"show_cnt",
				"click_cnt",
				"ctr",
				"convert_cnt",
				"conversion_cost",
				"conversion_rate",
			},
			"filter": map[string]interface{}{
				"advertiser":      map[string]interface{}{},
				"group":           map[string]interface{}{},
				"pricingCategory": []int{2},
				"campaign":        map[string]interface{}{},
				"is_active":       true,
			},
			"ocean_white":      true,
			"order_field":      "stat_cost",
			"platform_version": "2.0",
		}

		// 调用巨量 API
		var resp types.JuliangApiResponse
		err = client.Post(ctx, "/nbs/api/bm/promotion/ad/get_account_list", req, &resp)
		if err != nil {
			logx.Errorf("调用巨量 API 失败 (page=%d): %v", page, err)
			break
		}

		// 检查响应
		if resp.Code != 0 {
			logx.Errorf("巨量 API 返回错误 (page=%d): code=%d, message=%s", page, resp.Code, resp.Msg)
			break
		}

		// 累加数据
		for _, account := range resp.Data.DataList {
			totalAccounts++

			// 解析消耗（去除逗号）
			cost := parseNumber(account.StatCost)
			totalCost += cost

			// 解析曝光数
			showCnt := parseInt64(account.ShowCnt)
			totalShowCnt += showCnt

			// 解析点击数
			clickCnt := parseInt64(account.ClickCnt)
			totalClickCnt += clickCnt

			// 解析转化数
			convertCnt := parseInt64(account.ConvertCnt)
			totalConvertCnt += convertCnt

			logx.Infof("账户 %s (%d): 消耗=%.2f, 曝光=%d, 点击=%d, 转化=%d",
				account.AdvertiserName, account.AdvertiserId, cost, showCnt, clickCnt, convertCnt)
		}

		// 检查是否还有更多数据
		hasMore = resp.Data.Pagination.HasMore
		page++

		logx.Infof("已处理第 %d 页，本页账户数: %d，累计账户数: %d",
			page-1, len(resp.Data.DataList), totalAccounts)

		// 避免请求过快
		if hasMore {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// 打印汇总数据
	if totalAccounts == 0 {
		logx.Info("今日暂无巨量账户数据")
		return
	}

	// 计算转化成本和转化率
	var avgConversionCost float64
	var avgConversionRate float64

	if totalConvertCnt > 0 {
		avgConversionCost = totalCost / float64(totalConvertCnt)
	}
	if totalClickCnt > 0 {
		avgConversionRate = float64(totalConvertCnt) / float64(totalClickCnt) * 100
	}

	logx.Infof("巨量报表汇总 - 账户数: %d, 总消耗: %.2f, 总曝光: %d, 总点击: %d, 总转化: %d, 转化成本: %.2f, 转化率: %.2f%%",
		totalAccounts, totalCost, totalShowCnt, totalClickCnt, totalConvertCnt, avgConversionCost, avgConversionRate)

	// 发送钉钉通知
	sendJuliangDingTalkNotification(ctx, dingTalk, totalCost, totalShowCnt, totalClickCnt, totalConvertCnt,
		avgConversionCost, avgConversionRate, totalAccounts)

	logx.Infof("巨量报表任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// sendJuliangDingTalkNotification 发送巨量钉钉通知
func sendJuliangDingTalkNotification(ctx context.Context, dingConfig config.DingTalkConfig,
	totalCost float64, totalShowCnt, totalClickCnt, totalConvertCnt int64,
	avgConversionCost, avgConversionRate float64, totalAccounts int) {

	if !dingConfig.Enabled || dingConfig.WebhookURL == "" {
		logx.Info("钉钉通知未启用，跳过发送")
		return
	}

	// 获取当前时间
	now := time.Now()
	timeStr := now.Format("2006-01-02 15时")

	// 构建钉钉消息
	markdownText := fmt.Sprintf(
		"#### 巨量时报  \n---\n"+
			"**时间**：%s  \n"+
			"**账户数**：%d  \n"+
			"**消耗金额**：%.2f  \n"+
			"**注册转化数**：%d  \n"+
			"**转化成本**：%.2f  \n"+
			"**曝光量**：%d  \n"+
			"**点击量**：%d  \n"+
			"**转化率**：%.2f%%",
		timeStr,
		totalAccounts,
		totalCost,
		totalConvertCnt,
		avgConversionCost,
		totalShowCnt,
		totalClickCnt,
		avgConversionRate,
	)

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "巨量时报",
			"text":  markdownText,
		},
	}

	// 创建 HTTP 客户端发送消息
	client := httpclient.NewClient("", 30)
	var result map[string]interface{}
	err := client.Post(ctx, dingConfig.WebhookURL, msg, &result)
	if err != nil {
		logx.Errorf("发送巨量钉钉消息失败: %v", err)
		return
	}

	logx.Infof("巨量钉钉消息发送成功: %v", result)
}

// parseNumber 解析带逗号的数字字符串为 float64
func parseNumber(s string) float64 {
	// 去除逗号
	s = strings.ReplaceAll(s, ",", "")
	// 去除百分号（如果有）
	s = strings.TrimSuffix(s, "%")

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

// parseInt64 解析带逗号的数字字符串为 int64
func parseInt64(s string) int64 {
	// 去除逗号
	s = strings.ReplaceAll(s, ",", "")

	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
