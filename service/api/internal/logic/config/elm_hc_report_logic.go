package config

import (
	"context"
	"fmt"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetElmHcReportListLogic 获取汇川饿了么数据报表列表
type GetElmHcReportListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetElmHcReportListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetElmHcReportListLogic {
	return &GetElmHcReportListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetElmHcReportListLogic) GetElmHcReportList() (*types.ElmHcReportListResp, error) {
	reports, err := model.GetAllElmHcPerformanceReports(l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("查询汇川饿了么数据报表列表失败: %v", err)
		return &types.ElmHcReportListResp{
			Code:    500,
			Message: "查询失败",
		}, nil
	}

	var respList []*types.ElmHcReportResp
	for _, report := range reports {
		respList = append(respList, &types.ElmHcReportResp{
			ID:                report.ID,
			CustomerName:      report.CustomerName,
			CustomerShort:     report.CustomerShort,
			AgentName:         report.AgentName,
			AgentShort:        report.AgentShort,
			MediaPlatformName: report.MediaPlatformName,
			CreateTime:        report.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:        report.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ElmHcReportListResp{
		Code:    0,
		Message: "success",
		Data:    respList,
	}, nil
}

// CreateElmHcReportLogic 创建汇川饿了么数据报表
type CreateElmHcReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateElmHcReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateElmHcReportLogic {
	return &CreateElmHcReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateElmHcReportLogic) CreateElmHcReport(req *types.CreateElmHcReportReq) (*types.ElmHcReportCommonResp, error) {
	// 参数校验
	if req.CustomerName == "" || req.MediaPlatformName == "" {
		return &types.ElmHcReportCommonResp{
			Code:    400,
			Message: "客户名称和媒体平台名称不能为空",
		}, nil
	}

	report := &model.ElmHcPerformanceReport{
		CustomerName:      req.CustomerName,
		CustomerShort:     req.CustomerShort,
		AgentName:         req.AgentName,
		AgentShort:        req.AgentShort,
		MediaPlatformName: req.MediaPlatformName,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}

	if err := model.CreateElmHcPerformanceReport(l.svcCtx.DB, report); err != nil {
		l.Logger.Errorf("创建汇川饿了么数据报表失败: %v", err)
		return &types.ElmHcReportCommonResp{
			Code:    500,
			Message: fmt.Sprintf("创建失败: %v", err),
		}, nil
	}

	return &types.ElmHcReportCommonResp{
		Code:    0,
		Message: "创建成功",
	}, nil
}

// UpdateElmHcReportLogic 更新汇川饿了么数据报表
type UpdateElmHcReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateElmHcReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateElmHcReportLogic {
	return &UpdateElmHcReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateElmHcReportLogic) UpdateElmHcReport(req *types.UpdateElmHcReportReq) (*types.ElmHcReportCommonResp, error) {
	// 参数校验
	if req.ID <= 0 {
		return &types.ElmHcReportCommonResp{
			Code:    400,
			Message: "ID不能为空",
		}, nil
	}

	if req.CustomerName == "" || req.MediaPlatformName == "" {
		return &types.ElmHcReportCommonResp{
			Code:    400,
			Message: "客户名称和媒体平台名称不能为空",
		}, nil
	}

	report := &model.ElmHcPerformanceReport{
		ID:                req.ID,
		CustomerName:      req.CustomerName,
		CustomerShort:     req.CustomerShort,
		AgentName:         req.AgentName,
		AgentShort:        req.AgentShort,
		MediaPlatformName: req.MediaPlatformName,
	}

	if err := model.UpdateElmHcPerformanceReport(l.svcCtx.DB, report); err != nil {
		l.Logger.Errorf("更新汇川饿了么数据报表失败: %v", err)
		return &types.ElmHcReportCommonResp{
			Code:    500,
			Message: fmt.Sprintf("更新失败: %v", err),
		}, nil
	}

	return &types.ElmHcReportCommonResp{
		Code:    0,
		Message: "更新成功",
	}, nil
}

// DeleteElmHcReportLogic 删除汇川饿了么数据报表
type DeleteElmHcReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteElmHcReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteElmHcReportLogic {
	return &DeleteElmHcReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteElmHcReportLogic) DeleteElmHcReport(req *types.DeleteElmHcReportReq) (*types.ElmHcReportCommonResp, error) {
	if req.ID <= 0 {
		return &types.ElmHcReportCommonResp{
			Code:    400,
			Message: "ID不能为空",
		}, nil
	}

	if err := model.DeleteElmHcPerformanceReport(l.svcCtx.DB, req.ID); err != nil {
		l.Logger.Errorf("删除汇川饿了么数据报表失败: %v", err)
		return &types.ElmHcReportCommonResp{
			Code:    500,
			Message: fmt.Sprintf("删除失败: %v", err),
		}, nil
	}

	return &types.ElmHcReportCommonResp{
		Code:    0,
		Message: "删除成功",
	}, nil
}
