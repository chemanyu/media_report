package logic

import (
	"context"
	"fmt"
	"math"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetKsAccountReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKsAccountReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetKsAccountReportLogic {
	return &GetKsAccountReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKsAccountReportLogic) GetKsAccountReport(req *types.KsAccountReportReq) (resp *types.KsAccountReportResp, err error) {
	// 参数验证
	if req.StartDate == "" || req.EndDate == "" {
		return &types.KsAccountReportResp{
			Code:    400,
			Message: "start_date and end_date are required",
			Data:    []*types.KsReportDataItem{},
		}, nil
	}

	if req.AdvertiserId <= 0 {
		return &types.KsAccountReportResp{
			Code:    400,
			Message: "advertiser_id is required and must be positive",
			Data:    []*types.KsReportDataItem{},
		}, nil
	}

	// 设置默认值
	if req.TemporalGranularity == "" {
		req.TemporalGranularity = "DAILY"
	}

	// 构建快手 API 请求参数
	ksReq := map[string]interface{}{
		"start_date":           req.StartDate,
		"end_date":             req.EndDate,
		"advertiser_id":        req.AdvertiserId,
		"temporal_granularity": req.TemporalGranularity,
	}

	// 从数据库获取当前有效的 access token
	accessToken, _, err := model.GetTokensByMedia(l.svcCtx.DB, "kuaishou")
	if err != nil {
		return &types.KsAccountReportResp{
			Code:    500,
			Message: "获取accessToken失败：: " + err.Error(),
			Data:    []*types.KsReportDataItem{},
		}, nil
	}

	// 设置快手 API 需要的请求头
	l.svcCtx.HTTPClient.SetHeader("Access-Token", accessToken)
	l.svcCtx.HTTPClient.SetHeader("Content-Type", "application/json")

	// 调用快手 API
	var ksResp types.KsApiResponse
	err = l.svcCtx.HTTPClient.Post(l.ctx, "/rest/openapi/v1/report/account_report", ksReq, &ksResp)
	if err != nil {
		logx.Errorf("call kuaishou api failed: %v", err)
		return &types.KsAccountReportResp{
			Code:    500,
			Message: "call kuaishou api failed: " + err.Error(),
			Data:    []*types.KsReportDataItem{},
		}, nil
	}

	// 检查快手 API 返回的业务状态码
	if ksResp.Code != 0 {
		logx.Errorf("kuaishou api returned error code: %d, message: %s", ksResp.Code, ksResp.Message)
		return &types.KsAccountReportResp{
			Code:    500,
			Message: fmt.Sprintf("kuaishou api error: %s", ksResp.Message),
			Data:    []*types.KsReportDataItem{},
		}, nil
	}

	// 转换数据并计算
	dataItems := make([]*types.KsReportDataItem, 0, len(ksResp.Data.Details))
	currentHour := time.Now().Format("15") // 获取当前小时
	for _, detail := range ksResp.Data.Details {
		// 计算转化率并格式化为百分比字符串
		conversionRatio := fmt.Sprintf("%.2f%%", detail.ConversionRatio*100)

		// 计算消耗和转化成本，保留两位小数
		charge := math.Round(detail.Charge*1*100) / 100
		conversionCost := math.Round(detail.ConversionCost*1*100) / 100

		// 在日期后添加小时
		timeWithHour := fmt.Sprintf("%s %s", detail.StatDate, currentHour)

		dataItems = append(dataItems, &types.KsReportDataItem{
			Time:            timeWithHour,         // 统计日期 + 小时
			Account:         "美致dsp",              // 账户名称（可配置）
			Charge:          charge,               // 消耗 * 1，保留两位小数
			Activation:      detail.Activation,    // 注册转化数（激活数）
			ConversionCost:  conversionCost,       // 转化成本 * 1，保留两位小数
			AdShow:          int64(detail.AdShow), // 曝光数
			Bclick:          detail.Bclick,        // 点击数
			ConversionRatio: conversionRatio,      // 转化率（格式化）
		})
	}

	return &types.KsAccountReportResp{
		Code:    0,
		Message: "success",
		Data:    dataItems,
	}, nil
}
