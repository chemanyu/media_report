package logic

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"media_report/common/httpclient"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

type GetJDAttributionDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJDAttributionDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJDAttributionDataLogic {
	return &GetJDAttributionDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetJDAttributionData 获取京东归因数据
func (l *GetJDAttributionDataLogic) GetJDAttributionData(req *types.JDAttributionRequest) (*types.JDAttributionResponse, error) {
	// 格式化日期
	formattedDate := formatDate(req.Date)

	// 调用归因API
	apiResp, err := getAttributionDataFromAPI(formattedDate)
	if err != nil {
		logx.Errorf("获取归因数据失败: %v", err)
		return nil, err
	}

	// 解析响应数据
	resp := parseAttributionResponse(apiResp)
	return resp, nil
}

// formatDate 格式化日期字符串
func formatDate(dateStr string) string {
	// 如果已经是8位数字格式，直接返回
	if len(dateStr) == 8 {
		if _, err := strconv.Atoi(dateStr); err == nil {
			return dateStr
		}
	}

	// 尝试解析常见的日期格式
	formats := []string{"2006-01-02", "2006/01/02", "2006.01.02"}
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t.Format("20060102")
		}
	}

	// 如果都无法解析，返回原字符串
	return dateStr
}

// getAttributionDataFromAPI 从API获取归因数据
func getAttributionDataFromAPI(date string) (*types.JDAttributionAPIResponse, error) {
	baseURL := "http://ad-ocpx.zhltech.net"
	client := httpclient.NewClient(baseURL, 30)

	params := map[string]string{
		"date": date,
	}

	var resp types.JDAttributionAPIResponse
	err := client.Get(context.Background(), "/attribution/data", params, &resp)
	if err != nil {
		return nil, fmt.Errorf("调用归因API失败: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API返回错误: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp, nil
}

// parseAttributionResponse 解析归因API响应
func parseAttributionResponse(apiResp *types.JDAttributionAPIResponse) *types.JDAttributionResponse {
	return &types.JDAttributionResponse{
		Date:          apiResp.Data.Date,
		TotalRequests: apiResp.Data.TotalRequests,
		ErrorCounts:   apiResp.Data.ErrorCounts,
	}
}

type ExportJDErrorCountsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportJDErrorCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportJDErrorCountsLogic {
	return &ExportJDErrorCountsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ExportJDErrorCounts 导出京东错误统计数据为Excel
func (l *ExportJDErrorCountsLogic) ExportJDErrorCounts(req *types.JDErrorExportRequest) ([]byte, string, error) {
	numDays := req.NumDays
	if numDays <= 0 {
		numDays = 10
	}

	// 获取过去N天的数据
	multiDaysData, err := l.getMultiDaysErrorData(numDays)
	if err != nil {
		return nil, "", err
	}

	// 创建Excel文件
	excelData, err := l.createExcelFile(multiDaysData)
	if err != nil {
		return nil, "", err
	}

	// 生成文件名
	filename := fmt.Sprintf("jd_error_counts_%s.xlsx", time.Now().Format("20060102_150405"))

	return excelData, filename, nil
}

// getMultiDaysErrorData 获取过去N天的错误数据
func (l *ExportJDErrorCountsLogic) getMultiDaysErrorData(numDays int) (map[string]map[string]map[string]int64, error) {
	allData := make(map[string]map[string]map[string]int64)

	for i := 0; i < numDays; i++ {
		// 计算日期
		targetDate := time.Now().AddDate(0, 0, -i)
		dateStr := targetDate.Format("20060102")

		// 获取该天数据
		apiResp, err := getAttributionDataFromAPI(dateStr)
		if err != nil {
			logx.Errorf("获取 %s 的数据失败: %v", dateStr, err)
			continue
		}

		if apiResp.Code == 0 {
			allData[dateStr] = apiResp.Data.ErrorCounts
			logx.Infof("成功获取 %s 的数据", dateStr)
		} else {
			logx.Errorf("获取 %s 的数据失败: %s", dateStr, apiResp.Message)
		}
	}

	return allData, nil
}

// createExcelFile 创建Excel文件
func (l *ExportJDErrorCountsLogic) createExcelFile(multiDaysData map[string]map[string]map[string]int64) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "错误统计"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	// 设置数据单元格样式
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 12)
	f.SetColWidth(sheetName, "B", "B", 18)
	f.SetColWidth(sheetName, "C", "C", 25)
	f.SetColWidth(sheetName, "D", "D", 12)
	f.SetColWidth(sheetName, "E", "E", 12)

	// 写入表头
	headers := []string{"日期", "账户ID", "错误类型", "事件", "错误数量"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// 按日期降序排列
	var dates []string
	for date := range multiDaysData {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	row := 2
	for _, dateStr := range dates {
		errorCounts := multiDaysData[dateStr]

		// 按账户ID排序
		var advertiserIds []string
		for advertiserId := range errorCounts {
			advertiserIds = append(advertiserIds, advertiserId)
		}
		sort.Strings(advertiserIds)

		for _, advertiserId := range advertiserIds {
			errors := errorCounts[advertiserId]

			// 按错误类型排序
			var errorTypes []string
			for errorType := range errors {
				errorTypes = append(errorTypes, errorType)
			}
			sort.Strings(errorTypes)

			for _, errorType := range errorTypes {
				errorCount := errors[errorType]

				// 拆分错误类型和事件
				var errorCategory, eventNumber string
				parts := strings.Split(errorType, "_")
				if len(parts) >= 2 {
					lastPart := parts[len(parts)-1]
					if _, err := strconv.Atoi(lastPart); err == nil {
						// 最后一部分是数字
						errorCategory = strings.Join(parts[:len(parts)-1], "_")
						eventNumber = lastPart
					} else {
						errorCategory = errorType
						eventNumber = ""
					}
				} else {
					errorCategory = errorType
					eventNumber = ""
				}

				// 写入数据
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), dateStr)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), advertiserId)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), errorCategory)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), eventNumber)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), errorCount)

				// 应用样式
				for col := 1; col <= 5; col++ {
					cell, _ := excelize.CoordinatesToCellName(col, row)
					f.SetCellStyle(sheetName, cell, cell, dataStyle)
				}

				row++
			}
		}
	}

	// 保存到内存
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
