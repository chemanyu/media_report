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

// GetElmHcMediaListLogic 获取汇川饿了么媒体账户列表
type GetElmHcMediaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetElmHcMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetElmHcMediaListLogic {
	return &GetElmHcMediaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetElmHcMediaListLogic) GetElmHcMediaList() (*types.ElmHcMediaListResp, error) {
	medias, err := model.GetAllElmHcMediaReports(l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("查询汇川饿了么媒体账户列表失败: %v", err)
		return &types.ElmHcMediaListResp{
			Code:    500,
			Message: "查询失败",
		}, nil
	}

	var respList []*types.ElmHcMediaResp
	for _, media := range medias {
		respList = append(respList, &types.ElmHcMediaResp{
			ID:            media.ID,
			PerformanceID: media.PerformanceID,
			MediaAdvId:    media.MediaAdvId,
			MediaAdvName:  media.MediaAdvName,
			HuichuanAdvId: media.HuichuanAdvId,
			RedirectNum:   media.RedirectNum,
			PayNum:        media.PayNum,
			CreateTime:    media.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:    media.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ElmHcMediaListResp{
		Code:    0,
		Message: "success",
		Data:    respList,
	}, nil
}

// CreateElmHcMediaLogic 创建汇川饿了么媒体账户
type CreateElmHcMediaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateElmHcMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateElmHcMediaLogic {
	return &CreateElmHcMediaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateElmHcMediaLogic) CreateElmHcMedia(req *types.CreateElmHcMediaReq) (*types.ElmHcMediaCommonResp, error) {
	// 参数校验
	if req.PerformanceID <= 0 || req.MediaAdvId == "" || req.MediaAdvName == "" {
		return &types.ElmHcMediaCommonResp{
			Code:    400,
			Message: "客户ID、媒体账户ID和媒体账户名称不能为空",
		}, nil
	}

	media := &model.ElmHcMediaReport{
		PerformanceID: req.PerformanceID,
		MediaAdvId:    req.MediaAdvId,
		MediaAdvName:  req.MediaAdvName,
		HuichuanAdvId: req.HuichuanAdvId,
		RedirectNum:   req.RedirectNum,
		PayNum:        req.PayNum,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := model.CreateElmHcMediaReport(l.svcCtx.DB, media); err != nil {
		l.Logger.Errorf("创建汇川饿了么媒体账户失败: %v", err)
		return &types.ElmHcMediaCommonResp{
			Code:    500,
			Message: fmt.Sprintf("创建失败: %v", err),
		}, nil
	}

	return &types.ElmHcMediaCommonResp{
		Code:    0,
		Message: "创建成功",
	}, nil
}

// DeleteElmHcMediaLogic 删除汇川饿了么媒体账户
type DeleteElmHcMediaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteElmHcMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteElmHcMediaLogic {
	return &DeleteElmHcMediaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteElmHcMediaLogic) DeleteElmHcMedia(req *types.DeleteElmHcMediaReq) (*types.ElmHcMediaCommonResp, error) {
	if req.ID <= 0 {
		return &types.ElmHcMediaCommonResp{
			Code:    400,
			Message: "ID不能为空",
		}, nil
	}

	if err := model.DeleteElmHcMediaReport(l.svcCtx.DB, req.ID); err != nil {
		l.Logger.Errorf("删除汇川饿了么媒体账户失败: %v", err)
		return &types.ElmHcMediaCommonResp{
			Code:    500,
			Message: fmt.Sprintf("删除失败: %v", err),
		}, nil
	}

	return &types.ElmHcMediaCommonResp{
		Code:    0,
		Message: "删除成功",
	}, nil
}

// UpdateElmHcMediaLogic 更新汇川饿了么媒体账户
type UpdateElmHcMediaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateElmHcMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateElmHcMediaLogic {
	return &UpdateElmHcMediaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateElmHcMediaLogic) UpdateElmHcMedia(req *types.UpdateElmHcMediaReq) (*types.ElmHcMediaCommonResp, error) {
	// 参数校验
	if req.ID <= 0 || req.PerformanceID <= 0 || req.MediaAdvId == "" || req.MediaAdvName == "" {
		return &types.ElmHcMediaCommonResp{
			Code:    400,
			Message: "ID、客户ID、媒体账户ID和媒体账户名称不能为空",
		}, nil
	}

	media := &model.ElmHcMediaReport{
		ID:            req.ID,
		PerformanceID: req.PerformanceID,
		MediaAdvId:    req.MediaAdvId,
		MediaAdvName:  req.MediaAdvName,
		HuichuanAdvId: req.HuichuanAdvId,
		RedirectNum:   req.RedirectNum,
		PayNum:        req.PayNum,
	}

	if err := model.UpdateElmHcMediaReport(l.svcCtx.DB, media); err != nil {
		l.Logger.Errorf("更新汇川饿了么媒体账户失败: %v", err)
		return &types.ElmHcMediaCommonResp{
			Code:    500,
			Message: fmt.Sprintf("更新失败: %v", err),
		}, nil
	}

	return &types.ElmHcMediaCommonResp{
		Code:    0,
		Message: "更新成功",
	}, nil
}
