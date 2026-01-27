package logic

import (
	"context"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateJuliangCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateJuliangCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateJuliangCookieLogic {
	return &UpdateJuliangCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateJuliangCookieLogic) UpdateJuliangCookie(req *types.UpdateJuliangCookieReq) (resp *types.UpdateJuliangCookieResp, err error) {
	// 更新 media_token 表中 media = 'juliang_pachong' 的记录
	result := l.svcCtx.DB.Model(&model.MediaToken{}).
		Where("media = ? AND del_flag = ?", "juliang_pachong", 0).
		Updates(map[string]interface{}{
			"token":         req.Cookie,
			"refresh_token": req.CsrfToken,
		})

	if result.Error != nil {
		l.Logger.Errorf("更新聚量Cookie失败: %v", result.Error)
		return &types.UpdateJuliangCookieResp{
			Code:    500,
			Message: "更新失败: " + result.Error.Error(),
		}, nil
	}

	if result.RowsAffected == 0 {
		l.Logger.Error("未找到media='juliang_pachong'的记录")
		return &types.UpdateJuliangCookieResp{
			Code:    404,
			Message: "未找到对应的媒体记录",
		}, nil
	}

	l.Logger.Infof("成功更新聚量Cookie，影响行数: %d", result.RowsAffected)

	return &types.UpdateJuliangCookieResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}
