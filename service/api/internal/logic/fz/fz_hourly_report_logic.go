package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"media_report/service/api/internal/honor"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/oppo"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
	"media_report/service/api/internal/xiaomi"
)

// FzHourlyReportLogic 飞猪时报业务逻辑
type FzHourlyReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFzHourlyReportLogic 创建飞猪时报逻辑实例
func NewFzHourlyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FzHourlyReportLogic {
	return &FzHourlyReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SyncOppoData 同步OPPO媒体数据
func (l *FzHourlyReportLogic) SyncOppoData(reportDate string) (int, error) {
	// 1. 从数据库获取所有OPPO账户
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)
	advertisers, err := advertiserModel.FindByMedia("oppo")
	if err != nil {
		return 0, fmt.Errorf("查询OPPO账户失败: %w", err)
	}

	if len(advertisers) == 0 {
		return 0, fmt.Errorf("未找到OPPO账户")
	}

	// 2. 创建OPPO API客户端（从配置中获取API密钥）
	oppoClient := oppo.NewOppoAPIClient(
		l.svcCtx.Config.OppoAPI.OwnerId,
		l.svcCtx.Config.OppoAPI.ApiId,
		l.svcCtx.Config.OppoAPI.ApiKey,
	)

	// 3. 遍历账户，调用API获取数据
	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)
	successCount := 0

	for _, advertiser := range advertisers {
		// 将媒体账户ID转换为int64
		ownerId, err := strconv.ParseInt(advertiser.MediaAdvId, 10, 64)
		if err != nil {
			fmt.Printf("账户ID转换失败: %s, err: %v\n", advertiser.MediaAdvId, err)
			continue
		}

		// 构建请求参数
		req := oppo.AdDataRequest{
			BeginTime: reportDate,
			EndTime:   reportDate,
			TimeLevel: "DAY",
			OwnerId:   ownerId,
			ParaMap: map[string]interface{}{
				"filter_zero": 0,
			},
		}

		// 调用API
		adData, err := oppoClient.QueryAdData(req)
		if err != nil {
			fmt.Printf("查询账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		// 将日期字符串转换为int（例如: "20260211" -> 20260211）
		reportDateInt, _ := strconv.Atoi(reportDate)

		// 构建报表数据
		report := &model.FzHourlyReport{
			Media:           "oppo",
			MediaAdvId:      advertiser.MediaAdvId,
			MediaAdvName:    advertiser.MediaAdvName,
			ReportDate:      reportDateInt,
			Cost:            adData.Cost,
			ConvertDp:       adData.ConvertDp,
			DpAppOrderNums:  adData.DpAppOrderNums,
			Click:           adData.Click,
			Expose:          adData.Expose,
			ConvertDpPrice:  adData.ConvertDpPrice,
			DpAppOrderPrice: adData.DpAppOrderPrice,
		}

		// 保存到数据库（插入或更新）
		if err := reportModel.InsertOrUpdate(report); err != nil {
			fmt.Printf("保存账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		successCount++
		fmt.Printf("成功同步账户 %s(%s) 数据: 消耗=%.2f, 拉活=%d, 订单=%d\n",
			advertiser.MediaAdvName, advertiser.MediaAdvId, adData.Cost, adData.ConvertDp, adData.DpAppOrderNums)
	}

	return successCount, nil
}

// SyncTodayOppoData 同步今天的OPPO数据
func (l *FzHourlyReportLogic) SyncTodayOppoData() (int, error) {
	// 使用北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Now().In(loc).Format("20060102")

	return l.SyncOppoData(today)
}

// SyncXiaomiData 同步小米媒体数据
func (l *FzHourlyReportLogic) SyncXiaomiData(reportDate string) (int, error) {
	// 1. 从数据库获取所有小米账户
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)
	advertisers, err := advertiserModel.FindByMedia("xiaomi")
	if err != nil {
		return 0, fmt.Errorf("查询小米账户失败: %w", err)
	}

	if len(advertisers) == 0 {
		return 0, fmt.Errorf("未找到小米账户")
	}

	// 2. 创建小米API客户端（从配置中获取API密钥）
	xiaomiClient := xiaomi.NewXiaomiAPIClient(
		l.svcCtx.Config.XiaomiAPI.SignId,
		l.svcCtx.Config.XiaomiAPI.SecretKey,
		0, // CustomerId 会从账户列表中获取
	)

	// 3. 遍历账户，调用API获取数据
	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)
	successCount := 0

	// 将日期格式从 20260211 转换为 2026-02-11
	var formattedDate string
	if len(reportDate) == 8 {
		formattedDate = fmt.Sprintf("%s-%s-%s", reportDate[0:4], reportDate[4:6], reportDate[6:8])
	} else {
		formattedDate = reportDate
	}

	for _, advertiser := range advertisers {
		// 将媒体账户ID转换为int64
		customerId, err := strconv.ParseInt(advertiser.MediaAdvId, 10, 64)
		if err != nil {
			fmt.Printf("账户ID转换失败: %s, err: %v\n", advertiser.MediaAdvId, err)
			continue
		}

		// 更新客户端的customerId
		xiaomiClient.CustomerId = customerId

		// 构建请求参数 - 使用指标code: 2017(消耗), 2012(下载量)等
		req := xiaomi.ReportDataRequest{
			CustomerId:  customerId,
			SDate:       formattedDate,
			EDate:       formattedDate,
			MetricsList: "2017,2012,1863,2018", // 消耗、下载量、点击量、曝光量
			Page:        1,
			PageSize:    1000,
		}

		// 调用API
		result, err := xiaomiClient.GetReportData(&req)
		if err != nil {
			fmt.Printf("查询账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		// 检查API调用是否成功
		if !result.Success || result.Code != 0 {
			fmt.Printf("查询账户 %s(%s) API返回错误: success=%v, code=%d\n", advertiser.MediaAdvName, advertiser.MediaAdvId, result.Success, result.Code)
			continue
		}

		// 汇总数据 - 从 result.Result.Total 获取
		totalData := result.Result.Total

		// 将字符串类型的 cost 转换为 float64
		totalCost, _ := strconv.ParseFloat(totalData.Cost, 64)

		// reActiveSumjFormat 是激活数（对应拉活数）
		convertDp, _ := strconv.ParseInt(totalData.ReActiveSumjFormat, 10, 64)

		dpAppOrderNums, _ := strconv.ParseInt(totalData.PaySumjFormat, 10, 64)

		// 将日期字符串转换为int（例如: "20260211" -> 20260211）
		reportDateInt, _ := strconv.Atoi(reportDate)

		// 计算成本指标
		convertDpPrice, _ := strconv.ParseFloat(totalData.CostPerReActivejFormat, 64)
		dpAppOrderPrice, _ := strconv.ParseFloat(totalData.CostPerPayzFormat, 64)

		// 构建报表数据
		report := &model.FzHourlyReport{
			Media:           "xiaomi",
			MediaAdvId:      advertiser.MediaAdvId,
			MediaAdvName:    advertiser.MediaAdvName,
			ReportDate:      reportDateInt,
			Cost:            totalCost * 100,
			ConvertDp:       convertDp,
			DpAppOrderNums:  dpAppOrderNums,
			Click:           totalData.ClickNum,
			Expose:          totalData.ExposeNum,
			ConvertDpPrice:  convertDpPrice * 100,
			DpAppOrderPrice: dpAppOrderPrice * 100,
		}

		// 保存到数据库（插入或更新）
		if err := reportModel.InsertOrUpdate(report); err != nil {
			fmt.Printf("保存账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		successCount++
		fmt.Printf("成功同步账户 %s(%s) 数据: 消耗=%.2f, 激活=%d, 下载=%d, 点击=%d, 曝光=%d\n",
			advertiser.MediaAdvName, advertiser.MediaAdvId, totalCost, convertDp, totalData.DownloadNum, totalData.ClickNum, totalData.ExposeNum)
	}

	return successCount, nil
}

// SyncTodayXiaomiData 同步今天的小米数据
func (l *FzHourlyReportLogic) SyncTodayXiaomiData() (int, error) {
	// 使用北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Now().In(loc).Format("20060102")

	return l.SyncXiaomiData(today)
}

// SaveAdnData 保存ADN媒体数据
func (l *FzHourlyReportLogic) SaveAdnData(req *types.FzSyncAdnDataReq) error {
	// 将日期字符串转换为int（例如: "20260211" -> 20260211）
	reportDateInt, err := strconv.Atoi(req.ReportDate)
	if err != nil {
		return fmt.Errorf("报表日期格式错误: %w", err)
	}

	// 构建报表数据
	report := &model.FzHourlyReport{
		Media:           "adn",
		MediaAdvId:      req.MediaAdvId,
		MediaAdvName:    req.MediaAdvName,
		ReportDate:      reportDateInt,
		Cost:            req.Cost * 100,
		ConvertDp:       req.ConvertDp,
		DpAppOrderNums:  req.DpAppOrderNums,
		Click:           req.Click,
		Expose:          req.Expose,
		ConvertDpPrice:  req.ConvertDpPrice * 100,
		DpAppOrderPrice: req.DpAppOrderPrice * 100,
	}

	// 保存到数据库（插入或更新）
	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)
	if err := reportModel.InsertOrUpdate(report); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}

	fmt.Printf("成功保存ADN账户 %s(%s) 数据: 消耗=%.2f, 拉活=%d, 订单=%d, 点击=%d, 曝光=%d\n",
		req.MediaAdvName, req.MediaAdvId, req.Cost, req.ConvertDp, req.DpAppOrderNums, req.Click, req.Expose)

	return nil
}

// SyncHonorData 同步荣耀媒体数据
func (l *FzHourlyReportLogic) SyncHonorData(reportDate string) (int, error) {
	// 1. 从数据库获取所有荣耀账户
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)
	advertisers, err := advertiserModel.FindByMedia("honor")
	if err != nil {
		return 0, fmt.Errorf("查询荣耀账户失败: %w", err)
	}

	if len(advertisers) == 0 {
		return 0, fmt.Errorf("未找到荣耀账户")
	}

	// 2. 将日期格式从 20260211 转换为 2026-02-11
	var formattedDate string
	if len(reportDate) == 8 {
		formattedDate = fmt.Sprintf("%s-%s-%s", reportDate[0:4], reportDate[4:6], reportDate[6:8])
	} else {
		formattedDate = reportDate
	}

	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)
	reportDateInt, _ := strconv.Atoi(reportDate)
	successCount := 0

	for _, advertiser := range advertisers {
		if advertiser.ClientID == "" || advertiser.ClientSecret == "" {
			fmt.Printf("荣耀账户 %s(%s) 未配置 ClientID/ClientSecret，跳过\n", advertiser.MediaAdvName, advertiser.MediaAdvId)
			continue
		}

		// 3. 每个账户使用自己的凭据创建独立客户端
		honorClient := honor.NewHonorAPIClient(advertiser.ClientID, advertiser.ClientSecret)

		// 4. 查询该账户广告主报表（honorPull/payment 非默认字段，须显式指定）
		req := honor.ReportRequest{
			StartTime:       formattedDate,
			EndTime:         formattedDate,
			TimeDimension:   0,
			PageIndex:       1,
			PageSize:        100,
			IndexScreenList: []string{"honorPull", "honorPullCost", "payment", "paymentCost"},
		}

		items, err := honorClient.QueryAdvertiserReport(req, advertiser.MediaAdvId)
		if err != nil {
			fmt.Printf("查询荣耀账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		if len(items) == 0 {
			fmt.Printf("荣耀账户 %s(%s) 暂无数据\n", advertiser.MediaAdvName, advertiser.MediaAdvId)
			continue
		}

		// 5. 汇总所有返回条目（API 按版位拆分，需汇总）
		var totalCost float64
		var totalImpression, totalClick, totalHonorPull, totalPayment int64
		var totalHonorPullCost, totalPaymentCost float64

		for _, item := range items {
			totalCost += honor.ToFloat64(item.Metrics.AdBilling)
			totalImpression += honor.ToInt64(item.Metrics.Impression)
			totalClick += honor.ToInt64(item.Metrics.Click)
			totalHonorPull += honor.ToInt64(item.Metrics.HonorPull)
			totalHonorPullCost += honor.ToFloat64(item.Metrics.HonorPullCost)
			totalPayment += honor.ToInt64(item.Metrics.Payment)
			totalPaymentCost += honor.ToFloat64(item.Metrics.PaymentCost)
		}

		var avgHonorPullCost float64
		if totalHonorPull > 0 {
			avgHonorPullCost = totalHonorPullCost / float64(totalHonorPull)
		}
		var avgPaymentCost float64
		if totalPayment > 0 {
			avgPaymentCost = totalPaymentCost / float64(totalPayment)
		}

		report := &model.FzHourlyReport{
			Media:           "honor",
			MediaAdvId:      advertiser.MediaAdvId,
			MediaAdvName:    advertiser.MediaAdvName,
			ReportDate:      reportDateInt,
			Cost:            totalCost / 10000,             // 微→分（1元=1000000微，1分=10000微）
			ConvertDp:       totalHonorPull,                // 全网首唤数
			DpAppOrderNums:  totalPayment,                  // 付费数
			Click:           totalClick,
			Expose:          totalImpression,
			ConvertDpPrice:  avgHonorPullCost / 10000,      // 微→分
			DpAppOrderPrice: avgPaymentCost / 10000,        // 微→分
		}

		if err := reportModel.InsertOrUpdate(report); err != nil {
			fmt.Printf("保存荣耀账户 %s(%s) 数据失败: %v\n", advertiser.MediaAdvName, advertiser.MediaAdvId, err)
			continue
		}

		successCount++
		fmt.Printf("成功同步荣耀账户 %s(%s) 数据: 消耗=%.2f分, 全网首唤=%d, 付费=%d, 点击=%d, 曝光=%d\n",
			advertiser.MediaAdvName, advertiser.MediaAdvId, totalCost, totalHonorPull, totalPayment, totalClick, totalImpression)
	}

	return successCount, nil
}

// SyncTodayHonorData 同步今天的荣耀数据
func (l *FzHourlyReportLogic) SyncTodayHonorData() (int, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Now().In(loc).Format("20060102")
	return l.SyncHonorData(today)
}

// GetReportList 获取报表列表（最多查询单天数据，以 end_date 为准）
func (l *FzHourlyReportLogic) GetReportList(req *types.FzHourlyReportListReq) ([]*model.FzHourlyReport, error) {
	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)

	// 以 end_date 为查询日期，不传则取今天
	var queryDate int
	if req.EndDate != "" {
		queryDate, _ = strconv.Atoi(req.EndDate)
	} else if req.StartDate != "" {
		queryDate, _ = strconv.Atoi(req.StartDate)
	} else {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		queryDate, _ = strconv.Atoi(time.Now().In(loc).Format("20060102"))
	}

	return reportModel.FindByDateRange(req.Media, queryDate, queryDate)
}
