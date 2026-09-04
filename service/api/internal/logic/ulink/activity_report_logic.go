package ulink

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"
)

const taobaoAPIURL = "http://gw.api.taobao.com/router/rest"

// ActivityReportLogic 淘宝客活动报表查询（纯 Go 实现）
type ActivityReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivityReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivityReportLogic {
	return &ActivityReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// pidQueryResult 单个 PID 的双日期查询结果
type pidQueryResult struct {
	pid     string
	current *activityReportItem
	prev    *activityReportItem
}

// activityReportItem 活动报表数据项
type activityReportItem struct {
	Pid           string
	BizDate       string
	Union30dLxUv  int64
	RewardAmount  string
	ExtInfoParsed map[string]interface{}
}

// GetActivityReport 活动报表批量查询（Excel 上传 → Excel 下载）
func (l *ActivityReportLogic) GetActivityReport(w http.ResponseWriter, r *http.Request) {
	cfg := l.svcCtx.Config.Ulink

	// 解析 form 参数
	var req types.ActivityReportReq
	req.BizDate = r.FormValue("biz_date")
	req.QueryType = r.FormValue("query_type")
	req.EventId = r.FormValue("event_id")

	if req.BizDate == "" {
		writeJSONError(w, 400, "biz_date 不能为空（格式: yyyyMMdd）")
		return
	}
	if req.QueryType == "" {
		req.QueryType = "1"
	}
	if req.EventId == "" {
		req.EventId = cfg.DefaultEventId
	}

	// 读取上传的 Excel 文件
	inputPath, cleanup, err := saveUploadedFile(r, "link_file", cfg.TempDir, "*.xlsx")
	if err != nil {
		writeJSONError(w, 400, "读取上传文件失败: "+err.Error())
		return
	}
	defer cleanup()

	// 解析 Excel 获取 PID 列表
	pids, err := readPIDsFromExcel(inputPath)
	if err != nil {
		writeJSONError(w, 400, "解析 Excel 文件失败: "+err.Error())
		return
	}
	if len(pids) == 0 {
		writeJSONError(w, 400, "Excel 文件中没有找到 pid 数据")
		return
	}

	// 计算前一天日期
	prevDate, err := getPreviousDate(req.BizDate)
	if err != nil {
		writeJSONError(w, 400, "日期格式错误: "+err.Error())
		return
	}

	queryType, _ := strconv.Atoi(req.QueryType)

	// 并发查询每个 PID
	ctx := r.Context()
	results := make([]pidQueryResult, len(pids))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 最大并发 5

	for i, pid := range pids {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
				defer func() { <-sem }()
			}

			curData := queryActivityReport(ctx, cfg.TaobaoAppKey, cfg.TaobaoAppSecret, req.EventId, req.BizDate, queryType, p)
			prevData := queryActivityReport(ctx, cfg.TaobaoAppKey, cfg.TaobaoAppSecret, req.EventId, prevDate, queryType, p)

			results[idx] = pidQueryResult{pid: p, current: curData, prev: prevData}
		}(i, pid)
	}
	wg.Wait()

	// 生成 Excel 输出
	outputPath := filepath.Join(cfg.TempDir, fmt.Sprintf("activity_report_%s_%d.xlsx", req.BizDate, time.Now().UnixNano()))
	defer os.Remove(outputPath)

	if err := buildActivityReportExcel(outputPath, results, req.EventId, req.BizDate); err != nil {
		l.Errorf("生成报表 Excel 失败: %v", err)
		writeJSONError(w, 500, "生成报表失败: "+err.Error())
		return
	}

	serveExcelFile(w, outputPath, fmt.Sprintf("activity_report_%s.xlsx", req.BizDate))
}

// ==================== 淘宝客 API 调用 ====================

// signTaobaoAPI 生成淘宝 API 签名（MD5，ASCII 排序）
func signTaobaoAPI(params map[string]string, appSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	sb := appSecret
	for _, k := range keys {
		sb += k + params[k]
	}
	sb += appSecret

	h := md5.Sum([]byte(sb))
	return strings.ToUpper(fmt.Sprintf("%x", h))
}

