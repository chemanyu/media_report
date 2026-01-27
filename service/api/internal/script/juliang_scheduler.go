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

// AccountReportData 账户报表数据
type AccountReportData struct {
	AdvertiserId    int64   // 账户ID
	AdvertiserName  string  // 账户名称
	Subject         string  // 主体
	Port            string  // 端口
	ServiceProvider string  // 服务商
	TaskCode        string  // 任务代码
	Cost            float64 // 消耗
	CashCost        float64 // 现金消耗
	RebateCost      float64 // 返点消耗
	ShowCnt         int64   // 曝光
	ClickCnt        int64   // 点击
	Ctr             string  // 点击率
	ConvertCnt      int64   // 转化
	ConversionCost  string  // 转化成本
	ConversionRate  string  // 转化率
	ServiceFeeCost  float64 // 服务费成本
	Revenue         float64 // 预估收入
	Profit          float64 // 预估利润
	ProfitRate      float64 // 预估利润率
}

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

	// 预加载配置数据
	rebateMap, err := model.LoadRebateConfigMap(db)
	if err != nil {
		logx.Errorf("加载返点配置失败: %v", err)
		return
	}
	serviceFeeMap, err := model.LoadServiceFeeConfigMap(db)
	if err != nil {
		logx.Errorf("加载服务费配置失败: %v", err)
		return
	}
	taskTypeMap, err := model.LoadTaskTypeConfigMap(db)
	if err != nil {
		logx.Errorf("加载任务类型配置失败: %v", err)
		return
	}

	// 累加统计数据
	var totalCost float64           // 总消耗
	var totalCashCost float64       // 总现金消耗
	var totalRebateCost float64     // 总返后消耗
	var totalShowCnt int64          // 总曝光
	var totalClickCnt int64         // 总点击
	var totalConvertCnt int64       // 总转化
	var totalConversionCost float64 // 总转化成本
	var totalAccounts int           // 总账户数
	var totalServiceFeeCost float64 // 总服务商成本
	var totalRevenue float64        // 总预估收入
	var totalProfit float64         // 总预估利润
	var skippedAccounts int         // 跳过的账户数

	// 保存每条账户数据
	var accountReports []AccountReportData

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
			// 解析备注字段：主体-端口-服务商-任务
			remark := strings.TrimSpace(account.AdvertiserRemark)
			parts := strings.Split(remark, "-")

			// 如果分割后小于4个部分，跳过
			if len(parts) < 4 {
				skippedAccounts++
				continue
			}

			subject := strings.TrimSpace(parts[0])         // 主体
			port := strings.TrimSpace(parts[1])            // 端口
			serviceProvider := strings.TrimSpace(parts[2]) // 服务商
			taskCode := strings.TrimSpace(parts[3])        // 任务代码

			totalAccounts++

			// 解析消耗（去除逗号）
			cost := parseNumber(account.StatCost)         // 消耗
			cashCost := parseNumber(account.StatCashCost) // 先进消耗
			totalCost += cost
			totalCashCost += cashCost

			// 解析曝光数、点击数、转化数
			showCnt := parseInt64(account.ShowCnt) // 曝光
			totalShowCnt += showCnt
			clickCnt := parseInt64(account.ClickCnt) // 点击
			totalClickCnt += clickCnt
			convertCnt := parseInt64(account.ConvertCnt) // 转化
			totalConvertCnt += convertCnt

			// 累加转化成本
			conversionCost := parseNumber(account.ConversionCost)
			totalConversionCost += conversionCost

			// 查询返点率（主体-端口）
			rebateKey := fmt.Sprintf("%s-%s", subject, port)
			rebateRate := rebateMap[rebateKey]

			// 计算返点消耗 = 现金消耗 / (各端口各主体对应的返点率)
			var rebateCost float64
			if rebateRate > 0 {
				rebateCost = cashCost / rebateRate
			} else {
				rebateCost = cashCost
			}
			totalRebateCost += rebateCost

			// 查询服务费率
			serviceFeeRate := serviceFeeMap[serviceProvider]
			// 计算服务商成本 = 现金消耗 * 服务费率（如果服务费是0，则现金消耗*0.04）
			var serviceFeeCost float64
			if serviceFeeRate > 0 {
				serviceFeeCost = cashCost * serviceFeeRate
			} else {
				serviceFeeCost = cashCost
			}
			totalServiceFeeCost += serviceFeeCost

			// 查询结算单价
			settlementPrice := taskTypeMap[taskCode]

			// 计算预估收入 = (转化数+扣量数) * 结算单价
			// 注：结算单价不固定，主体或任务不一样对应的结算单价也不一样
			revenue := float64(convertCnt) * settlementPrice
			totalRevenue += revenue

			// 计算预估利润 = (预估收入 * 0.95) - 服务商成本 - 返点消耗
			profit := (revenue * 0.95) - serviceFeeCost - rebateCost
			totalProfit += profit

			// 计算预估利润率 = 预估利润/预估收入
			profitRate := profit / revenue

			logx.Infof("账户 %s (%d) [%s-%s-%s-%s]: 消耗=%.2f, 现金消耗=%.2f, 返点消耗=%.2f, 曝光=%d, 点击=%d, 转化=%d, 转化成本=%.2f, 转化率=%s, 服务费=%.2f, 预估收入=%.2f, 预估利润=%.2f, 预估利润率=%.2f",
				account.AdvertiserName, account.AdvertiserId, subject, port, serviceProvider, taskCode,
				cost, cashCost, rebateCost, showCnt, clickCnt, convertCnt, conversionCost, account.ConversionRate, serviceFeeCost, revenue, profit, profitRate)

			// 保存账户数据
			accountReports = append(accountReports, AccountReportData{
				AdvertiserId:    account.AdvertiserId,
				AdvertiserName:  account.AdvertiserName,
				Subject:         subject,
				Port:            port,
				ServiceProvider: serviceProvider,
				TaskCode:        taskCode,
				Cost:            cost,
				CashCost:        cashCost,
				RebateCost:      rebateCost,
				ShowCnt:         showCnt,
				ClickCnt:        clickCnt,
				Ctr:             account.Ctr,
				ConvertCnt:      convertCnt,
				ConversionCost:  account.ConversionCost,
				ConversionRate:  account.ConversionRate,
				ServiceFeeCost:  serviceFeeCost,
				Revenue:         revenue,
				Profit:          profit,
				ProfitRate:      profitRate,
			})
		}

		// 检查是否还有更多数据
		hasMore = resp.Data.Pagination.HasMore
		page++

		logx.Infof("已处理第 %d 页，本页账户数: %d，累计账户数: %d",
			page-1, len(resp.Data.DataList), totalAccounts)
	}

	// 打印汇总数据
	if totalAccounts == 0 {
		logx.Infof("今日暂无巨量账户数据，跳过的账户数: %d", skippedAccounts)
		return
	}

	// 计算总点击率
	var avgCtr float64
	if totalShowCnt > 0 {
		avgCtr = float64(totalClickCnt) / float64(totalShowCnt) * 100
	}

	// 计算平均转化成本
	var avgConversionCost float64
	if totalConvertCnt > 0 {
		avgConversionCost = totalConversionCost / float64(totalConvertCnt)
	}

	// 计算总转化率
	var avgConversionRate float64
	if totalClickCnt > 0 {
		avgConversionRate = float64(totalConvertCnt) / float64(totalClickCnt) * 100
	}

	// 计算预估利润率
	var profitRate float64
	if totalRevenue > 0 {
		profitRate = (totalProfit / totalRevenue) * 100
	}

	logx.Infof("巨量报表汇总 - 账户数: %d, 跳过: %d, 总消耗: %.2f, 现金消耗: %.2f, 返点消耗: %.2f, 总曝光: %d, 总点击: %d, 点击率: %.2f%%, 总转化: %d, 平均转化成本: %.2f, 转化率: %.2f%%, 服务费成本: %.2f, 预估收入: %.2f, 预估利润: %.2f, 利润率: %.2f%%",
		totalAccounts, skippedAccounts, totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt, avgCtr, totalConvertCnt,
		avgConversionCost, avgConversionRate, totalServiceFeeCost, totalRevenue, totalProfit, profitRate)

	logx.Infof("已保存 %d 条账户数据，待后续生成Excel报表", len(accountReports))

	// TODO: 后续在这里调用生成Excel并上传到服务器的函数，获取下载URL
	// excelDownloadURL := generateAndUploadExcelReport(ctx, accountReports)
	excelDownloadURL := "" // 临时占位符，等实现Excel生成和上传后填入真实URL

	// 发送钉钉通知
	sendJuliangDingTalkNotification(ctx, dingTalk, totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt,
		totalConvertCnt, avgConversionCost, avgConversionRate, avgCtr, totalAccounts, totalServiceFeeCost, totalRevenue, totalProfit, profitRate, excelDownloadURL)

	logx.Infof("巨量报表任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// sendJuliangDingTalkNotification 发送巨量钉钉通知
func sendJuliangDingTalkNotification(ctx context.Context, dingConfig config.DingTalkConfig,
	totalCost, totalCashCost, totalRebateCost float64, totalShowCnt, totalClickCnt, totalConvertCnt int64,
	avgConversionCost, avgConversionRate, avgCtr float64, totalAccounts int,
	totalServiceFeeCost, totalRevenue, totalProfit, profitRate float64, excelDownloadURL string) {

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
			"**总消耗**：%.2f  \n"+
			"**现金消耗**：%.2f  \n"+
			"**返后消耗**：%.2f  \n"+
			"**曝光量**：%d  \n"+
			"**点击量**：%d  \n"+
			"**点击率**：%.2f%%  \n"+
			"**注册转化数**：%d  \n"+
			"**转化成本**：%.2f  \n"+
			"**转化率**：%.2f%%  \n"+
			"**服务费成本**：%.2f  \n"+
			"**预估收入**：%.2f  \n"+
			"**预估利润**：%.2f  \n"+
			"**利润率**：%.2f%%  \n\n"+
			"详细账户信息请下载文件：[下载](%s)",
		timeStr,
		totalAccounts,
		totalCost,
		totalCashCost,
		totalRebateCost,
		totalShowCnt,
		totalClickCnt,
		avgCtr,
		totalConvertCnt,
		avgConversionCost,
		avgConversionRate,
		totalServiceFeeCost,
		totalRevenue,
		totalProfit,
		profitRate,
		excelDownloadURL,
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
