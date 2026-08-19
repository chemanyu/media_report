package logic

import (
	"context"
	"errors"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// jingchengFutouMedia media_token 表中京橙复投的媒体标识
const jingchengFutouMedia = "jingcheng_futou"

type UpdateJingchengFutouCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateJingchengFutouCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateJingchengFutouCookieLogic {
	return &UpdateJingchengFutouCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateJingchengFutouCookieLogic) UpdateJingchengFutouCookie(req *types.UpdateJingchengFutouCookieReq) (resp *types.UpdateJingchengFutouCookieResp, err error) {
	// 先查记录是否存在：不存在则插入，存在则更新。
	// 不能只靠 Updates 的 RowsAffected 判断——Cookie 未变化时 MySQL 也返回 0，会误判成"记录不存在"。
	var token model.MediaToken
	queryResult := l.svcCtx.DB.Where("media = ? AND del_flag = ?", jingchengFutouMedia, 0).First(&token)

	if queryResult.Error != nil {
		if !errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			l.Logger.Errorf("查询京橙复投Cookie失败: %v", queryResult.Error)
			return &types.UpdateJingchengFutouCookieResp{
				Code:    500,
				Message: "更新失败: " + queryResult.Error.Error(),
			}, nil
		}

		// 首次上报，自动补一条 media_token 记录
		newToken := model.MediaToken{
			Media:        jingchengFutouMedia,
			Token:        req.Token,
			RefreshToken: req.RefreshToken,
		}
		if createResult := l.svcCtx.DB.Create(&newToken); createResult.Error != nil {
			l.Logger.Errorf("新建京橙复投Cookie记录失败: %v", createResult.Error)
			return &types.UpdateJingchengFutouCookieResp{
				Code:    500,
				Message: "更新失败: " + createResult.Error.Error(),
			}, nil
		}

		l.Logger.Infof("已新建 media='%s' 的Cookie记录，id: %d", jingchengFutouMedia, newToken.ID)
		return &types.UpdateJingchengFutouCookieResp{
			Code:    200,
			Message: "更新成功",
		}, nil
	}

	result := l.svcCtx.DB.Model(&model.MediaToken{}).
		Where("id = ?", token.ID).
		Updates(map[string]interface{}{
			"token":         req.Token,
			"refresh_token": req.RefreshToken,
			"update_time":   time.Now().Unix(),
		})

	if result.Error != nil {
		l.Logger.Errorf("更新京橙复投Cookie失败: %v", result.Error)
		return &types.UpdateJingchengFutouCookieResp{
			Code:    500,
			Message: "更新失败: " + result.Error.Error(),
		}, nil
	}

	l.Logger.Infof("成功更新京橙复投Cookie，影响行数: %d", result.RowsAffected)

	return &types.UpdateJingchengFutouCookieResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}
