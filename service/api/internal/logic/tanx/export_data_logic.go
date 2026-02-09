package tanx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/gomail.v2"
)

type DataRow struct {
	Ds            string
	Pid           string
	AdzoneName    string
	Qingqiupv     int64
	ActiveRatioDf string
	TanxEffectPv  int64
	TanxClk       int64
	DongfengEf    float64
}

type ExportDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportDataLogic {
	return &ExportDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ExportData 查询数据库并导出Excel，发送邮件
func (l *ExportDataLogic) ExportData(req *types.TanxExportDataReq) (resp *types.TanxExportDataResp, err error) {
	// 计算30天前的日期
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	// 使用模型方法查询数据
	monitors, err := model.GetTanxMonitorsByDateRange(l.svcCtx.DB, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %w", err)
	}

	// 转换为 DataRow 格式
	dataRows := make([]DataRow, len(monitors))
	for i, m := range monitors {
		dataRows[i] = DataRow{
			Ds:            m.Ds,
			Pid:           m.Pid,
			AdzoneName:    m.AdzoneName,
			Qingqiupv:     m.Qingqiupv,
			ActiveRatioDf: m.ActiveRatioDf,
			TanxEffectPv:  m.TanxEffectPv,
			TanxClk:       m.TanxClk,
			DongfengEf:    m.DongfengEf,
		}
	}

	// 创建Excel文件
	filePath, err := l.createExcelFile(dataRows)
	if err != nil {
		return nil, fmt.Errorf("创建Excel文件失败: %w", err)
	}

	// 发送邮件
	if err := l.sendEmail(filePath); err != nil {
		return nil, fmt.Errorf("发送邮件失败: %w", err)
	}

	l.Logger.Infof("数据导出并发送邮件成功: %s", filePath)

	return &types.TanxExportDataResp{
		Message: fmt.Sprintf("数据导出成功，已发送邮件，共 %d 条数据", len(dataRows)),
	}, nil
}

// createExcelFile 创建Excel文件
func (l *ExportDataLogic) createExcelFile(dataRows []DataRow) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "详细数据"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}

	// 设置表头
	headers := []string{"日期", "广告位", "广告位名称", "tanx有效请求",
		"东风手淘换端率-同步点击", "TANX曝光数", "TANX点击数", "TANX预估收益"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for rowIdx, row := range dataRows {
		rowNum := rowIdx + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row.Ds)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row.Pid)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), row.AdzoneName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), row.Qingqiupv)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), row.ActiveRatioDf)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), row.TanxEffectPv)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), row.TanxClk)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), row.DongfengEf)
	}

	// 设置默认sheet
	f.SetActiveSheet(index)

	// 调整列宽
	l.adjustColumnWidth(f, sheetName)

	// 保存文件
	tmpDir := "/tmp"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		tmpDir = "."
	}

	filePath := filepath.Join(tmpDir, fmt.Sprintf("tanx_data_%s.xlsx",
		time.Now().Format("20060102_150405")))

	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

// adjustColumnWidth 调整列宽
func (l *ExportDataLogic) adjustColumnWidth(f *excelize.File, sheetName string) {
	columnWidths := map[string]float64{
		"A": 12, // 日期
		"B": 25, // 广告位
		"C": 30, // 广告位名称
		"D": 15, // tanx有效请求
		"E": 25, // 东风手淘换端率-同步点击
		"F": 15, // TANX曝光数
		"G": 15, // TANX点击数
		"H": 15, // TANX预估收益
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}
}

// sendEmail 发送邮件
func (l *ExportDataLogic) sendEmail(filePath string) error {
	config := GetTanxConfig()
	smtpConfig := GetSMTPConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.User)
	m.SetHeader("To", config.Recipients...)
	m.SetHeader("Subject", "Tanx Data Export")
	m.SetBody("text/plain", "Please find the attached Excel file containing the Tanx data for the last 30 days.")
	m.Attach(filePath)

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	l.Logger.Infof("邮件已发送到: %v", config.Recipients)
	return nil
}
