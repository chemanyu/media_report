package logic

import (
	"context"
	"fmt"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

// FzAdvertiserLogic 飞猪媒体账户业务逻辑
type FzAdvertiserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFzAdvertiserLogic 创建飞猪媒体账户逻辑实例
func NewFzAdvertiserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FzAdvertiserLogic {
	return &FzAdvertiserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetList 获取账户列表
func (l *FzAdvertiserLogic) GetList() ([]*model.FzMediaAdvertiser, error) {
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)
	return advertiserModel.FindAll()
}

// GetAdnList 获取ADN账户列表
func (l *FzAdvertiserLogic) GetAdnList() (*types.FzAdnAdvertiserListResp, error) {
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)
	list, err := advertiserModel.FindByMedia("adn")
	if err != nil {
		return nil, err
	}
	items := make([]*types.FzAdnAdvertiserItem, 0, len(list))
	for _, a := range list {
		items = append(items, &types.FzAdnAdvertiserItem{
			Id:           a.Id,
			MediaAdvId:   a.MediaAdvId,
			MediaAdvName: a.MediaAdvName,
		})
	}
	return &types.FzAdnAdvertiserListResp{List: items}, nil
}

// Add 添加账户
func (l *FzAdvertiserLogic) Add(req *types.FzAdvertiserAddReq) (*types.CommonResp, error) {
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)

	// 检查账户ID是否已存在
	existing, _ := advertiserModel.FindByMediaAdvId(req.MediaAdvId)
	if existing != nil {
		return nil, fmt.Errorf("账户ID已存在: %s", req.MediaAdvId)
	}

	// 创建新账户
	advertiser := &model.FzMediaAdvertiser{
		Media:        req.Media,
		MediaAdvId:   req.MediaAdvId,
		MediaAdvName: req.MediaAdvName,
	}

	if err := advertiserModel.Insert(advertiser); err != nil {
		return nil, fmt.Errorf("添加账户失败: %w", err)
	}

	return &types.CommonResp{
		Code:    0,
		Message: "添加成功",
	}, nil
}

// Update 更新账户
func (l *FzAdvertiserLogic) Update(req *types.FzAdvertiserUpdateReq) (*types.CommonResp, error) {
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)

	// 检查账户是否存在
	existing, err := advertiserModel.FindOne(req.Id)
	if err != nil {
		return nil, fmt.Errorf("账户不存在: id=%d", req.Id)
	}

	// 更新账户名称
	existing.MediaAdvName = req.MediaAdvName

	if err := advertiserModel.Update(existing); err != nil {
		return nil, fmt.Errorf("更新账户失败: %w", err)
	}

	return &types.CommonResp{
		Code:    0,
		Message: "更新成功",
	}, nil
}

// Delete 删除账户
func (l *FzAdvertiserLogic) Delete(req *types.FzAdvertiserDeleteReq) (*types.CommonResp, error) {
	advertiserModel := model.NewFzMediaAdvertiserModel(l.svcCtx.DB)

	if err := advertiserModel.Delete(req.Id); err != nil {
		return nil, fmt.Errorf("删除账户失败: %w", err)
	}

	return &types.CommonResp{
		Code:    0,
		Message: "删除成功",
	}, nil
}
