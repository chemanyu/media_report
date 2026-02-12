package script

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	fzlogic "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
)

// ExecuteFzDataSyncJob 执行飞猪外投数据同步任务
func ExecuteFzDataSyncJob(db *gorm.DB, config config.Config) {
	ctx := context.Background()
	logx.Infof("开始执行飞猪外投数据同步任务 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 创建服务上下文
	svcCtx := &svc.ServiceContext{
		Config: config,
		DB:     db,
	}

	// 创建业务逻辑实例
	l := fzlogic.NewFzHourlyReportLogic(ctx, svcCtx)

	var oppoCount, xiaomiCount int
	var oppoErr, xiaomiErr error

	// 同步今天的OPPO数据
	oppoCount, oppoErr = l.SyncTodayOppoData()
	if oppoErr != nil {
		logx.Errorf("同步OPPO数据失败: %v", oppoErr)
	} else {
		logx.Infof("成功同步OPPO数据，账户数量: %d", oppoCount)
	}

	// 同步今天的小米数据
	xiaomiCount, xiaomiErr = l.SyncTodayXiaomiData()
	if xiaomiErr != nil {
		logx.Errorf("同步小米数据失败: %v", xiaomiErr)
	} else {
		logx.Infof("成功同步小米数据，账户数量: %d", xiaomiCount)
	}

	// 汇总结果
	if oppoErr == nil && xiaomiErr == nil {
		logx.Infof("飞猪外投数据同步任务执行完成 - OPPO: %d个账户, 小米: %d个账户 - %s",
			oppoCount, xiaomiCount, time.Now().Format("2006-01-02 15:04:05"))
	} else if oppoErr != nil && xiaomiErr != nil {
		logx.Errorf("飞猪外投数据同步任务失败 - OPPO和小米数据均同步失败")
	} else {
		logx.Infof("飞猪外投数据同步任务部分完成 - OPPO: %d个账户, 小米: %d个账户 - %s",
			oppoCount, xiaomiCount, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 发送钉钉通知
	SendFzDingTalkNotification(ctx, db, config.DingTalk)
}

// sendFzDingTalkNotification 发送飞猪数据钉钉通知
func SendFzDingTalkNotification(ctx context.Context, db *gorm.DB, dingConfig config.DingTalkConfig) {
	// 获取今天的日期（格式：20260212）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	today := now.Format("20060102")
	todayInt, _ := strconv.Atoi(today)

	// 格式化日期用于显示（2月12日-14时）
	displayDate := fmt.Sprintf("%d月%d日-%d时", now.Month(), now.Day(), now.Hour())

	// 从数据库查询今天的所有数据
	reportModel := model.NewFzHourlyReportModel(db)
	reports, err := reportModel.FindByDateRange("", todayInt, todayInt)
	if err != nil {
		logx.Errorf("查询今天的飞猪数据失败: %v", err)
		return
	}

	if len(reports) == 0 {
		logx.Info("今天暂无飞猪数据，跳过发送钉钉通知")
		return
	}

	// 分组统计：集能量活动 vs 常规活动
	var regularCost, regularConvertDp, regularDpAppOrderNums float64
	var energyCost, energyConvertDp, energyDpAppOrderNums float64
	var regularCount, energyCount int64

	for _, report := range reports {
		if report.MediaAdvId == "1073770254_348117" {
			// 集能量活动
			energyCost += report.Cost
			energyConvertDp += float64(report.ConvertDp)
			energyDpAppOrderNums += float64(report.DpAppOrderNums)
			energyCount++
		} else {
			// 常规活动
			regularCost += report.Cost
			regularConvertDp += float64(report.ConvertDp)
			regularDpAppOrderNums += float64(report.DpAppOrderNums)
			regularCount++
		}
	}

	// 计算成本（数据库中是分，需要除以100）
	// 实际消耗（元）
	regularCostYuan := regularCost / 100
	// 现金消耗 = 实际消耗 * 1.2 * 0.85
	regularCashCost := regularCostYuan * 1.2 * 0.85
	regularConvertDpPrice := 0.0
	if regularConvertDp > 0 {
		regularConvertDpPrice = regularCashCost / regularConvertDp
	}
	regularDpOrderPrice := 0.0
	if regularDpAppOrderNums > 0 {
		regularDpOrderPrice = regularCashCost / regularDpAppOrderNums
	}

	// 实际消耗（元）
	energyCostYuan := energyCost / 100
	// 现金消耗 = 实际消耗 * 1.2 * 0.85
	energyCashCost := energyCostYuan * 1.2 * 0.85
	energyConvertDpPrice := 0.0
	if energyConvertDp > 0 {
		energyConvertDpPrice = energyCashCost / energyConvertDp
	}

	// 汇总数据（常规 + 集能量）
	totalConvertDp := regularConvertDp + energyConvertDp
	totalDpAppOrderNums := regularDpAppOrderNums + energyDpAppOrderNums
	totalCostYuan := (regularCost + energyCost) / 100
	totalCashCost := totalCostYuan * 1.2 * 0.85
	totalConvertDpPrice := 0.0
	if totalConvertDp > 0 {
		totalConvertDpPrice = totalCashCost / totalConvertDp
	}
	totalDpOrderPrice := 0.0
	if totalDpAppOrderNums > 0 {
		totalDpOrderPrice = totalCashCost / totalDpAppOrderNums
	}

	// 构建钉钉消息
	var markdownText string

	// 添加标头
	markdownText += "#### 飞猪活动时报  \n---\n"

	// 汇总数据
	if regularCount+energyCount > 0 {
		markdownText += fmt.Sprintf(
			"**汇总-海纳【飞猪app拉活 %s简报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 6500）  \n"+
				"**唤起成本**：%.2f（考核 0.5）  \n"+
				"**下单 pv 成本**：%.2f（考核 35）  \n\n",
			displayDate,
			int64(totalConvertDp),
			totalCashCost,
			totalConvertDpPrice,
			totalDpOrderPrice,
		)
		if regularCount > 0 || energyCount > 0 {
			markdownText += "---\n"
		}
	}

	// 常规活动数据
	if regularCount > 0 {
		markdownText += fmt.Sprintf(
			"**常规活动-海纳【飞猪常规活动 %s简报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 6000）  \n"+
				"**唤起成本**：%.2f（考核 0.5）  \n"+
				"**下单 pv 成本**：%.2f（考核 35）  \n\n",
			displayDate,
			int64(regularConvertDp),
			regularCashCost,
			regularConvertDpPrice,
			regularDpOrderPrice,
		)
		if energyCount > 0 {
			markdownText += "---\n"
		}
	}

	// 集能量活动数据
	if energyCount > 0 {
		markdownText += fmt.Sprintf(
			"**集能量-海纳【飞猪集能量 %s简报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 500）  \n"+
				"**唤起成本**：%.2f（考核 0.2）  \n",
			displayDate,
			int64(energyConvertDp),
			energyCashCost,
			energyConvertDpPrice,
		)
	}

	if markdownText == "" {
		logx.Info("没有需要发送的飞猪数据")
		return
	}

	// 检查钉钉配置
	if !dingConfig.Enabled || dingConfig.FzWebhookURL == "" {
		logx.Info("飞猪钉钉通知未启用或未配置webhook，跳过发送")
		return
	}

	// 发送到钉钉
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "飞猪外投简报",
			"text":  markdownText,
		},
	}

	client := httpclient.NewClient("", 30)
	var result map[string]interface{}
	err = client.Post(ctx, dingConfig.FzWebhookURL, msg, &result)
	if err != nil {
		logx.Errorf("发送飞猪钉钉消息失败: %v", err)
		return
	}

	logx.Infof("飞猪钉钉消息发送成功: %v", result)
}

// SendFzDailyReport 发送飞猪昨日日报
func SendFzDailyReport(ctx context.Context, db *gorm.DB, dingConfig config.DingTalkConfig) {
	// 获取昨天的日期（格式：20260211）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("20060102")
	yesterdayInt, _ := strconv.Atoi(yesterdayStr)

	// 格式化日期用于显示（2月11日）
	displayDate := fmt.Sprintf("%d月%d日", yesterday.Month(), yesterday.Day())

	// 从数据库查询昨天的所有数据
	reportModel := model.NewFzHourlyReportModel(db)
	reports, err := reportModel.FindByDateRange("", yesterdayInt, yesterdayInt)
	if err != nil {
		logx.Errorf("查询昨天的飞猪数据失败: %v", err)
		return
	}

	if len(reports) == 0 {
		logx.Info("昨天暂无飞猪数据，跳过发送钉钉日报")
		return
	}

	// 分组统计：集能量活动 vs 常规活动
	var regularCost, regularConvertDp, regularDpAppOrderNums float64
	var energyCost, energyConvertDp, energyDpAppOrderNums float64
	var regularCount, energyCount int64

	for _, report := range reports {
		if report.MediaAdvId == "1073770254_348117" {
			// 集能量活动
			energyCost += report.Cost
			energyConvertDp += float64(report.ConvertDp)
			energyDpAppOrderNums += float64(report.DpAppOrderNums)
			energyCount++
		} else {
			// 常规活动
			regularCost += report.Cost
			regularConvertDp += float64(report.ConvertDp)
			regularDpAppOrderNums += float64(report.DpAppOrderNums)
			regularCount++
		}
	}

	// 计算成本（数据库中是分，需要除以100）
	// 实际消耗（元）
	regularCostYuan := regularCost / 100
	// 现金消耗 = 实际消耗 * 1.2 * 0.85
	regularCashCost := regularCostYuan * 1.2 * 0.85
	regularConvertDpPrice := 0.0
	if regularConvertDp > 0 {
		regularConvertDpPrice = regularCashCost / regularConvertDp
	}
	regularDpOrderPrice := 0.0
	if regularDpAppOrderNums > 0 {
		regularDpOrderPrice = regularCashCost / regularDpAppOrderNums
	}

	// 实际消耗（元）
	energyCostYuan := energyCost / 100
	// 现金消耗 = 实际消耗 * 1.2 * 0.85
	energyCashCost := energyCostYuan * 1.2 * 0.85
	energyConvertDpPrice := 0.0
	if energyConvertDp > 0 {
		energyConvertDpPrice = energyCashCost / energyConvertDp
	}

	// 汇总数据（常规 + 集能量）
	totalConvertDp := regularConvertDp + energyConvertDp
	totalDpAppOrderNums := regularDpAppOrderNums + energyDpAppOrderNums
	totalCostYuan := (regularCost + energyCost) / 100
	totalCashCost := totalCostYuan * 1.2 * 0.85
	totalConvertDpPrice := 0.0
	if totalConvertDp > 0 {
		totalConvertDpPrice = totalCashCost / totalConvertDp
	}
	totalDpOrderPrice := 0.0
	if totalDpAppOrderNums > 0 {
		totalDpOrderPrice = totalCashCost / totalDpAppOrderNums
	}

	// 构建钉钉消息
	var markdownText string

	// 添加标头
	markdownText += "#### 飞猪活动日报  \n---\n"

	// 汇总数据
	if regularCount+energyCount > 0 {
		markdownText += fmt.Sprintf(
			"**汇总-海纳【飞猪app拉活 %s日报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 6500）  \n"+
				"**唤起成本**：%.2f（考核 0.5）  \n"+
				"**下单 pv 成本**：%.2f（考核 35）  \n\n",
			displayDate,
			int64(totalConvertDp),
			totalCashCost,
			totalConvertDpPrice,
			totalDpOrderPrice,
		)
		if regularCount > 0 || energyCount > 0 {
			markdownText += "---\n"
		}
	}

	// 常规活动数据
	if regularCount > 0 {
		markdownText += fmt.Sprintf(
			"**常规活动-海纳【飞猪常规活动 %s日报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 6000）  \n"+
				"**唤起成本**：%.2f（考核 0.5）  \n"+
				"**下单 pv 成本**：%.2f（考核 35）  \n\n",
			displayDate,
			int64(regularConvertDp),
			regularCashCost,
			regularConvertDpPrice,
			regularDpOrderPrice,
		)
		if energyCount > 0 {
			markdownText += "---\n"
		}
	}

	// 集能量活动数据
	if energyCount > 0 {
		markdownText += fmt.Sprintf(
			"**集能量-海纳【飞猪集能量 %s日报】**  \n"+
				"**唤起量**：%d  \n"+
				"**现金消耗**：%.2f（日预算 500）  \n"+
				"**唤起成本**：%.2f（考核 0.2）  \n",
			displayDate,
			int64(energyConvertDp),
			energyCashCost,
			energyConvertDpPrice,
		)
	}

	if markdownText == "" {
		logx.Info("没有需要发送的飞猪日报数据")
		return
	}

	// 检查钉钉配置
	if !dingConfig.Enabled || dingConfig.FzWebhookURL == "" {
		logx.Info("飞猪钉钉日报未启用或未配置webhook，跳过发送")
		return
	}

	// 发送到钉钉
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "飞猪外投日报",
			"text":  markdownText,
		},
	}

	client := httpclient.NewClient("", 30)
	var result map[string]interface{}
	err = client.Post(ctx, dingConfig.FzWebhookURL, msg, &result)
	if err != nil {
		logx.Errorf("发送飞猪钉钉日报失败: %v", err)
		return
	}

	logx.Infof("飞猪钉钉日报发送成功: %v", result)
}
