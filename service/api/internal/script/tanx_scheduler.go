package script

import (
	"context"
	"time"

	"media_report/service/api/internal/logic/tanx"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// ExecuteTanxFetchDataJob Tanx 数据抓取任务
func ExecuteTanxFetchDataJob(db *gorm.DB) {
	logx.Info("开始执行 Tanx 数据抓取任务")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
	}

	// 调用 FetchData 逻辑
	logic := tanx.NewFetchDataLogic(ctx, svcCtx)
	resp, err := logic.FetchData(&types.TanxFetchDataReq{})
	if err != nil {
		logx.Errorf("Tanx 数据抓取任务失败: %v", err)
		return
	}

	logx.Infof("Tanx 数据抓取任务执行完成: %s", resp.Message)
}

// ExecuteTanxExportDataJob Tanx 数据导出任务
func ExecuteTanxExportDataJob(db *gorm.DB) {
	logx.Info("开始执行 Tanx 数据导出任务")

	// 等待5秒，确保数据已写入数据库
	time.Sleep(5 * time.Second)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
	}

	// 调用 ExportData 逻辑
	logic := tanx.NewExportDataLogic(ctx, svcCtx)
	resp, err := logic.ExportData(&types.TanxExportDataReq{})
	if err != nil {
		logx.Errorf("Tanx 数据导出任务失败: %v", err)
		return
	}

	logx.Infof("Tanx 数据导出任务执行完成: %s", resp.Message)
}
