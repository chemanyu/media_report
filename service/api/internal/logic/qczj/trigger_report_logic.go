package qczj

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/script"
	"media_report/service/api/internal/svc"
)

type TriggerQczjReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerQczjReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerQczjReportLogic {
	return &TriggerQczjReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerQczjReportLogic) TriggerReport() (map[string]interface{}, error) {
	l.Logger.Info("手动触发 QCZJ 分时监控报表任务")

	go func() {
		script.ExecuteQczjReportJob(l.svcCtx.DB, l.svcCtx.Config.DingTalk, l.svcCtx.Config.FileServer)
	}()

	return map[string]interface{}{
		"code":    0,
		"message": "QCZJ 报表任务已触发，正在后台执行",
	}, nil
}
