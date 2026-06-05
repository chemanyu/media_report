package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FzConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFzConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FzConfigLogic {
	return &FzConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FzConfigLogic) GetConfig() (*types.FzConfigResp, error) {
	config, err := model.GetFzConfig(l.svcCtx.DB)
	if err != nil {
		return nil, err
	}
	return &types.FzConfigResp{
		Id:          config.Id,
		Coefficient: config.Coefficient,
		BaseNum:     config.BaseNum,
		UpdateTime:  config.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

func (l *FzConfigLogic) UpdateConfig(req *types.FzUpdateConfigReq) error {
	return model.UpdateFzConfig(l.svcCtx.DB, req.Coefficient, req.BaseNum)
}
