package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
			Cost:            totalCost,
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

// SyncYesterdayXiaomiData 同步昨天的小米数据
func (l *FzHourlyReportLogic) SyncYesterdayXiaomiData() (int, error) {
	// 使用北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1).Format("20060102")

	return l.SyncXiaomiData(yesterday)
}

// GetReportList 获取报表列表
func (l *FzHourlyReportLogic) GetReportList(req *types.FzHourlyReportListReq) ([]*model.FzHourlyReport, error) {
	reportModel := model.NewFzHourlyReportModel(l.svcCtx.DB)

	// 解析日期参数
	var startDate, endDate int
	if req.StartDate != "" {
		startDate, _ = strconv.Atoi(req.StartDate)
	}
	if req.EndDate != "" {
		endDate, _ = strconv.Atoi(req.EndDate)
	}

	// 查询数据
	return reportModel.FindByDateRange(req.Media, startDate, endDate)
}
