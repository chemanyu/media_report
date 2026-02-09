package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateZfbCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateZfbCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateZfbCookieLogic {
	return &UpdateZfbCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateZfbCookieLogic) UpdateZfbCookie(req *types.UpdateZfbCookieReq) (resp *types.UpdateZfbCookieResp, err error) {
	// 更新 media_token 表中 media = 'zfb_pachong' 的记录
	//l.Logger.Infof("zfb: ", req.Token)
	result := l.svcCtx.DB.Model(&model.MediaToken{}).
		Where("media = ? AND del_flag = ?", "zfb_pachong", 0).
		Updates(map[string]interface{}{
			"token":         req.Token,
			"refresh_token": req.RefreshToken,
		})

	if result.Error != nil {
		l.Logger.Errorf("更新支付宝Cookie失败: %v", result.Error)
		return &types.UpdateZfbCookieResp{
			Code:    500,
			Message: "更新失败: " + result.Error.Error(),
		}, nil
	}

	if result.RowsAffected == 0 {
		l.Logger.Error("未找到media='zfb_pachong'的记录")
		return &types.UpdateZfbCookieResp{
			Code:    404,
			Message: "未找到对应的媒体记录",
		}, nil
	}

	l.Logger.Infof("成功更新支付宝Cookie，影响行数: %d", result.RowsAffected)

	return &types.UpdateZfbCookieResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}