// queryActivityReport 调用淘宝客活动报表 API，返回第一条数据
func queryActivityReport(ctx context.Context, appKey, appSecret, eventId, bizDate string, queryType int, pid string) *activityReportItem {
	params := map[string]string{
		"method":      "taobao.tbk.dg.cpa.activity.report",
		"app_key":     appKey,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"format":      "json",
		"v":           "2.0",
		"sign_method": "md5",
		"partner_id":  "top-apitools",
		"event_id":    eventId,
		"biz_date":    bizDate,
		"query_type":  strconv.Itoa(queryType),
		"page_no":     "1",
		"page_size":   "10",
	}
	if pid != "" {
		params["pid"] = pid
	}
	params["sign"] = signTaobaoAPI(params, appSecret)

	queryValues := url.Values{}
	for k, v := range params {
		queryValues.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, taobaoAPIURL+"?"+queryValues.Encode(), nil)
	if err != nil {
		return nil
	}
	resp, err := taobaoHTTPClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}

	if _, hasErr := result["error_response"]; hasErr {
		return nil
	}

	respData, ok := result["tbk_dg_cpa_activity_report_response"].(map[string]interface{})
	if !ok {
		return nil
	}

	resultData, ok := respData["result"].(map[string]interface{})
	if !ok {
		return nil
	}

	dataMap, ok := resultData["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	resultsMap, ok := dataMap["results"].(map[string]interface{})
	if !ok {
		return nil
	}

	dtoList, ok := resultsMap["vegas_cpa_report_d_t_o"].([]interface{})
	if !ok || len(dtoList) == 0 {
		return nil
	}

	dto, ok := dtoList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	item := &activityReportItem{
		Pid:          getString(dto, "pid"),
		BizDate:      getString(dto, "biz_date"),
		RewardAmount: getString(dto, "reward_amount"),
	}
	if v, ok := dto["union_30d_lx_uv"]; ok {
		item.Union30dLxUv = toInt64(v)
	}

	extInfoStr := getString(dto, "ext_info")
	item.ExtInfoParsed = parseExtInfo(extInfoStr, eventId)

	return item
}

// parseExtInfo 解析 ext_info JSON 字符串，按活动类型提取字段
func parseExtInfo(extInfoStr, eventId string) map[string]interface{} {
	result := map[string]interface{}{}
	if extInfoStr == "" {
		return result
	}

	var extData map[string]interface{}
	if err := json.Unmarshal([]byte(extInfoStr), &extData); err != nil {
		return result
	}

	switch eventId {
	case "7610424": // 福利购
		result["crowd1_reward_uv"] = getFirstNonEmpty(extData, "人群1结算奖励uv", "人群1预估奖励uv")
		result["crowd2_reward_uv"] = getFirstNonEmpty(extData, "人群2结算奖励uv", "人群2预估奖励uv")
		result["crowd3_reward_uv"] = getFirstNonEmpty(extData, "人群3结算奖励uv", "人群3预估奖励uv")
		result["crowd4_reward_uv"] = getFirstNonEmpty(extData, "人群4结算奖励uv", "人群4预估奖励uv")
		result["crowd5_reward_uv"] = getFirstNonEmpty(extData, "人群5结算奖励uv", "人群5预估奖励uv")
		result["account_draw_rate"] = extData["账号总开奖率（奖励计算用)"]
		result["draw_rate"] = extData["开奖率"]
		result["update_time"] = extData["更新时间"]

	case "4297311": // 超级红包
		result["user_quality_level"] = extData["用户质量等级"]
		result["account_draw_rate"] = extData["账号总开奖率（奖励计算用)"]
		result["settlement_reward_uv"] = getFirstNonEmpty(extData, "结算奖励uv", "预估奖励uv")
		result["draw_rate"] = extData["开奖率"]
		result["update_time"] = extData["更新时间"]
	}

	return result
}

