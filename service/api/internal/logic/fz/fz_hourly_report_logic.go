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

// SyncYesterdayOppoData 同步昨天的OPPO数据
func (l *FzHourlyReportLogic) SyncYesterdayOppoData() (int, error) {
	// 使用北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1).Format("20060102")

	return l.SyncOppoData(yesterday)
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
