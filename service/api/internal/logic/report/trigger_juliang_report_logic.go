package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/script"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

type TriggerJuliangReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerJuliangReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerJuliangReportLogic {
	return &TriggerJuliangReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerJuliangReportLogic) TriggerJuliangReport() (resp *types.TriggerJuliangReportResp, err error) {
	// 调用巨量报表任务
	script.ExecuteJuliangReportJob(l.svcCtx.DB, l.svcCtx.Config.DingTalk, l.svcCtx.Config.FileServer)

	return &types.TriggerJuliangReportResp{
		Code:    0,
		Message: "巨量报表任务已触发",
	}, nil
}
