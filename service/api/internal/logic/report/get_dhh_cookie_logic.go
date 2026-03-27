package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDhhCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDhhCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDhhCookieLogic {
	return &GetDhhCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDhhCookieLogic) GetDhhCookie() (resp *types.GetDhhCookieResp, err error) {
	// 从 media_token 表中获取 media = 'Dhh_pachong' 的 token
	var token model.MediaToken
	result := l.svcCtx.DB.Where("media = ? AND del_flag = ?", "dhh_pachong", 0).First(&token)

	if result.Error != nil {
		l.Logger.Errorf("获取京橙Cookie失败: %v", result.Error)
		return &types.GetDhhCookieResp{
			Code:    500,
			Message: "查询失败: " + result.Error.Error(),
			Data:    "",
		}, nil
	}

	l.Logger.Infof("成功获取京橙Cookie")

	return &types.GetDhhCookieResp{
		Code:    200,
		Message: "查询成功",
		Data:    token.Token,
		Csrf:    token.RefreshToken,
	}, nil
}
