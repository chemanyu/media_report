package config

import (
	"context"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ==================== 获取服务费配置列表 ====================

type GetServiceFeeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServiceFeeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceFeeListLogic {
	return &GetServiceFeeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServiceFeeListLogic) GetServiceFeeList() (resp *types.ServiceFeeListResp, err error) {
	// 增加 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("[GetServiceFeeList] Panic recovered: %v", r)
			resp = &types.ServiceFeeListResp{
				Code:    500,
				Message: "服务器内部错误",
				Data:    []*types.ServiceFeeResp{},
			}
		}
	}()

	l.Logger.Info("[GetServiceFeeList] 开始查询服务费配置列表")

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Logger.Error("[GetServiceFeeList] 数据库连接为空")
		return &types.ServiceFeeListResp{
			Code:    500,
			Message: "数据库连接异常",
			Data:    []*types.ServiceFeeResp{},
		}, nil
	}

	fees, err := model.GetAllServiceFees(l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("[GetServiceFeeList] 获取服务费配置列表失败: %v", err)
		return &types.ServiceFeeListResp{
			Code:    500,
			Message: "获取失败: " + err.Error(),
			Data:    []*types.ServiceFeeResp{},
		}, nil
	}

	l.Logger.Infof("[GetServiceFeeList] 成功获取到 %d 条服务费配置记录", len(fees))

	var list []*types.ServiceFeeResp
	for idx, f := range fees {
		// 安全的时间格式化
		updateTimeStr := ""
		createTimeStr := ""

		if !f.UpdateTime.IsZero() {
			updateTimeStr = f.UpdateTime.Format(time.RFC3339)
		} else {
			l.Logger.Infof("[GetServiceFeeList] 记录[%d] ID=%d UpdateTime为零值", idx, f.ID)
		}

		if !f.CreateTime.IsZero() {
			createTimeStr = f.CreateTime.Format(time.RFC3339)
		} else {
			l.Logger.Infof("[GetServiceFeeList] 记录[%d] ID=%d CreateTime为零值", idx, f.ID)
		}

		list = append(list, &types.ServiceFeeResp{
			ID:              f.ID,
			ServiceProvider: f.ServiceProvider,
			FeeRate:         f.FeeRate,
			Remark:          f.Remark,
			UpdateTime:      updateTimeStr,
			CreateTime:      createTimeStr,
		})
	}

	l.Logger.Infof("[GetServiceFeeList] 成功构建响应，共 %d 条记录", len(list))

	return &types.ServiceFeeListResp{
		Code:    200,
		Message: "成功",
		Data:    list,
	}, nil
}

// ==================== 创建服务费配置 ====================

type CreateServiceFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateServiceFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateServiceFeeLogic {
	return &CreateServiceFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateServiceFeeLogic) CreateServiceFee(req *types.CreateServiceFeeReq) (resp *types.ServiceFeeCommonResp, err error) {
	fee := &model.ServiceFee{
		ServiceProvider: req.ServiceProvider,
		FeeRate:         req.FeeRate,
		Remark:          req.Remark,
	}

	if err := model.CreateServiceFee(l.svcCtx.DB, fee); err != nil {
		l.Logger.Errorf("创建服务费配置失败: %v", err)
		return &types.ServiceFeeCommonResp{
			Code:    500,
			Message: "创建失败: " + err.Error(),
		}, nil
	}

	return &types.ServiceFeeCommonResp{
		Code:    200,
		Message: "创建成功",
	}, nil
}

// ==================== 更新服务费配置 ====================

type UpdateServiceFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateServiceFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServiceFeeLogic {
	return &UpdateServiceFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateServiceFeeLogic) UpdateServiceFee(req *types.UpdateServiceFeeReq) (resp *types.ServiceFeeCommonResp, err error) {
	fee := &model.ServiceFee{
		ID:              req.ID,
		ServiceProvider: req.ServiceProvider,
		FeeRate:         req.FeeRate,
		Remark:          req.Remark,
	}

	if err := model.UpdateServiceFee(l.svcCtx.DB, fee); err != nil {
		l.Logger.Errorf("更新服务费配置失败: %v", err)
		return &types.ServiceFeeCommonResp{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		}, nil
	}

	return &types.ServiceFeeCommonResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// ==================== 删除服务费配置 ====================

type DeleteServiceFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteServiceFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServiceFeeLogic {
	return &DeleteServiceFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteServiceFeeLogic) DeleteServiceFee(req *types.DeleteServiceFeeReq) (resp *types.ServiceFeeCommonResp, err error) {
	if err := model.DeleteServiceFee(l.svcCtx.DB, req.ID); err != nil {
		l.Logger.Errorf("删除服务费配置失败: %v", err)
		return &types.ServiceFeeCommonResp{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		}, nil
	}

	return &types.ServiceFeeCommonResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}
