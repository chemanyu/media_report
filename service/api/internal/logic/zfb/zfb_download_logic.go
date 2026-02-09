package zfb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type ZfbDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewZfbDownloadLogic 创建 ZFB 下载逻辑实例
func NewZfbDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ZfbDownloadLogic {
	return &ZfbDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ZfbDownload ZFB 独立下载接口的业务逻辑
func (l *ZfbDownloadLogic) ZfbDownload(req *types.ZfbDownloadReq) (filePath string, filename string, err error) {
	logx.Infof("[ZFB Download] 收到下载请求 - Uid: %s, StartDate: %s, EndDate: %s", req.Uid, req.StartDate, req.EndDate)

	// 1. 参数验证
	if req.Uid == "" || req.StartDate == "" || req.EndDate == "" {
		return "", "", fmt.Errorf("参数不完整，需要 uid, start_date, end_date")
	}

	// 2. 调用支付宝接口获取数据
	alipayData, err := l.fetchAlipayData(req.Uid, req.StartDate, req.EndDate)
	if err != nil {
		l.Logger.Errorf("调用支付宝接口失败: %v", err)
		return "", "", fmt.Errorf("获取数据失败: %v", err)
	}

	// 3. 匹配 UID 并生成 Excel
	filePath, filename, err = l.generateExcel(alipayData, req.Uid, req.StartDate, req.EndDate)
	if err != nil {
		l.Logger.Errorf("生成Excel失败: %v", err)
		return "", "", fmt.Errorf("生成Excel失败: %v", err)
	}

	return filePath, filename, nil
}

// fetchAlipayData 调用支付宝接口获取数据
func (l *ZfbDownloadLogic) fetchAlipayData(uid, startDate, endDate string) (*types.AlipayApiResponse, error) {
	// 构造请求体
	requestBody := map[string]interface{}{
		"queries": []map[string]interface{}{
			{
				"chartId": "D202601120016130100044932943",
				"query": map[string]interface{}{
					"select": []map[string]interface{}{
						{"displayId": "d2f62acc-7111-42d2-8588-5a35eb7f2583", "displayName": "日期", "aggregateMethod": "", "columnId": "D2024082600161505000042778079", "dateGranularity": "day", "type": "dimension"},
						{"displayId": "5323949a-9859-4fae-a9c6-07727e8793d6", "displayName": "机构id", "aggregateMethod": "", "columnId": "D2024082600161505000042778084", "type": "dimension"},
						{"displayId": "9ca18e6d-1974-406d-a5e6-2f19b3e0dcba", "displayName": "机构名称", "aggregateMethod": "", "columnId": "D2024082600161505000042778086", "type": "dimension"},
						{"displayId": "70b4f200-4fbf-4f8d-842a-78cd9491c390", "displayName": "推广客uid", "aggregateMethod": "", "columnId": "D2024082600161505000042778081", "type": "dimension"},
						{"displayId": "2a0714b1-7a90-4fc9-ac17-94a03db93228", "displayName": "5SUV", "aggregateMethod": "count_distinct", "columnId": "D2024102800161505000045179491", "type": "measure"},
						{"displayId": "c043f060-c18d-417e-bcd1-9f32f86bc52d", "displayName": "分享人数", "aggregateMethod": "", "columnId": "D2024082600161505000042783737", "type": "measure"},
						{"displayId": "00893fcd-bf09-4349-8393-f80e862cb3b0", "displayName": "人均播放时长", "columnId": "D2024082600161505000042783734", "aggregateMethod": "", "type": "measure"},
						{"displayId": "090f6c71-4f9d-468e-ac7d-965cab38edb6", "displayName": "人均vV数", "columnId": "D2024082600161505000042782892", "expression": "[分享裂变5sVV数]\n/[助力成功征集数]", "aggregateMethod": "", "type": "measure"},
						{"displayId": "ab9ec617-f631-4a7a-94fa-81cba4b2b735", "displayName": "0到1分钟分享人数", "columnId": "D2024082600161505000042783735", "aggregateMethod": "", "type": "measure"},
						{"displayId": "64024caa-ec51-4c95-b317-1c92b2ad66fb", "displayName": "1到2分钟分享人数", "columnId": "D2024082600161505000042783736", "aggregateMethod": "", "type": "measure"},
						{"displayId": "c7e7d34b-9d16-4804-8927-2b3cc3d22a42", "displayName": "2到3分钟分享人数", "columnId": "D2024082600161505000042783738", "aggregateMethod": "", "type": "measure"},
						{"displayId": "4b83263e-42da-4739-ad1b-fed699e0845b", "displayName": "3到4分钟分享人数", "columnId": "D2024082600161505000042783739", "aggregateMethod": "", "type": "measure"},
						{"displayId": "f8683f6f-b0e6-44c5-bcd5-9b0f07634db8", "displayName": "4到5分钟分享人数", "columnId": "D2024082600161505000042783740", "aggregateMethod": "", "type": "measure"},
						{"displayId": "bb811ab1-6868-41dd-ba91-94f444c6830e", "displayName": "5到6分钟分享人数", "columnId": "D2024082600161505000042783741", "aggregateMethod": "", "type": "measure"},
						{"displayId": "a39acb59-753c-47b6-89f6-3043351a6531", "displayName": "6到7分钟分享人数", "columnId": "D2024082600161505000042783742", "aggregateMethod": "", "type": "measure"},
						{"displayId": "cd4e4d92-5473-423d-866d-41cdb8f1ba9d", "displayName": "7到8分钟分享人数", "columnId": "D2024082600161505000042783743", "aggregateMethod": "", "type": "measure"},
						{"displayId": "24fc19a2-0f07-4f6b-8a73-991be7929aa2", "displayName": "8分钟以上分享人数", "columnId": "D2024082600161505000042783744", "aggregateMethod": "", "type": "measure"},
					},
					"where": map[string]interface{}{
						"children": []map[string]interface{}{
							{
								"columnId":          "D2024082600161505000042778079",
								"granularity":       "day",
								"offsetStringValue": "d-0",
								"operator":          ">=",
								"quickGranularity":  "day",
								"right":             startDate + " 00:00:00",
							},
							{
								"columnId":          "D2024082600161505000042778079",
								"granularity":       "day",
								"offsetStringValue": "d-0",
								"operator":          "<=",
								"quickGranularity":  "day",
								"right":             endDate + " 23:59:59",
							},
							{
								"columnId":          "D2024082600161505000042778081",
								"granularity":       "day",
								"offsetStringValue": "d-0",
								"operator":          "in",
								"right":             []string{uid},
							},
						},
						"relation": "and",
					},
				},
				"reportId": "D2026011200161406000006322268",
				"version":  "FORMAL",
			},
		},
		"spreadsheetDefinition": map[string]interface{}{
			"fields": map[string]interface{}{
				"rows": []string{
					"d2f62acc-7111-42d2-8588-5a35eb7f2583",
					"5323949a-9859-4fae-a9c6-07727e8793d6",
					"9ca18e6d-1974-406d-a5e6-21f9b3e0dcba",
					"70b4f200-4bfb-4f8d-842a-78cd9491c390",
				},
				"values": []string{
					"2a0714b1-7a90-4fc9-ac17-94a03db93228",
					"c043f060-c18d-417e-bcd1-9f32f86bc52d",
					"00893fcd-bf09-4349-8393-f80e862cb3b0",
					"090f6c71-4f9d-468e-ac7d-965cab38edb6",
					"ab9ec617-f631-4a7a-94fa-81cba4b2b735",
					"64024caa-ce51-4c95-b317-1c92b2ad66fb",
					"c7e7d34b-9d16-4804-8927-2b3cc3d22a42",
					"4b83263e-42da-4739-ad1b-fed699e0845b",
					"f8683f6f-b0e6-44c5-bcd5-9b0f07634db8",
					"bb811ab1-6868-41dd-ba91-94f444c6830e",
					"a39acb59-753c-47b6-89f6-3043351a6531",
					"cd4e4d92-5473-423d-866d-41cdb8f1ba9d",
					"24fc19a2-0f07-4f6b-8a73-991be7929aa2",
				},
			},
			"tableOrder": map[string]interface{}{
				"globalOrder": []map[string]interface{}{
					{
						"baseDisplayId":  "c043f060-c18d-417e-bcd1-9f32f86bc52d",
						"order":          "DESC",
						"orderDisplayId": "c043f060-c18d-417e-bcd1-9f32f86bc52d",
					},
				},
				"groupOrder": []interface{}{},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 发送HTTP请求
	req, err := http.NewRequest("POST", "https://diopen.alipay.com/dipublic/adame/charts/spread_query_data", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头（Cookie需要从配置或数据库获取）
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", l.getAlipayChookie())

	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp types.AlipayApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回失败: %s", apiResp.ErrorDesc)
	}

	return &apiResp, nil
}

// generateExcel 生成Excel文件
func (l *ZfbDownloadLogic) generateExcel(data *types.AlipayApiResponse, uid, startDate, endDate string) (string, string, error) {
	if len(data.Data.SyncResult) == 0 {
		return "", "", fmt.Errorf("没有数据")
	}

	result := data.Data.SyncResult[0]

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	f.SetSheetName("Sheet1", sheetName)

	// 定义固定的表头顺序（按照图片中的顺序）
	type ColumnDef struct {
		DisplayID string
		Name      string
	}

	columns := []ColumnDef{
		{"d2f62acc-7111-42d2-8588-5a35eb7f2583", "日期"},
		{"5323949a-9859-4fae-a9c6-07727e8793d6", "机构id"},
		{"9ca18e6d-1974-406d-a5e6-2f19b3e0dcba", "机构名称"},
		{"70b4f200-4fbf-4f8d-842a-78cd9491c390", "推广客uid"},
		{"2a0714b1-7a90-4fc9-ac17-94a03db93228", "5SSUV"},
		{"c043f060-c18d-417e-bcd1-9f32f86bc52d", "分享人数"},
		{"00893fcd-bf09-4349-8393-f80e862cb3b0", "人均播放时长"},
		{"090f6c71-4f9d-468e-ac7d-965cab38edb6", "人均uv数"},
		{"ab9ec617-f631-4a7a-94fa-81cba4b2b735", "0到1分钟分享人数"},
		{"64024caa-ec51-4c95-b317-1c92b2ad66fb", "1到2分钟分享人数"},
		{"c7e7d34b-9d16-4804-8927-2b3cc3d22a42", "2到3分钟分享人数"},
		{"4b83263e-42da-4739-ad1b-fed699e0845b", "3到4分钟分享人数"},
		{"f8683f6f-b0e6-44c5-bcd5-9b0f07634db8", "4到5分钟分享人数"},
		{"bb811ab1-6868-41dd-ba91-94f444c6830e", "5到6分钟分享人数"},
		{"a39acb59-753c-47b6-89f6-3043351a6531", "6到7分钟分享人数"},
		{"cd4e4d92-5473-423d-866d-41cdb8f1ba9d", "7到8分钟分享人数"},
		{"24fc19a2-0f07-4f6b-8a73-991be7929aa2", "8分钟以上分享人数"},
	}

	// 写入表头
	for colIdx, col := range columns {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, col.Name)
	}

	// 写入数据行 - 匹配UID
	const uidFieldId = "70b4f200-4fbf-4f8d-842a-78cd9491c390"
	rowIdx := 2
	log.Print("dataValue: ", result.DataValue)
	for _, dataRow := range result.DataValue {
		// 检查UID是否匹配
		if uidValue, ok := dataRow[uidFieldId]; ok && fmt.Sprintf("%v", uidValue) == uid {
			// 写入这一行数据
			for colIdx, col := range columns {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx)
				if value, exists := dataRow[col.DisplayID]; exists {
					f.SetCellValue(sheetName, cell, value)
				}
			}
			rowIdx++
		}
	}

	// 保存文件
	downloadDir := filepath.Join("service", "api", "download_files")
	os.MkdirAll(downloadDir, 0755)

	filename := fmt.Sprintf("alipay_data_%s_%s_%s.xlsx", uid, startDate, endDate)
	filePath := filepath.Join(downloadDir, filename)

	if err := f.SaveAs(filePath); err != nil {
		return "", "", fmt.Errorf("保存Excel失败: %w", err)
	}

	return filePath, filename, nil
}

// getAlipayChookie 获取支付宝Cookie（从数据库的 media_token 表读取）
func (l *ZfbDownloadLogic) getAlipayChookie() string {
	token, _, err := model.GetTokensByMedia(l.svcCtx.DB, "zfb_pachong")
	if err != nil {
		l.Logger.Errorf("获取支付宝Cookie失败: %v", err)
		return ""
	}
	return token
}
