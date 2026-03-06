package qczj

import (
	"context"
	"fmt"
	"strconv"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncAdnDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncAdnDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncAdnDataLogic {
	return &SyncAdnDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncAdnDataLogic) SaveData(req *types.QczjSyncDataReq) error {
	reportDate, err := strconv.Atoi(req.ReportDate)
	if err != nil {
		return fmt.Errorf("report_date 格式错误: %w", err)
	}

	data := &model.QczjReportData{
		ReportDate: reportDate,
		View:       req.View,
		Click:      req.Click,
		Action:     req.Action,
	}

	if err := model.InsertOrUpdateQczjReportData(l.svcCtx.DB, data); err != nil {
		return err
	}

	l.Infof("qczj: 保存成功 date=%s view=%d click=%d action=%d", req.ReportDate, req.View, req.Click, req.Action)
	return nil
}
