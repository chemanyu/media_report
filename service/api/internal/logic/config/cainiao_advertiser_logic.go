package config

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// =============== 获取列表 ===============

type CainiaoAdvertiserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCainiaoAdvertiserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CainiaoAdvertiserListLogic {
	return &CainiaoAdvertiserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CainiaoAdvertiserListLogic) GetList() (resp []model.CainiaoAdvertiser, err error) {
	return model.ListAdvertisers(l.svcCtx.DB)
}

// =============== 添加账户 ===============

type CainiaoAdvertiserAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCainiaoAdvertiserAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CainiaoAdvertiserAddLogic {
	return &CainiaoAdvertiserAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CainiaoAdvertiserAddLogic) Add(req *types.CainiaoAdvertiserAddReq) (resp *types.CommonResp, err error) {
	err = model.AddAdvertiser(l.svcCtx.DB, req.MediaAdvId)
	if err != nil {
		return &types.CommonResp{
			Code:    500,
			Message: "添加失败: " + err.Error(),
		}, nil
	}

	return &types.CommonResp{
		Code:    0,
		Message: "添加成功",
	}, nil
}

// =============== 删除账户 ===============

type CainiaoAdvertiserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCainiaoAdvertiserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CainiaoAdvertiserDeleteLogic {
	return &CainiaoAdvertiserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CainiaoAdvertiserDeleteLogic) Delete(req *types.CainiaoAdvertiserDeleteReq) (resp *types.CommonResp, err error) {
	err = model.DeleteAdvertiser(l.svcCtx.DB, req.Id)
	if err != nil {
		return &types.CommonResp{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		}, nil
	}

	return &types.CommonResp{
		Code:    0,
		Message: "删除成功",
	}, nil
}
