package config

import (
	"context"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ==================== 获取返点配置列表 ====================

type RebateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRebateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RebateLogic) GetRebateList() (resp *types.RebateListResp, err error) {
	rebates, err := model.GetAllRebates(l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("获取返点配置列表失败: %v", err)
		return &types.RebateListResp{
			Code:    500,
			Message: "获取失败: " + err.Error(),
			Data:    []*types.RebateResp{},
		}, nil
	}

	var list []*types.RebateResp
	for _, r := range rebates {
		list = append(list, &types.RebateResp{
			ID:          r.ID,
			Subject:     r.Subject,
			Port:        r.Port,
			RebateRate:  r.RebateRate,
			SubjectType: r.SubjectType,
			Remark:      r.Remark,
			UpdateTime:  r.UpdateTime.Format(time.RFC3339),
			CreateTime:  r.CreateTime.Format(time.RFC3339),
		})
	}

	return &types.RebateListResp{
		Code:    200,
		Message: "成功",
		Data:    list,
	}, nil
}

// ==================== 创建返点配置 ====================
func NewCreateRebateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RebateLogic) CreateRebate(req *types.CreateRebateReq) (resp *types.RebateCommonResp, err error) {
	rebate := &model.Rebate{
		Subject:     req.Subject,
		Port:        req.Port,
		RebateRate:  req.RebateRate,
		SubjectType: req.SubjectType,
		Remark:      req.Remark,
	}

	if err := model.CreateRebate(l.svcCtx.DB, rebate); err != nil {
		l.Logger.Errorf("创建返点配置失败: %v", err)
		return &types.RebateCommonResp{
			Code:    500,
			Message: "创建失败: " + err.Error(),
		}, nil
	}

	return &types.RebateCommonResp{
		Code:    200,
		Message: "创建成功",
	}, nil
}

// ==================== 更新返点配置 ====================

func NewUpdateRebateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RebateLogic) UpdateRebate(req *types.UpdateRebateReq) (resp *types.RebateCommonResp, err error) {
	rebate := &model.Rebate{
		ID:          req.ID,
		Subject:     req.Subject,
		Port:        req.Port,
		RebateRate:  req.RebateRate,
		SubjectType: req.SubjectType,
		Remark:      req.Remark,
	}

	if err := model.UpdateRebate(l.svcCtx.DB, rebate); err != nil {
		l.Logger.Errorf("更新返点配置失败: %v", err)
		return &types.RebateCommonResp{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		}, nil
	}

	return &types.RebateCommonResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// ==================== 删除返点配置 ====================
func NewDeleteRebateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RebateLogic) DeleteRebate(req *types.DeleteRebateReq) (resp *types.RebateCommonResp, err error) {
	if err := model.DeleteRebate(l.svcCtx.DB, req.ID); err != nil {
		l.Logger.Errorf("删除返点配置失败: %v", err)
		return &types.RebateCommonResp{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		}, nil
	}

	return &types.RebateCommonResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}
