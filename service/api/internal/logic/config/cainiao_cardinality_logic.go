package config

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// =============== 获取基数配置 ===============

type CainiaoCardinalityGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCainiaoCardinalityGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CainiaoCardinalityGetLogic {
	return &CainiaoCardinalityGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CainiaoCardinalityGetLogic) Get() (resp *types.CainiaoCardinalityResp, err error) {
	config, err := model.GetCardinality(l.svcCtx.DB)
	if err != nil {
		return &types.CainiaoCardinalityResp{
			Code:    500,
			Message: "获取失败: " + err.Error(),
		}, nil
	}

	return &types.CainiaoCardinalityResp{
		Code:    0,
		Message: "success",
		Data: &types.CardinalityData{
			ID:          config.ID,
			Cardinality: config.Cardinality,
		},
	}, nil
}

// =============== 更新基数配置 ===============

type CainiaoCardinalityUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCainiaoCardinalityUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CainiaoCardinalityUpdateLogic {
	return &CainiaoCardinalityUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CainiaoCardinalityUpdateLogic) Update(req *types.CainiaoCardinalityUpdateReq) (resp *types.CommonResp, err error) {
	// 验证基数值范围（可选）
	if req.Cardinality <= 0 || req.Cardinality > 99.9 {
		return &types.CommonResp{
			Code:    400,
			Message: "基数值必须在0.1到99.9之间",
		}, nil
	}

	err = model.UpdateCardinality(l.svcCtx.DB, req.Cardinality)
	if err != nil {
		return &types.CommonResp{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		}, nil
	}

	return &types.CommonResp{
		Code:    0,
		Message: "更新成功",
	}, nil
}
