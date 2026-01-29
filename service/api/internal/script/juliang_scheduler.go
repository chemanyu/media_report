package script

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

// AccountReportData 账户报表数据
type AccountReportData struct {
	AdvertiserId    string  // 账户ID
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

// ExecuteJuliangReportJob 执行巨量报表任务 (导出供外部调用)
func ExecuteJuliangReportJob(db *gorm.DB, dingTalk config.DingTalkConfig, fileServer config.FileServerConfig) {
	executeJuliangReportJob(db, dingTalk, fileServer)
}

// executeJuliangReportJob 执行巨量报表任务
func executeJuliangReportJob(db *gorm.DB, dingTalk config.DingTalkConfig, fileServer config.FileServerConfig) {
	ctx := context.Background()
	logx.Infof("开始执行巨量报表任务 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 调用归因接口获取扣量数据
	attributionMap := fetchAttributionData(ctx)

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

	// 第一步：先请求第一页获取总数
	limit := 100
	req := map[string]interface{}{
		"start_time":   startTime,
		"end_time":     endTime,
		"offset":       1,
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

	var firstResp types.JuliangApiResponse
	err = client.Post(ctx, "/nbs/api/bm/promotion/ad/get_account_list", req, &firstResp)
	if err != nil {
		logx.Errorf("调用巨量 API 获取第一页失败: %v", err)
		return
	}

	if firstResp.Code != 0 {
		logx.Errorf("巨量 API 返回错误: code=%d, message=%s", firstResp.Code, firstResp.Msg)
		return
	}

	// 获取总数和计算总页数
	total := firstResp.Data.Pagination.Total
	if total == 0 {
		logx.Info("今日暂无巨量账户数据")
		return
	}

	totalPages := (total + limit - 1) / limit
	logx.Infof("共有 %d 条账户数据，分 %d 页，开始并发请求...", total, totalPages)

	// 第二步：并发请求所有页面
	type pageResult struct {
		page int
		data []types.JuliangAccountData
		err  error
	}

	resultChan := make(chan pageResult, totalPages)
	semaphore := make(chan struct{}, 10) // 限制并发数为10

	// 处理第一页数据（已经获取到了）
	go func() {
		resultChan <- pageResult{page: 1, data: firstResp.Data.DataList, err: nil}
	}()

	// 并发请求剩余页面
	for page := 2; page <= totalPages; page++ {
		semaphore <- struct{}{} // 获取信号量
		go func(p int) {
			defer func() { <-semaphore }() // 释放信号量

			pageReq := map[string]interface{}{
				"start_time":   startTime,
				"end_time":     endTime,
				"offset":       p,
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

			var pageResp types.JuliangApiResponse
			err := client.Post(ctx, "/nbs/api/bm/promotion/ad/get_account_list", pageReq, &pageResp)
			if err != nil {
				resultChan <- pageResult{page: p, data: nil, err: err}
				return
			}

			if pageResp.Code != 0 {
				resultChan <- pageResult{page: p, data: nil, err: fmt.Errorf("API返回错误: code=%d, message=%s", pageResp.Code, pageResp.Msg)}
				return
			}

			resultChan <- pageResult{page: p, data: pageResp.Data.DataList, err: nil}
		}(page)
	}

	// 第三步：收集所有结果
	allAccounts := make(map[int][]types.JuliangAccountData)
	for i := 0; i < totalPages; i++ {
		result := <-resultChan
		if result.err != nil {
			logx.Errorf("第 %d 页请求失败: %v", result.page, result.err)
			continue
		}
		allAccounts[result.page] = result.data
		logx.Infof("第 %d/%d 页数据获取成功，账户数: %d", result.page, totalPages, len(result.data))
	}
	close(resultChan)

	logx.Infof("所有页面请求完成，开始处理账户数据...")

	// 第四步：按页码顺序处理所有账户数据
	for page := 1; page <= totalPages; page++ {
		accounts, exists := allAccounts[page]
		if !exists {
			continue
		}

		for _, account := range accounts {
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
			taskName := strings.TrimSpace(parts[3])        // 任务代码

			// 查询返点率（主体-端口）
			rebateKey := fmt.Sprintf("%s-%s", subject, port)
			rebateRate, rebateExists := rebateMap[rebateKey]

			// 查询服务费率
			serviceFeeRate, serviceFeeExists := serviceFeeMap[serviceProvider]

			// 查询结算单价
			settlementPrice, taskTypeExists := taskTypeMap[taskName]

			// 校验：如果主体-端口、服务商、任务不在数据库配置中，跳过此条数据
			if !rebateExists {
				skippedAccounts++
				continue
			}
			if !serviceFeeExists {
				skippedAccounts++
				continue
			}
			if !taskTypeExists {
				skippedAccounts++
				continue
			}

			// 只有通过所有校验后才计入总账户数
			totalAccounts++

			// 解析消耗（去除逗号）
			cost := parseNumber(account.StatCost)         // 消耗
			cashCost := parseNumber(account.StatCashCost) // 现金消耗
			totalCost += cost
			totalCashCost += cashCost

			// 解析曝光数、点击数、转化数
			showCnt := parseInt64(account.ShowCnt)       // 曝光
			clickCnt := parseInt64(account.ClickCnt)     // 点击
			convertCnt := parseInt64(account.ConvertCnt) // 转化
			totalShowCnt += showCnt
			totalClickCnt += clickCnt
			totalConvertCnt += convertCnt

			// 累加转化成本
			conversionCost := parseNumber(account.ConversionCost)
			totalConversionCost += conversionCost

			// 计算返点消耗 = 消耗 / (1 + 返点率)
			// 例如：返点率0.04（4个点），则 消耗/1.04
			var rebateCost float64
			if rebateRate > 0 {
				rebateCost = cost / (1 + rebateRate)
			} else {
				rebateCost = cost
			}
			totalRebateCost += rebateCost

			// 计算服务费 = 返点消耗 * 服务费率
			var serviceFeeCost float64
			if serviceFeeRate > 0 {
				serviceFeeCost = cost * serviceFeeRate
			} else {
				serviceFeeCost = cost
			}
			totalServiceFeeCost += serviceFeeCost

			// 获取归因扣量数据 (advertiser_rate_false_4)
			advertiserIdStr := strconv.FormatInt(account.AdvertiserId, 10)
			deductionCount := attributionMap[advertiserIdStr]

			// 计算预估收入 = (转化数+扣量数) * 结算单价
			// 注：结算单价不固定，主体或任务不一样对应的结算单价也不一样
			revenue := float64(convertCnt+deductionCount) * settlementPrice
			totalRevenue += revenue

			// 计算预估利润 = (预估收入 * 0.95) - 服务商成本 - 返点消耗
			profit := (revenue * 0.95) - serviceFeeCost - rebateCost
			totalProfit += profit

			// 计算预估利润率 = 预估利润/预估收入
			profitRate := profit / revenue

			// logx.Infof("账户 %s (%d) [%s-%s-%s-%s]: 消耗=%.2f, 现金消耗=%.2f, 返点消耗=%.2f, 曝光=%d, 点击=%d, 转化=%d, 转化成本=%.2f, 转化率=%s, 服务费=%.2f, 预估收入=%.2f, 预估利润=%.2f, 预估利润率=%.2f",
			// 	account.AdvertiserName, account.AdvertiserId, subject, port, serviceProvider, taskCode,
			// 	cost, cashCost, rebateCost, showCnt, clickCnt, convertCnt, conversionCost, account.ConversionRate, serviceFeeCost, revenue, profit, profitRate)

			// 保存账户数据
			accountReports = append(accountReports, AccountReportData{
				AdvertiserId:    strconv.Itoa(int(account.AdvertiserId)),
				AdvertiserName:  account.AdvertiserName,
				Subject:         subject,
				Port:            port,
				ServiceProvider: serviceProvider,
				TaskCode:        taskName,
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
	}

	logx.Infof("数据处理完成 - 有效账户数: %d, 跳过账户数: %d", totalAccounts, skippedAccounts)

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

	logx.Infof("已保存 %d 条账户数据，待后续生成Excel报表", len(accountReports))

	// 生成Excel报表并获取下载URL
	excelDownloadURL := generateAndUploadExcelReport(ctx, accountReports, fileServer,
		totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt, avgCtr, totalConvertCnt, avgConversionCost, avgConversionRate,
		totalServiceFeeCost, totalRevenue, totalProfit, profitRate)

	// 发送钉钉通知
	sendJuliangDingTalkNotification(ctx, dingTalk, totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt,
		totalConvertCnt, avgConversionCost, avgConversionRate, avgCtr, totalAccounts, totalServiceFeeCost, totalRevenue, totalProfit, profitRate, skippedAccounts, excelDownloadURL)

	logx.Infof("巨量报表任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// sendJuliangDingTalkNotification 发送巨量钉钉通知
func sendJuliangDingTalkNotification(ctx context.Context, dingConfig config.DingTalkConfig,
	totalCost, totalCashCost, totalRebateCost float64, totalShowCnt, totalClickCnt, totalConvertCnt int64,
	avgConversionCost, avgConversionRate, avgCtr float64, totalAccounts int,
	totalServiceFeeCost, totalRevenue, totalProfit, profitRate float64, skippedAccounts int, excelDownloadURL string) {

	if !dingConfig.Enabled || dingConfig.JDReportWebhookURL == "" {
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
			"**转化数**：%d  \n"+
			"**转化成本**：%.2f  \n"+
			"**转化率**：%.2f%%  \n"+
			"**服务费成本**：%.2f  \n"+
			"**预估收入**：%.2f  \n"+
			"**预估利润**：%.2f  \n"+
			"**利润率**：%.2f%%  \n"+
			"**备注不符合标准跳过账户数**：%d  \n\n"+
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
		skippedAccounts,
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
	err := client.Post(ctx, dingConfig.JDReportWebhookURL, msg, &result)
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

// generateAndUploadExcelReport 生成Excel报表并返回下载URL
func generateAndUploadExcelReport(ctx context.Context, accountReports []AccountReportData, fileServer config.FileServerConfig,
	totalCost, totalCashCost, totalRebateCost float64, totalShowCnt, totalClickCnt int64, avgCtr float64, totalConvertCnt int64,
	avgConversionCost, avgConversionRate, totalServiceFeeCost, totalRevenue, totalProfit, profitRate float64) string {
	if len(accountReports) == 0 {
		logx.Error("账户数据为空，无法生成Excel报表")
		return ""
	}

	// 创建新的Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logx.Errorf("关闭Excel文件失败: %v", err)
		}
	}()

	sheetName := "巨量账户报表"
	// 创建工作表
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logx.Errorf("创建工作表失败: %v", err)
		return ""
	}

	// 设置表头
	headers := []string{
		"账户ID", "账户名称", "主体", "任务", "服务商",
		"消耗汇总", "现金消耗汇总", "返后消耗汇总",
		"曝光汇总", "点击汇总", "点击率汇总",
		"转化数汇总", "转化成本汇总", "转化率",
		"服务商成本", "预估收入", "预估利润", "预估利润率",
	}

	// 写入表头
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, report := range accountReports {
		row := i + 2 // 从第2行开始（第1行是表头）
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), report.AdvertiserId)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.AdvertiserName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), report.Subject)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), report.TaskCode)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), report.ServiceProvider)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), report.Cost)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), report.CashCost)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), report.RebateCost)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), report.ShowCnt)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), report.ClickCnt)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), report.Ctr)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), report.ConvertCnt)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), report.ConversionCost)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), report.ConversionRate)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), report.ServiceFeeCost)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), report.Revenue)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), report.Profit)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), fmt.Sprintf("%.2f%%", report.ProfitRate*100))
	}

	// 在最后一行添加汇总数据
	totalRow := len(accountReports) + 2
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", totalRow), "")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), "汇总")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", totalRow), "")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), "")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", totalRow), "")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", totalRow), totalCost)
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", totalRow), totalCashCost)
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", totalRow), totalRebateCost)
	f.SetCellValue(sheetName, fmt.Sprintf("I%d", totalRow), totalShowCnt)
	f.SetCellValue(sheetName, fmt.Sprintf("J%d", totalRow), totalClickCnt)
	f.SetCellValue(sheetName, fmt.Sprintf("K%d", totalRow), fmt.Sprintf("%.2f%%", avgCtr))
	f.SetCellValue(sheetName, fmt.Sprintf("L%d", totalRow), totalConvertCnt)
	f.SetCellValue(sheetName, fmt.Sprintf("M%d", totalRow), avgConversionCost)
	f.SetCellValue(sheetName, fmt.Sprintf("N%d", totalRow), avgConversionRate)
	f.SetCellValue(sheetName, fmt.Sprintf("O%d", totalRow), totalServiceFeeCost)
	f.SetCellValue(sheetName, fmt.Sprintf("P%d", totalRow), totalRevenue)
	f.SetCellValue(sheetName, fmt.Sprintf("Q%d", totalRow), totalProfit)
	f.SetCellValue(sheetName, fmt.Sprintf("R%d", totalRow), fmt.Sprintf("%.2f%%", profitRate))

	// 设置默认活动工作表
	f.SetActiveSheet(index)
	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 确保保存目录存在
	savePath := fileServer.Path
	if err := os.MkdirAll(savePath, 0755); err != nil {
		logx.Errorf("创建报表目录失败: %v", err)
		return ""
	}

	// 生成文件名（包含时间戳）
	now := time.Now()
	filename := fmt.Sprintf("juliang_report_%s.xlsx", now.Format("20060102_150405"))
	filepath := filepath.Join(savePath, filename)

	// 保存文件
	if err := f.SaveAs(filepath); err != nil {
		logx.Errorf("保存Excel文件失败: %v", err)
		return ""
	}

	logx.Infof("Excel报表已生成: %s", filepath)

	// 生成下载URL
	baseURL := fileServer.BaseURL
	downloadURL := fmt.Sprintf("%s/download/%s", baseURL, filename)

	return downloadURL
}

// fetchAttributionData 获取归因扣量数据
func fetchAttributionData(ctx context.Context) map[string]int64 {
	// 创建 HTTP 客户端
	client := httpclient.NewClient("http://ad-ocpx.zhltech.net", 30)

	// 调用归因接口
	var resp types.AttributionResponse
	err := client.Get(ctx, "/attribution/data", nil, &resp)
	if err != nil {
		logx.Errorf("调用归因接口失败: %v", err)
		return make(map[string]int64)
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("归因接口返回错误: code=%d, message=%s", resp.Code, resp.Message)
		return make(map[string]int64)
	}

	// 构建账户ID -> advertiser_rate_false_4 的映射
	attributionMap := make(map[string]int64)
	for advertiserId, errorCounts := range resp.Data.ErrorCounts {
		if count, exists := errorCounts["advertiser_rate_false_4"]; exists {
			attributionMap[advertiserId] = count
		}
	}

	logx.Infof("成功获取归因数据，共 %d 个账户有扣量记录", len(attributionMap))
	return attributionMap
}
