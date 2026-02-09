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
	startDate := time.Now().AddDate(0, 0, -30).Format("20060102")
	endDate := time.Now().Format("20060102")

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

// createExcelFile 创建Excel文件（完全复刻Python逻辑）
func (l *ExportDataLogic) createExcelFile(dataRows []DataRow) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 1. 处理数据，添加媒体统计（复刻 add_media_statistics）
	processedRows := l.addMediaStatistics(dataRows)

	// 2. 创建详细数据 sheet
	detailSheetName := "详细数据"
	detailIndex, err := f.NewSheet(detailSheetName)
	if err != nil {
		return "", err
	}

	// 设置详细数据表头
	headers := []string{"日期", "广告位", "广告位名称", "tanx有效请求",
		"东风手淘换端率-同步点击", "TANX曝光数", "TANX点击数", "TANX预估收益"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(detailSheetName, cell, header)
	}

	// 创建样式：左对齐（普通数据）和右对齐（统计总数）
	leftAlignStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})

	rightAlignStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})

	// 写入处理后的数据（包含统计行和空行）
	rowNum := 2
	for _, row := range processedRows {
		// 判断是否为统计总数行
		isSummaryRow := row.Pid == "统计总数"

		// 选择对齐样式
		alignStyle := leftAlignStyle
		if isSummaryRow {
			alignStyle = rightAlignStyle
		}

		// 写入数据并设置样式
		f.SetCellValue(detailSheetName, fmt.Sprintf("A%d", rowNum), row.Ds)
		f.SetCellStyle(detailSheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("A%d", rowNum), alignStyle)

		f.SetCellValue(detailSheetName, fmt.Sprintf("B%d", rowNum), row.Pid)
		f.SetCellStyle(detailSheetName, fmt.Sprintf("B%d", rowNum), fmt.Sprintf("B%d", rowNum), alignStyle)

		f.SetCellValue(detailSheetName, fmt.Sprintf("C%d", rowNum), row.AdzoneName)
		f.SetCellStyle(detailSheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), alignStyle)

		if row.Qingqiupv > 0 {
			f.SetCellValue(detailSheetName, fmt.Sprintf("D%d", rowNum), row.Qingqiupv)
			f.SetCellStyle(detailSheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("D%d", rowNum), alignStyle)
		}

		f.SetCellValue(detailSheetName, fmt.Sprintf("E%d", rowNum), row.ActiveRatioDf)
		f.SetCellStyle(detailSheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("E%d", rowNum), alignStyle)

		if row.TanxEffectPv > 0 {
			f.SetCellValue(detailSheetName, fmt.Sprintf("F%d", rowNum), row.TanxEffectPv)
			f.SetCellStyle(detailSheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("F%d", rowNum), alignStyle)
		}

		if row.TanxClk > 0 {
			f.SetCellValue(detailSheetName, fmt.Sprintf("G%d", rowNum), row.TanxClk)
			f.SetCellStyle(detailSheetName, fmt.Sprintf("G%d", rowNum), fmt.Sprintf("G%d", rowNum), alignStyle)
		}

		if row.DongfengEf > 0 {
			f.SetCellValue(detailSheetName, fmt.Sprintf("H%d", rowNum), row.DongfengEf)
			f.SetCellStyle(detailSheetName, fmt.Sprintf("H%d", rowNum), fmt.Sprintf("H%d", rowNum), alignStyle)
		}

		rowNum++
	}

	// 3. 创建汇总统计 sheet（从详细数据中提取"统计总数"行）
	summarySheetName := "汇总统计"
	_, err = f.NewSheet(summarySheetName)
	if err != nil {
		l.Logger.Errorf("创建汇总统计sheet失败: %v", err)
	} else {
		l.createSummarySheetFromDetails(f, summarySheetName, processedRows)
	}

	// 设置默认sheet为详细数据
	f.SetActiveSheet(detailIndex)

	// 删除默认的 Sheet1
	f.DeleteSheet("Sheet1")

	// 调整所有sheet的列宽
	l.adjustAllSheetsColumnWidth(f)

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

	l.Logger.Infof("Excel文件创建成功，包含详细数据和汇总统计")
	return filePath, nil
}

