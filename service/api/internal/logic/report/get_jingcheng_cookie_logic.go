package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJingchengCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJingchengCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJingchengCookieLogic {
	return &GetJingchengCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJingchengCookieLogic) GetJingchengCookie() (resp *types.GetJingchengCookieResp, err error) {
	// 从 media_token 表中获取 media = 'jingcheng_pachong' 的 token
	var token model.MediaToken
	result := l.svcCtx.DB.Where("media = ? AND del_flag = ?", "jingcheng_pachong", 0).First(&token)

	if result.Error != nil {
		l.Logger.Errorf("获取京橙Cookie失败: %v", result.Error)
		return &types.GetJingchengCookieResp{
			Code:    500,
			Message: "查询失败: " + result.Error.Error(),
			Data:    "",
		}, nil
	}

	l.Logger.Infof("成功获取京橙Cookie")

	return &types.GetJingchengCookieResp{
		Code:    200,
		Message: "查询成功",
		Data:    token.Token,
	}, nil
}
