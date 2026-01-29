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

// executeJuliangReportJobV2 执行巨量报表任务（使用下载接口方式）
func executeJuliangReportJobV2(db *gorm.DB, dingTalk config.DingTalkConfig, fileServer config.FileServerConfig) {
	ctx := context.Background()
	logx.Infof("开始执行巨量报表任务V2 - %s", time.Now().Format("2006-01-02 15:04:05"))

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

	// 第一步：调用下载接口获取 download_id
	logx.Info("第一步：调用下载接口获取 download_id")
	downloadId, err := requestDownloadTask(ctx, client, startTime, endTime)
	if err != nil {
		logx.Errorf("请求下载任务失败: %v", err)
		return
	}
	logx.Infof("获取到 download_id: %s", downloadId)

	// 第二步：轮询下载任务列表，获取 schedulerId
	logx.Info("第二步：轮询下载任务列表，等待任务完成")
	schedulerId, err := waitForDownloadTask(ctx, client, downloadId, startOfDay, endOfDay)
	if err != nil {
		logx.Errorf("等待下载任务失败: %v", err)
		return
	}
	logx.Infof("获取到 schedulerId: %d", schedulerId)

	// 第三步：下载 Excel 文件
	logx.Info("第三步：下载 Excel 文件")
	excelPath, err := downloadExcelFile(ctx, client, schedulerId, fileServer)
	if err != nil {
		logx.Errorf("下载 Excel 文件失败: %v", err)
		return
	}
	logx.Infof("Excel 文件已下载: %s", excelPath)

	// 读取 Excel 文件并处理数据
	logx.Info("第四步：读取 Excel 文件并处理数据")
	accountReports, err := processExcelData(ctx, excelPath, db, attributionMap)
	if err != nil {
		logx.Errorf("处理 Excel 数据失败: %v", err)
		return
	}

	if len(accountReports) == 0 {
		logx.Info("Excel 文件中没有有效的账户数据")
		return
	}

	// 计算汇总数据并发送钉钉通知
	sendSummaryNotification(ctx, accountReports, dingTalk, fileServer)

	logx.Infof("巨量报表任务V2执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// requestDownloadTask 第一步：请求下载任务，获取 download_id
func requestDownloadTask(ctx context.Context, client *httpclient.Client, startTime, endTime string) (string, error) {
	req := map[string]interface{}{
		"start_time":   startTime,
		"end_time":     endTime,
		"offset":       1,
		"limit":        100,
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
		"download":         true, // 关键参数
	}

	var resp types.JuliangApiResponse
	err := client.Post(ctx, "/nbs/api/bm/promotion/ad/get_account_list", req, &resp)
	if err != nil {
		return "", fmt.Errorf("调用下载接口失败: %w", err)
	}

	if resp.Code != 0 {
		return "", fmt.Errorf("下载接口返回错误: code=%d, message=%s", resp.Code, resp.Msg)
	}

	downloadId := resp.Data.DownloadId
	if downloadId == "" {
		return "", fmt.Errorf("未获取到 download_id")
	}

	// 去掉 "V3" 前缀
	downloadId = strings.TrimPrefix(downloadId, "V3")

	return downloadId, nil
}

// waitForDownloadTask 第二步：轮询下载任务列表，等待任务完成
func waitForDownloadTask(ctx context.Context, client *httpclient.Client, downloadId string, startOfDay, endOfDay time.Time) (int64, error) {
	// 最多等待 5 分钟
	maxRetries := 60
	retryInterval := 5 * time.Second

	stTimeStr := startOfDay.Format("2006-01-02 15:04:05")
	etTimeStr := endOfDay.Format("2006-01-02 15:04:05")

	for i := 0; i < maxRetries; i++ {
		req := map[string]interface{}{
			"status": []int{2}, // 2 表示任务成功
			"type":   []int{5}, // 5 表示 bp推广管理下载
			"st":     stTimeStr,
			"et":     etTimeStr,
			"order": []map[string]interface{}{
				{
					"orderField": "task_create_time",
					"orderType":  1,
				},
			},
			"page": 1,
			"size": 20,
		}

		var resp types.DownloadTaskListResponse
		err := client.Post(ctx, "/nbs/api/bm/task_center/download/list", req, &resp)
		if err != nil {
			logx.Errorf("查询下载任务列表失败 (第%d次尝试): %v", i+1, err)
			time.Sleep(retryInterval)
			continue
		}

		if resp.Code != 0 {
			logx.Errorf("下载任务列表接口返回错误 (第%d次尝试): code=%d, message=%s", i+1, resp.Code, resp.Msg)
			time.Sleep(retryInterval)
			continue
		}

		// 查找匹配的任务
		for _, task := range resp.Data.List {
			if task.TaskId == downloadId && task.SchedulerStatus == 2 {
				logx.Infof("找到匹配的下载任务: taskId=%s, schedulerId=%d, status=%s",
					task.TaskId, task.SchedulerId, task.SchedulerStatusName)
				return task.SchedulerId, nil
			}
		}

		logx.Infof("未找到匹配的下载任务 (第%d次尝试)，等待%v后重试...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return 0, fmt.Errorf("等待下载任务超时，未找到 taskId=%s 的已完成任务", downloadId)
}

// downloadExcelFile 第三步：下载 Excel 文件
func downloadExcelFile(ctx context.Context, client *httpclient.Client, schedulerId int64, fileServer config.FileServerConfig) (string, error) {
	req := map[string]interface{}{
		"operation":     1,
		"schedulerType": 5,
		"schedulerId":   strconv.FormatInt(schedulerId, 10),
	}

	// 确保保存目录存在
	savePath := fileServer.Path
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("创建报表目录失败: %w", err)
	}

	// 生成文件名
	now := time.Now()
	filename := fmt.Sprintf("juliang_download_%s.xlsx", now.Format("20060102_150405"))
	filepath := filepath.Join(savePath, filename)

	// 下载文件
	err := client.DownloadFile(ctx, "/nbs/api/bm/task_center/download/download_operation/", req, filepath)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}

	return filepath, nil
}

// processExcelData 处理 Excel 数据
func processExcelData(ctx context.Context, excelPath string, db *gorm.DB, attributionMap map[string]int64) ([]AccountReportData, error) {
	// 打开 Excel 文件
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("打开 Excel 文件失败: %w", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 文件中没有工作表")
	}
	sheetName := sheets[0]

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %w", err)
	}

	if len(rows) <= 1 {
		logx.Info("Excel 文件中没有数据行（只有表头或为空）")
		return []AccountReportData{}, nil
	}

	// 预加载配置数据
	rebateMap, err := model.LoadRebateConfigMap(db)
	if err != nil {
		return nil, fmt.Errorf("加载返点配置失败: %w", err)
	}
	serviceFeeMap, err := model.LoadServiceFeeConfigMap(db)
	if err != nil {
		return nil, fmt.Errorf("加载服务费配置失败: %w", err)
	}
	taskTypeMap, err := model.LoadTaskTypeConfigMap(db)
	if err != nil {
		return nil, fmt.Errorf("加载任务类型配置失败: %w", err)
	}

	var accountReports []AccountReportData
	var skippedAccounts int

	// 解析表头，找到各列的索引
	header := rows[0]
	colIndexMap := make(map[string]int)
	for i, col := range header {
		colIndexMap[col] = i
	}

	// 处理数据行（从第2行开始）
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// 获取单元格值的辅助函数
		getCellValue := func(colName string) string {
			if idx, exists := colIndexMap[colName]; exists && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		// 读取备注字段
		remark := getCellValue("账户备注")
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

		// 读取账户信息
		advertiserIdStr := getCellValue("账户ID")
		advertiserName := getCellValue("账户")

		// 读取数值字段
		cost := parseNumber(getCellValue("消耗"))
		cashCost := parseNumber(getCellValue("现金消耗"))
		showCnt := parseInt64(getCellValue("展示数"))
		clickCnt := parseInt64(getCellValue("点击数"))
		ctr := getCellValue("点击率(%)")
		convertCnt := parseInt64(getCellValue("转化数"))
		conversionCost := getCellValue("转化成本")
		conversionRate := getCellValue("转化率(%)")

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

		// 计算返点消耗 = 现金消耗 / (各端口各主体对应的返点率) - 使用消耗数据计算
		var rebateCost float64
		if rebateRate > 0 {
			rebateCost = cost / rebateRate
		} else {
			rebateCost = cost
		}

		// 计算服务商成本 = 现金消耗 * 服务费率
		var serviceFeeCost float64
		if serviceFeeRate > 0 {
			serviceFeeCost = cost * serviceFeeRate
		} else {
			serviceFeeCost = cost
		}

		// 获取归因扣量数据
		deductionCount := attributionMap[advertiserIdStr]

		// 计算预估收入 = (转化数+扣量数) * 结算单价
		revenue := float64(convertCnt+deductionCount) * settlementPrice

		// 计算预估利润 = (预估收入 * 0.95) - 服务商成本 - 返点消耗
		profit := (revenue * 0.95) - serviceFeeCost - rebateCost

		// 计算预估利润率 = 预估利润/预估收入
		var profitRate float64
		if revenue > 0 {
			profitRate = profit / revenue
		}

		// 保存账户数据
		accountReports = append(accountReports, AccountReportData{
			AdvertiserId:    advertiserIdStr,
			AdvertiserName:  advertiserName,
			Subject:         subject,
			Port:            port,
			ServiceProvider: serviceProvider,
			TaskCode:        taskName,
			Cost:            cost,
			CashCost:        cashCost,
			RebateCost:      rebateCost,
			ShowCnt:         showCnt,
			ClickCnt:        clickCnt,
			Ctr:             ctr,
			ConvertCnt:      convertCnt,
			ConversionCost:  conversionCost,
			ConversionRate:  conversionRate,
			ServiceFeeCost:  serviceFeeCost,
			Revenue:         revenue,
			Profit:          profit,
			ProfitRate:      profitRate,
		})
	}

	logx.Infof("从 Excel 文件中读取 %d 条账户数据，跳过 %d 条", len(accountReports), skippedAccounts)
	return accountReports, nil
}

// sendSummaryNotification 计算汇总数据并发送钉钉通知
func sendSummaryNotification(ctx context.Context, accountReports []AccountReportData, dingTalk config.DingTalkConfig, fileServer config.FileServerConfig) {
	var totalCost float64
	var totalCashCost float64
	var totalRebateCost float64
	var totalShowCnt int64
	var totalClickCnt int64
	var totalConvertCnt int64
	var totalConversionCost float64
	var totalServiceFeeCost float64
	var totalRevenue float64
	var totalProfit float64

	for _, report := range accountReports {
		totalCost += report.Cost
		totalCashCost += report.CashCost
		totalRebateCost += report.RebateCost
		totalShowCnt += report.ShowCnt
		totalClickCnt += report.ClickCnt
		totalConvertCnt += report.ConvertCnt
		totalConversionCost += parseNumber(report.ConversionCost)
		totalServiceFeeCost += report.ServiceFeeCost
		totalRevenue += report.Revenue
		totalProfit += report.Profit
	}

	totalAccounts := len(accountReports)

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

	// 生成Excel报表并获取下载URL
	excelDownloadURL := generateAndUploadExcelReport(ctx, accountReports, fileServer,
		totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt, avgCtr, totalConvertCnt, avgConversionCost, avgConversionRate,
		totalServiceFeeCost, totalRevenue, totalProfit, profitRate)

	// 发送钉钉通知
	sendJuliangDingTalkNotification(ctx, dingTalk, totalCost, totalCashCost, totalRebateCost, totalShowCnt, totalClickCnt,
		totalConvertCnt, avgConversionCost, avgConversionRate, avgCtr, totalAccounts, totalServiceFeeCost, totalRevenue, totalProfit, profitRate, 0, excelDownloadURL)
}
