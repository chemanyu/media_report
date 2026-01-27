package config

import (
	"context"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ==================== 获取任务类型列表 ====================

type GetTaskTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskTypeListLogic {
	return &GetTaskTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskTypeListLogic) GetTaskTypeList() (resp *types.TaskTypeListResp, err error) {
	taskTypes, err := model.GetAllTaskTypes(l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("获取任务类型列表失败: %v", err)
		return &types.TaskTypeListResp{
			Code:    500,
			Message: "获取失败: " + err.Error(),
			Data:    []*types.TaskTypeResp{},
		}, nil
	}

	var list []*types.TaskTypeResp
	for _, t := range taskTypes {
		list = append(list, &types.TaskTypeResp{
			ID:              t.ID,
			Name:            t.Name,
			Code:            t.Code,
			SettlementPrice: t.SettlementPrice,
			Media:           t.Media,
			Status:          t.Status,
			CreateTime:      t.CreateTime.Format(time.RFC3339),
			UpdateTime:      t.UpdateTime.Format(time.RFC3339),
		})
	}

	return &types.TaskTypeListResp{
		Code:    200,
		Message: "成功",
		Data:    list,
	}, nil
}

// ==================== 创建任务类型 ====================

type CreateTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskTypeLogic {
	return &CreateTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskTypeLogic) CreateTaskType(req *types.CreateTaskTypeReq) (resp *types.TaskTypeCommonResp, err error) {
	taskType := &model.TaskType{
		Name:            req.Name,
		Code:            req.Code,
		SettlementPrice: req.SettlementPrice,
		Media:           req.Media,
		Status:          req.Status,
	}

	if err := model.CreateTaskType(l.svcCtx.DB, taskType); err != nil {
		l.Logger.Errorf("创建任务类型失败: %v", err)
		return &types.TaskTypeCommonResp{
			Code:    500,
			Message: "创建失败: " + err.Error(),
		}, nil
	}

	return &types.TaskTypeCommonResp{
		Code:    200,
		Message: "创建成功",
	}, nil
}

// ==================== 更新任务类型 ====================

type UpdateTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskTypeLogic {
	return &UpdateTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskTypeLogic) UpdateTaskType(req *types.UpdateTaskTypeReq) (resp *types.TaskTypeCommonResp, err error) {
	taskType := &model.TaskType{
		ID:              req.ID,
		Name:            req.Name,
		Code:            req.Code,
		SettlementPrice: req.SettlementPrice,
		Media:           req.Media,
		Status:          req.Status,
	}

	if err := model.UpdateTaskType(l.svcCtx.DB, taskType); err != nil {
		l.Logger.Errorf("更新任务类型失败: %v", err)
		return &types.TaskTypeCommonResp{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		}, nil
	}

	return &types.TaskTypeCommonResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// ==================== 删除任务类型 ====================

type DeleteTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskTypeLogic {
	return &DeleteTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskTypeLogic) DeleteTaskType(req *types.DeleteTaskTypeReq) (resp *types.TaskTypeCommonResp, err error) {
	if err := model.DeleteTaskType(l.svcCtx.DB, req.ID); err != nil {
		l.Logger.Errorf("删除任务类型失败: %v", err)
		return &types.TaskTypeCommonResp{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		}, nil
	}

	return &types.TaskTypeCommonResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}
