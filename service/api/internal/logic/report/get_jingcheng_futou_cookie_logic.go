package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJingchengFutouCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJingchengFutouCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJingchengFutouCookieLogic {
	return &GetJingchengFutouCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJingchengFutouCookieLogic) GetJingchengFutouCookie() (resp *types.GetJingchengFutouCookieResp, err error) {
	// 从 media_token 表中获取 media = 'jingcheng_futou' 的 token
	var token model.MediaToken
	result := l.svcCtx.DB.Where("media = ? AND del_flag = ?", "jingcheng_futou", 0).First(&token)

	if result.Error != nil {
		l.Logger.Errorf("获取京橙复投Cookie失败: %v", result.Error)
		return &types.GetJingchengFutouCookieResp{
			Code:    500,
			Message: "查询失败: " + result.Error.Error(),
			Data:    "",
		}, nil
	}

	l.Logger.Infof("成功获取京橙复投Cookie")

	return &types.GetJingchengFutouCookieResp{
		Code:    200,
		Message: "查询成功",
		Data:    token.Token,
	}, nil
}
