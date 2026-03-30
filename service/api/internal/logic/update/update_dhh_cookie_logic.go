package logic

import (
	"context"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDhhCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDhhCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDhhCookieLogic {
	return &UpdateDhhCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDhhCookieLogic) UpdateDhhCookie(req *types.UpdateDhhCookieReq) (resp *types.UpdateDhhCookieResp, err error) {
	// 更新 media_token 表中 media = 'Dhh_pachong' 的记录
	result := l.svcCtx.DB.Model(&model.MediaToken{}).
		Where("media = ? AND del_flag = ?", "dhh_pachong", 0).
		Updates(map[string]interface{}{
			"token":         req.Token,
			"refresh_token": req.CsrfToken,
			"update_time":   time.Now().Unix(),
		})

	if result.Error != nil {
		l.Logger.Errorf("更新京橙Cookie失败: %v", result.Error)
		return &types.UpdateDhhCookieResp{
			Code:    500,
			Message: "更新失败: " + result.Error.Error(),
		}, nil
	}

	if result.RowsAffected == 0 {
		l.Logger.Error("未找到media='Dhh_pachong'的记录")
		return &types.UpdateDhhCookieResp{
			Code:    404,
			Message: "未找到对应的媒体记录",
		}, nil
	}

	l.Logger.Infof("成功更新京橙Cookie，影响行数: %d", result.RowsAffected)

	return &types.UpdateDhhCookieResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}