// extractMedia 从广告位名称中提取媒体类型
func (l *ExportDataLogic) extractMedia(adzoneName string) string {
	if adzoneName == "" {
		return "未知"
	}

	mediaKeys := []string{"佳投", "有境", "新数", "快友", "多盟", "浩睿", "美数"}
	for _, key := range mediaKeys {
		if contains(adzoneName, key) {
			return key
		}
	}
	return "其他"
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// addMediaStatistics 添加媒体统计（完全复刻Python的add_media_statistics）
func (l *ExportDataLogic) addMediaStatistics(dataRows []DataRow) []DataRow {
	var resultRows []DataRow

	// 按日期分组
	dateGroups := make(map[string][]DataRow)
	dateOrder := []string{}

	for _, row := range dataRows {
		if _, exists := dateGroups[row.Ds]; !exists {
			dateOrder = append(dateOrder, row.Ds)
		}
		dateGroups[row.Ds] = append(dateGroups[row.Ds], row)
	}

	// 按日期处理
	for _, date := range dateOrder {
		dateRows := dateGroups[date]

		// 格式化日期 2006-01-02 -> 2006/01/02
		formattedDate := formatDateSlash(date)

		// 按媒体分组
		mediaGroups := make(map[string][]DataRow)
		mediaOrder := []string{}

		for _, row := range dateRows {
			media := l.extractMedia(row.AdzoneName)
			if _, exists := mediaGroups[media]; !exists {
				mediaOrder = append(mediaOrder, media)
			}
			mediaGroups[media] = append(mediaGroups[media], row)
		}

		// 按媒体处理
		for _, media := range mediaOrder {
			mediaRows := mediaGroups[media]

			// 1. 添加该媒体当天的所有详细数据
			for _, row := range mediaRows {
				row.Ds = formattedDate
				resultRows = append(resultRows, row)
			}

			// 2. 添加该媒体当天的统计汇总行
			statsRow := l.calculateMediaStats(mediaRows, formattedDate, media)
			resultRows = append(resultRows, statsRow)
		}

		// 3. 每天结束后添加空行分隔
		emptyRow := DataRow{Ds: "", Pid: "", AdzoneName: "", Qingqiupv: 0, ActiveRatioDf: "", TanxEffectPv: 0, TanxClk: 0, DongfengEf: 0}
		resultRows = append(resultRows, emptyRow)
	}

	l.Logger.Infof("处理后的数据行数: %d", len(resultRows))
	return resultRows
}

// formatDateSlash 格式化日期：2006-01-02 -> 2006/01/02
func formatDateSlash(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("2006/01/02")
}

// calculateMediaStats 计算媒体统计数据
func (l *ExportDataLogic) calculateMediaStats(mediaRows []DataRow, formattedDate, media string) DataRow {
	var totalQingqiupv, totalTanxEffectPv, totalTanxClk int64
	var totalDongfengEf float64
	var ratioSum float64
	var ratioCount int

	for _, row := range mediaRows {
		totalQingqiupv += row.Qingqiupv
		totalTanxEffectPv += row.TanxEffectPv
		totalTanxClk += row.TanxClk
		totalDongfengEf += row.DongfengEf

		// 解析百分比格式的 ActiveRatioDf
		ratioStr := row.ActiveRatioDf
		if ratioStr != "" && len(ratioStr) > 1 && ratioStr[len(ratioStr)-1] == '%' {
			var ratio float64
			fmt.Sscanf(ratioStr[:len(ratioStr)-1], "%f", &ratio)
			ratioSum += ratio
			ratioCount++
		}
	}

	// 计算平均换端率
	avgRatio := 0.0
	if ratioCount > 0 {
		avgRatio = ratioSum / float64(ratioCount)
	}

	return DataRow{
		Ds:            formattedDate,
		Pid:           "统计总数",
		AdzoneName:    media,
		Qingqiupv:     totalQingqiupv,
		ActiveRatioDf: fmt.Sprintf("%.2f%%", avgRatio),
		TanxEffectPv:  totalTanxEffectPv,
		TanxClk:       totalTanxClk,
		DongfengEf:    totalDongfengEf,
	}
}

// createSummarySheetFromDetails 从详细数据中提取"统计总数"行创建汇总表
func (l *ExportDataLogic) createSummarySheetFromDetails(f *excelize.File, sheetName string, processedRows []DataRow) {
	// 筛选出所有统计汇总行（Pid为"统计总数"的行）
	var summaryRows []DataRow
	for _, row := range processedRows {
		if row.Pid == "统计总数" {
			summaryRows = append(summaryRows, row)
		}
	}

	if len(summaryRows) == 0 {
		l.Logger.Info("未找到统计汇总行")
		return
	}

	// 设置汇总表头
	headers := []string{"日期", "广告位", "广告位名称", "tanx有效请求合计",
		"东风手淘换端率-同步点击", "TANX曝光数合计", "TANX点击数合计", "TANX预估收益合计"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入汇总数据
	rowNum := 2
	for _, row := range summaryRows {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row.Ds)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row.Pid)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), row.AdzoneName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), row.Qingqiupv)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), row.ActiveRatioDf)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), row.TanxEffectPv)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), row.TanxClk)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), row.DongfengEf)
		rowNum++
	}

	l.Logger.Infof("汇总统计sheet创建成功，共 %d 条汇总数据", len(summaryRows))
}

// adjustAllSheetsColumnWidth 调整所有sheet的列宽
func (l *ExportDataLogic) adjustAllSheetsColumnWidth(f *excelize.File) {
	sheetList := f.GetSheetList()

	for _, sheetName := range sheetList {
		if sheetName == "Sheet1" {
			continue
		}

		// 获取所有行
		rows, err := f.GetRows(sheetName)
		if err != nil {
			l.Logger.Errorf("获取sheet %s 的行数据失败: %v", sheetName, err)
			continue
		}

		// 计算每列的最大宽度
		colWidths := make(map[int]int)
		for _, row := range rows {
			for colIdx, cell := range row {
				cellLen := len(cell)
				if cellLen > colWidths[colIdx] {
					colWidths[colIdx] = cellLen
				}
			}
		}

		// 设置列宽（最小10，最大50）
		for colIdx, width := range colWidths {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			adjustedWidth := float64(width + 2)
			if adjustedWidth < 10 {
				adjustedWidth = 10
			}
			if adjustedWidth > 50 {
				adjustedWidth = 50
			}
			f.SetColWidth(sheetName, colName, colName, adjustedWidth)
		}
	}

	l.Logger.Info("所有sheet列宽调整完成")
}

// adjustColumnWidth 调整列宽（已废弃，保留兼容）
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
