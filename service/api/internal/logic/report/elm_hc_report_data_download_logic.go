package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

type DownloadElmHcReportDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadElmHcReportDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadElmHcReportDataLogic {
	return &DownloadElmHcReportDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DownloadElmHcReportData 查询并导出饿了么汇川报表数据为 Excel
func (l *DownloadElmHcReportDataLogic) DownloadElmHcReportData(req *types.ElmHcReportDownloadReq) ([]byte, string, error) {
	params := model.QueryElmHcReportDataParams{
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		CustomerShort: req.CustomerShort,
	}

	records, err := model.QueryElmHcReportData(l.svcCtx.DB, params)
	if err != nil {
		l.Logger.Errorf("查询饿了么汇川报表数据失败: %v", err)
		return nil, "", fmt.Errorf("查询数据失败: %w", err)
	}

	l.Logger.Infof("查询到 %d 条饿了么汇川报表数据", len(records))

	excelData, err := buildElmHcReportExcel(records)
	if err != nil {
		return nil, "", fmt.Errorf("生成Excel失败: %w", err)
	}

	filename := fmt.Sprintf("elm_hc_report_%s.xlsx", time.Now().Format("20060102150405"))
	return excelData, filename, nil
}

func buildElmHcReportExcel(records []model.ElmHcReportData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "饿了么汇川报表"
	f.SetSheetName("Sheet1", sheetName)

	// 表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	// 数据行样式
	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "D9D9D9", Style: 1},
			{Type: "right", Color: "D9D9D9", Style: 1},
			{Type: "top", Color: "D9D9D9", Style: 1},
			{Type: "bottom", Color: "D9D9D9", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	headers := []string{
		"日期", "小时", "客户名称", "客户简称", "代理名称", "代理简称",
		"媒体平台", "媒体账户ID", "媒体账户名称", "汇川账户ID",
		"消费", "展示数", "点击数", "转化数", "深度转化数",
		"转化类型", "深度转化类型", "调起数", "付费数",
	}

	// 写入表头
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheetName, 1, 20)

	// 写入数据行
	for rowIdx, r := range records {
		row := rowIdx + 2
		values := []interface{}{
			r.Dt, r.Hh, r.CustomerName, r.CustomerShort, r.AgentName, r.AgentShort,
			r.MediaPlatformName, r.MediaAdvId, r.MediaAdvName, r.HuichuanAdvId,
			r.Cost, r.ShowNum, r.ClickNum, r.ConvertNum, r.DeepConvertNum,
			r.ConvertType, r.DeepConvertType, r.RedirectNum, r.PayNum,
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheetName, cell, val)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// 设置列宽
	colWidths := []float64{10, 6, 30, 15, 30, 10, 15, 15, 25, 15, 12, 12, 12, 12, 12, 10, 12, 10, 10}
	for colIdx, width := range colWidths {
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		f.SetColWidth(sheetName, colName, colName, width)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
