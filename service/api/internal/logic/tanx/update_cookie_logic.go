package tanx

import (
	"context"
	"fmt"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCookieLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCookieLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCookieLogic {
	return &UpdateCookieLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateCookie 更新Cookie到数据库
func (l *UpdateCookieLogic) UpdateCookie(req *types.TanxUpdateCookieReq) (resp *types.TanxUpdateCookieResp, err error) {
	if req.Cookie == "" {
		return nil, fmt.Errorf("Cookie不能为空")
	}

	// 更新数据库中的 token
	var mediaToken model.MediaToken
	err = l.svcCtx.DB.Where("media = ? AND del_flag = ?", "tanx_pachong", 0).First(&mediaToken).Error
	if err != nil {
		// 如果不存在，则创建新记录
		mediaToken = model.MediaToken{
			Media:   "tanx_pachong",
			Token:   req.Cookie,
			DelFlag: 0,
		}
		if err = l.svcCtx.DB.Create(&mediaToken).Error; err != nil {
			return nil, fmt.Errorf("创建 token 记录失败: %w", err)
		}
	} else {
		// 如果存在，则更新
		if err = l.svcCtx.DB.Model(&mediaToken).Update("token", req.Cookie).Error; err != nil {
			return nil, fmt.Errorf("更新 token 失败: %w", err)
		}
	}

	updateTime := time.Now().Format("2006-01-02 15:04:05")
	l.Logger.Infof("Cookie已更新到数据库 - 更新时间: %s, Cookie长度: %d", updateTime, len(req.Cookie))

	return &types.TanxUpdateCookieResp{
		Message:    "Cookie 已更新成功！",
		UpdateTime: updateTime,
	}, nil
}
