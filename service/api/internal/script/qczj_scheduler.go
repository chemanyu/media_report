package script

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
)

const (
	qczjSupplier = "美数"
	qczjResource = "dsp100%"
)

var qczjHours = []string{"09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}

// ExecuteQczjReportJob 执行 QCZJ 分时监控报表任务
func ExecuteQczjReportJob(db *gorm.DB, dingTalk config.DingTalkConfig, fileServer config.FileServerConfig) {
	ctx := context.Background()
	now := time.Now()
	date := now.Format("20060102")
	logx.Infof("开始执行 QCZJ 分时监控任务 - %s", now.Format("2006-01-02 15:04:05"))

	// 读取配置
	cfg, err := model.GetQczjConfig(db)
	if err != nil {
		logx.Errorf("获取 QCZJ 配置失败: %v", err)
		return
	}

	// 查询今日各小时数据
	todayData, err := model.ListQczjReportDataByDate(db, date)
	if err != nil {
		logx.Errorf("查询 QCZJ 今日数据失败: %v", err)
		return
	}

	// 生成 Excel 文件
	filename, err := generateQczjExcel(date, todayData, cfg, fileServer.Path)
	if err != nil {
		logx.Errorf("生成 QCZJ Excel 失败: %v", err)
		return
	}

	downloadURL := fmt.Sprintf("%s/download/%s", fileServer.BaseURL, filename)
	logx.Infof("QCZJ Excel 已生成: %s，下载地址: %s", filename, downloadURL)

	// 获取上一小时数据用于钉钉推送
	lastHour := fmt.Sprintf("%02d", now.Hour()-1)
	var lastData *model.QczjReportData
	if d, ok := todayData[lastHour]; ok {
		lastData = d
	}

	sendQczjDingTalkNotification(ctx, dingTalk, now, lastHour, lastData, downloadURL, cfg)
	logx.Infof("QCZJ 分时监控任务执行完成 - %s", time.Now().Format("2006-01-02 15:04:05"))
}

// generateQczjExcel 生成分时监控 Excel 文件，返回文件名
func generateQczjExcel(date string, todayData map[string]*model.QczjReportData, cfg *model.QczjConfig, savePath string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logx.Errorf("关闭 Excel 文件失败: %v", err)
		}
	}()

	sheetName := "分时监控"
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("创建工作表失败: %w", err)
	}
	f.SetActiveSheet(idx)
	f.DeleteSheet("Sheet1")

	// ---------- 列宽 ----------
	colWidths := map[string]float64{
		"A": 14, "B": 8, "C": 8, "D": 12, "E": 12,
		"F": 12, "G": 12, "H": 8, "I": 14, "J": 16,
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// ---------- 表头样式（行3）----------
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1F3864"}},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建表头样式失败: %w", err)
	}

	headers := map[string]string{
		"A3": "供应商名称", "B3": "日期", "C3": "时段",
		"D3": "曝光", "E3": "点击", "F3": "花费",
		"G3": "唤起uv", "H3": "今日预估总uv", "J3": "投放资源及占比",
	}
	for cell, val := range headers {
		f.SetCellValue(sheetName, cell, val)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	f.MergeCell(sheetName, "H3", "I3")
	f.MergeCell(sheetName, "J3", "K3")
	f.SetRowHeight(sheetName, 3, 20)

	// ---------- 数据行（行4起）----------
	startRow := 4
	hourCount := len(qczjHours)
	endRow := startRow + hourCount - 1

	// 解析日期中的天数
	t, _ := time.Parse("20060102", date)
	dayNum := t.Day()

	// 跨行合并
	f.MergeCell(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", endRow))
	f.MergeCell(sheetName, fmt.Sprintf("B%d", startRow), fmt.Sprintf("B%d", endRow))
	f.MergeCell(sheetName, fmt.Sprintf("H%d", startRow), fmt.Sprintf("I%d", endRow))
	f.MergeCell(sheetName, fmt.Sprintf("J%d", startRow), fmt.Sprintf("K%d", endRow))

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), qczjSupplier)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow), dayNum)
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", startRow), cfg.TotalUv)
	f.SetCellValue(sheetName, fmt.Sprintf("J%d", startRow), qczjResource)

	centerMiddleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return "", fmt.Errorf("创建居中样式失败: %w", err)
	}

	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", endRow), centerMiddleStyle)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", startRow), fmt.Sprintf("B%d", endRow), centerMiddleStyle)

	// 今日预估总UV - 红色加粗居中（先设置，整体边框后会被覆盖，最终在边框区域下方重新应用）
	uvStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FF0000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return "", fmt.Errorf("创建UV样式失败: %w", err)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("H%d", startRow), fmt.Sprintf("H%d", startRow), uvStyle)

	rightStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})
	if err != nil {
		return "", fmt.Errorf("创建右对齐样式失败: %w", err)
	}

	dataCenterStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return "", fmt.Errorf("创建数据居中样式失败: %w", err)
	}

	// 逐行填写时段数据
	for i, hour := range qczjHours {
		row := startRow + i
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), hour+":00")
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), rightStyle)

		if data, ok := todayData[hour]; ok {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.View)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.Click)
			// F列花费留空
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), int64(math.Round(float64(data.Action)*cfg.Ratio)))
			f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("G%d", row), dataCenterStyle)
		}
	}

	// 整体边框
	borderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建边框样式失败: %w", err)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("K%d", endRow), borderStyle)

	// 今日预估总UV - 红色加粗居中（含边框，覆盖整体边框后重新应用）
	uvStyleWithBorder, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FF0000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建UV样式失败: %w", err)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("H%d", startRow), fmt.Sprintf("H%d", startRow), uvStyleWithBorder)

	// 投放资源及占比 - 居中（含边框）
	centerBorderStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建居中边框样式失败: %w", err)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("J%d", startRow), fmt.Sprintf("J%d", startRow), centerBorderStyle)

	// ---------- 环比行 ----------
	compRow := endRow + 1
	f.MergeCell(sheetName, fmt.Sprintf("A%d", compRow), fmt.Sprintf("B%d", compRow))
	f.MergeCell(sheetName, fmt.Sprintf("D%d", compRow), fmt.Sprintf("K%d", compRow))
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", compRow), "环比")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", compRow), "环比昨天同时段数据的变化")

	compTextStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "FF0000"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return "", fmt.Errorf("创建环比文字样式失败: %w", err)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("D%d", compRow), fmt.Sprintf("D%d", compRow), compTextStyle)
	f.SetCellStyle(sheetName, fmt.Sprintf("C%d", compRow), fmt.Sprintf("C%d", compRow), rightStyle)
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", compRow), fmt.Sprintf("K%d", compRow), borderStyle)

	// ---------- 保存文件 ----------
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	filename := fmt.Sprintf("拉活分时监控-美数-%s.xlsx", date)
	fullPath := filepath.Join(savePath, filename)
	if err := f.SaveAs(fullPath); err != nil {
		return "", fmt.Errorf("保存 Excel 文件失败: %w", err)
	}

	return filename, nil
}