// buildActivityReportExcel 生成活动报表 Excel 文件
func buildActivityReportExcel(outputPath string, results []pidQueryResult, eventId, curDate string) error {
	f := excelize.NewFile()
	sheet := "活动报表"
	f.SetSheetName("Sheet1", sheet)

	// 构建表头
	var headers []string
	headers = append(headers, "pid", "统计日期", "符合奖励要求的累计用户数", "奖励金额")

	switch eventId {
	case "7610424": // 福利购
		headers = append(headers,
			"人群1奖励uv", "人群2奖励uv", "人群3奖励uv", "人群4奖励uv", "人群5奖励uv",
			"账号总开奖率", "开奖率", "更新时间",
			"累计用户数(日增)", "奖励金额(日增)",
			"人群1奖励uv(日增)", "人群2奖励uv(日增)", "人群3奖励uv(日增)",
			"人群4奖励uv(日增)", "人群5奖励uv(日增)",
			"状态",
		)
	case "4297311": // 超级红包
		headers = append(headers,
			"用户质量等级", "奖励uv", "账号总开奖率", "开奖率", "更新时间",
			"累计用户数(日增)", "奖励金额(日增)", "状态",
		)
	default:
		headers = append(headers, "状态")
	}

	// 写表头
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 写数据行
	for rowIdx, r := range results {
		row := rowIdx + 2
		data := map[string]interface{}{"pid": r.pid, "统计日期": curDate}

		if r.current == nil {
			data["状态"] = "No data"
			setRow(f, sheet, row, headers, data)
			continue
		}

		data["符合奖励要求的累计用户数"] = r.current.Union30dLxUv
		data["奖励金额"] = r.current.RewardAmount
		data["状态"] = "成功"

		ext := r.current.ExtInfoParsed

		switch eventId {
		case "7610424": // 福利购
			data["人群1奖励uv"] = ext["crowd1_reward_uv"]
			data["人群2奖励uv"] = ext["crowd2_reward_uv"]
			data["人群3奖励uv"] = ext["crowd3_reward_uv"]
			data["人群4奖励uv"] = ext["crowd4_reward_uv"]
			data["人群5奖励uv"] = ext["crowd5_reward_uv"]
			data["账号总开奖率"] = ext["account_draw_rate"]
			data["开奖率"] = ext["draw_rate"]
			data["更新时间"] = ext["update_time"]

			if r.prev != nil {
				data["累计用户数(日增)"] = safeSubtractInt64(r.current.Union30dLxUv, r.prev.Union30dLxUv)
				data["奖励金额(日增)"] = safeSubtractStr(r.current.RewardAmount, r.prev.RewardAmount)
				prevExt := r.prev.ExtInfoParsed
				data["人群1奖励uv(日增)"] = safeSubtractInterface(ext["crowd1_reward_uv"], prevExt["crowd1_reward_uv"])
				data["人群2奖励uv(日增)"] = safeSubtractInterface(ext["crowd2_reward_uv"], prevExt["crowd2_reward_uv"])
				data["人群3奖励uv(日增)"] = safeSubtractInterface(ext["crowd3_reward_uv"], prevExt["crowd3_reward_uv"])
				data["人群4奖励uv(日增)"] = safeSubtractInterface(ext["crowd4_reward_uv"], prevExt["crowd4_reward_uv"])
				data["人群5奖励uv(日增)"] = safeSubtractInterface(ext["crowd5_reward_uv"], prevExt["crowd5_reward_uv"])
			}

		case "4297311": // 超级红包
			data["用户质量等级"] = ext["user_quality_level"]
			data["奖励uv"] = ext["settlement_reward_uv"]
			data["账号总开奖率"] = ext["account_draw_rate"]
			data["开奖率"] = ext["draw_rate"]
			data["更新时间"] = ext["update_time"]

			if r.prev != nil {
				data["累计用户数(日增)"] = safeSubtractInt64(r.current.Union30dLxUv, r.prev.Union30dLxUv)
				data["奖励金额(日增)"] = safeSubtractStr(r.current.RewardAmount, r.prev.RewardAmount)
			}
		}

		setRow(f, sheet, row, headers, data)
	}

	return f.SaveAs(outputPath)
}

// setRow 按表头顺序写入一行数据
func setRow(f *excelize.File, sheet string, row int, headers []string, data map[string]interface{}) {
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, row)
		if v, ok := data[h]; ok && v != nil {
			f.SetCellValue(sheet, cell, v)
		}
	}
}

// ==================== Excel 读取 ====================

// readPIDsFromExcel 从 Excel 文件中读取 pid 列
func readPIDsFromExcel(filePath string) ([]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 文件没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel 文件无数据")
	}

	// 找到 pid 列
	pidColIdx := -1
	for i, header := range rows[0] {
		if strings.EqualFold(strings.TrimSpace(header), "pid") {
			pidColIdx = i
			break
		}
	}
	if pidColIdx == -1 {
		return nil, fmt.Errorf("未找到 'pid' 列")
	}

	var pids []string
	for _, row := range rows[1:] {
		if pidColIdx < len(row) {
			pid := strings.TrimSpace(row[pidColIdx])
			if pid != "" {
				pids = append(pids, pid)
			}
		}
	}

	return pids, nil
}

// ==================== 工具函数 ====================

func getPreviousDate(dateStr string) (string, error) {
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}
	return t.AddDate(0, 0, -1).Format("20060102"), nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func toInt64(v interface{}) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	}
	return 0
}

func getFirstNonEmpty(m map[string]interface{}, keys ...string) interface{} {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil && v != "" {
			return v
		}
	}
	return ""
}

func safeSubtractInt64(a, b int64) int64 {
	return a - b
}

func safeSubtractStr(a, b string) float64 {
	fa, _ := strconv.ParseFloat(a, 64)
	fb, _ := strconv.ParseFloat(b, 64)
	return fa - fb
}

func safeSubtractInterface(a, b interface{}) interface{} {
	fa, errA := toFloat(a)
	fb, errB := toFloat(b)
	if errA != nil || errB != nil {
		return ""
	}
	diff := fa - fb
	if diff == float64(int64(diff)) {
		return int64(diff)
	}
	return fmt.Sprintf("%.2f", diff)
}

func toFloat(v interface{}) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("nil")
	}
	switch x := v.(type) {
	case float64:
		return x, nil
	case int64:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(x, 64)
	}
	return 0, fmt.Errorf("unknown type")
}
