package qczj

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QczjConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQczjConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QczjConfigLogic {
	return &QczjConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QczjConfigLogic) GetConfig() (*types.QczjConfigResp, error) {
	config, err := model.GetQczjConfig(l.svcCtx.DB)
	if err != nil {
		return nil, err
	}
	return &types.QczjConfigResp{
		Id:         config.Id,
		TotalUv:    config.TotalUv,
		Ratio:      config.Ratio,
		UpdateTime: config.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

func (l *QczjConfigLogic) UpdateConfig(req *types.QczjUpdateConfigReq) error {
	return model.UpdateQczjConfig(l.svcCtx.DB, req.TotalUv, req.Ratio)
}

func (l *QczjConfigLogic) GetReportList() (*types.QczjReportListResp, error) {
	list, total, err := model.ListQczjReportData(l.svcCtx.DB)
	if err != nil {
		return nil, err
	}

	items := make([]*types.QczjReportItem, 0, len(list))
	for _, r := range list {
		items = append(items, &types.QczjReportItem{
			Id:         r.Id,
			ReportDate: r.ReportDate,
			View:       r.View,
			Click:      r.Click,
			Action:     r.Action,
			UpdateTime: r.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.QczjReportListResp{
		Total: total,
		List:  items,
	}, nil
}