// sendQczjDingTalkNotification 发送 QCZJ 分时监控钉钉通知
func sendQczjDingTalkNotification(ctx context.Context, dingTalk config.DingTalkConfig,
	now time.Time, lastHour string, data *model.QczjReportData, downloadURL string, cfg *model.QczjConfig) {

	if !dingTalk.Enabled || dingTalk.QczjWebhookURL == "" {
		logx.Info("QCZJ 钉钉通知未启用，跳过发送")
		return
	}

	timeStr := fmt.Sprintf("%s %s时", now.Format("2006-01-02"), lastHour)

	var view, click, action int64
	if data != nil {
		view = data.View
		click = data.Click
		action = int64(math.Round(float64(data.Action) * cfg.Ratio))
	}

	markdownText := fmt.Sprintf(
		"#### 美数-分时监控时报  \n---\n"+
			"**时间**：%s  \n"+
			"**曝光(view)**：%d  \n"+
			"**点击(click)**：%d  \n"+
			"**唤起UV(action)**：%d  \n\n"+
			"详细分时数据请下载文件：[下载](%s)",
		timeStr,
		view,
		click,
		action,
		downloadURL,
	)

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "QCZJ 分时监控时报",
			"text":  markdownText,
		},
	}

	client := httpclient.NewClient("", 30)
	var result map[string]interface{}
	if err := client.Post(ctx, dingTalk.QczjWebhookURL, msg, &result); err != nil {
		logx.Errorf("发送 QCZJ 钉钉消息失败: %v", err)
		return
	}

	logx.Infof("QCZJ 钉钉消息发送成功: %v", result)
}
