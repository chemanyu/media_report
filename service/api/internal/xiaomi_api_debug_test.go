package internal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

// XiaomiAPIClient 小米营销API客户端
type XiaomiAPIClient struct {
	BaseURL    string
	SignId     string
	SecretKey  string
	CustomerId int64
	HTTPClient *http.Client
}

// NewXiaomiAPIClient 创建小米API客户端
func NewXiaomiAPIClient(signId, secretKey string, customerId int64) *XiaomiAPIClient {
	return &XiaomiAPIClient{
		BaseURL:    "https://api.e.mi.com",
		SignId:     signId,
		SecretKey:  secretKey,
		CustomerId: customerId,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ReportDataRequest 报表数据请求参数
type ReportDataRequest struct {
	CustomerId     int64   `json:"customerId"`     // 账户id（必填）
	SDate          string  `json:"sdate"`          // 开始日期 yyyy-MM-dd（必填）
	EDate          string  `json:"edate"`          // 结束日期 yyyy-MM-dd（必填）
	MetricsList    string  `json:"metricsList"`    // 指标列表，多个用英文逗号隔开
	DimensionsList string  `json:"dimensionsList"` // 维度列表，多个用英文逗号隔开
	CampaignIds    []int64 `json:"campaignIds"`    // 计划id列表
	AdReportIds    []int64 `json:"adReportIds"`    // 广告id列表
	GroupBy        string  `json:"groupBy"`        // 分组类型
	OrderBy        string  `json:"orderBy"`        // 排序指标
	SortMode       int     `json:"sortMode"`       // 排序方式：0倒序，1正序
	Page           int     `json:"page"`           // 页数，默认1
	PageSize       int     `json:"pagesize"`       // 页面大小，默认20，最大1000
}

// ReportDataResponse 报表数据响应
type ReportDataResponse struct {
	List       []ReportData `json:"list"`       // 分日数据
	Total      TotalData    `json:"total"`      // 汇总数据
	RecordDate string       `json:"recordDate"` // 日期
	Page       int          `json:"page"`       // 页码
	PageSize   int          `json:"pagesize"`   // 页面大小
	TotalCount int          `json:"totalCount"` // 总数
}

// ReportData 报表数据
type ReportData struct {
	RecordDate        string  `json:"recordDate"`        // 日期
	CustomerId        int64   `json:"customerId"`        // 账户id
	CampaignName      string  `json:"campaignName"`      // 计划名称
	CampaignId        int64   `json:"campaignId"`        // 计划id
	AdReportName      string  `json:"adReportName"`      // 广告名称
	AdReportId        int64   `json:"adReportId"`        // 广告id
	Cost              float64 `json:"cost"`              // 消耗（元）
	ExposeNum         int64   `json:"exposeNum"`         // 曝光量
	ClickNum          int64   `json:"clickNum"`          // 点击量
	DownloadNum       int64   `json:"downloadNum"`       // 下载量
	CPC               float64 `json:"cpc"`               // 点击均价（元）
	CTR               float64 `json:"ctr"`               // 点击率
	ECPM              float64 `json:"ecpm"`              // ECPM
	CPD               float64 `json:"cpd"`               // 下载均价（元）
	DTR               float64 `json:"dtr"`               // 下载率
	CashCost          float64 `json:"cashCost"`          // 现金消耗（元）
	VirtualCost       float64 `json:"virtualCost"`       // 虚拟金消耗（元）
	Reservation       int64   `json:"reservation"`       // APP预约
	CancelReservation int64   `json:"cancelReservation"` // APP取消预约
	RealReservation   int64   `json:"realReservation"`   // APP实际预约量
	RPA               float64 `json:"rpa"`               // APP预约成本（元）
	DisplayType       string  `json:"displayType"`       // 版位
	PlatformType      string  `json:"platformType"`      // 流量范围
}

// TotalData 汇总数据
type TotalData struct {
	Cost        float64 `json:"cost"`        // 消耗（元）
	ExposeNum   int64   `json:"exposeNum"`   // 曝光量
	ClickNum    int64   `json:"clickNum"`    // 点击量
	DownloadNum int64   `json:"downloadNum"` // 下载量
	CPC         float64 `json:"cpc"`         // 点击均价（元）
	CTR         float64 `json:"ctr"`         // 点击率
	ECPM        float64 `json:"ecpm"`        // ECPM
}

// MetricsParamExplainRequest 指标参数说明请求
type MetricsParamExplainRequest struct {
	CustomerId  int64 `json:"customerId"`  // 账户id（必填）
	DataTopicId int   `json:"dataTopicId"` // 数据主题id（必填）1-广告创意数据，4-关键词数据，2-素材/广告数据
}

// MetricsParamExplainResponse 指标参数说明响应
type MetricsParamExplainResponse struct {
	DataTopicId   int        `json:"dataTopicId"`   // 数据主题ID
	DataTopicName string     `json:"dataTopicName"` // 数据主题名称
	Dimensions    []ParamDef `json:"dimensions"`    // 维度列表
	Metrics       []ParamDef `json:"metrics"`       // 指标列表
}

// ParamDef 参数定义
type ParamDef struct {
	Code        int    `json:"code"`        // 维度/指标code
	Description string `json:"description"` // 维度/指标描述
	Name        string `json:"name"`        // 维度/指标名称
}

// calculateSign 计算签名
func (c *XiaomiAPIClient) calculateSign(params map[string]string) string {
	fmt.Printf("\n[签名计算] 开始计算签名...")
	fmt.Printf("[签名计算] 原始参数: %+v\n", params)

	// 1. 将参数按key排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	fmt.Printf("[签名计算] 排序后的key列表: %v\n", keys)

	// 2. 拼接成字符串 key1value1key2value2...
	var paramStr strings.Builder
	for _, k := range keys {
		paramStr.WriteString(k)
		paramStr.WriteString(params[k])
	}
	paramStrValue := paramStr.String()
	fmt.Printf("[签名计算] 参数拼接字符串: %s\n", paramStrValue)

	// 3. 拼接 signId_key + signId_value + paramStr + secretKey_value
	signStr := paramStrValue + c.SecretKey
	fmt.Printf("[签名计算] 完整签名字符串: signId%s%s%s\n", c.SignId, paramStrValue, c.SecretKey)
	fmt.Printf("[签名计算] 完整签名字符串长度: %d\n", len(signStr))

	// 4. MD5加密
	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])
	fmt.Printf("[签名计算] MD5签名结果: %s\n\n", sign)
	return sign
}

// GetReportData 获取报表数据
func (c *XiaomiAPIClient) GetReportData(req *ReportDataRequest) (*ReportDataResponse, error) {
	fmt.Printf("\n========== 开始调用小米API ==========")
	fmt.Printf("[请求参数] CustomerId: %d\n", req.CustomerId)
	fmt.Printf("[请求参数] 日期范围: %s 至 %s\n", req.SDate, req.EDate)

	// 构建请求参数
	params := make(map[string]string)
	params["signId"] = c.SignId
	params["customerId"] = fmt.Sprintf("%d", req.CustomerId)
	params["sdate"] = req.SDate
	params["edate"] = req.EDate

	// 可选参数
	if req.MetricsList != "" {
		params["metricsList"] = req.MetricsList
		fmt.Printf("[请求参数] MetricsList: %s\n", req.MetricsList)
	}
	if req.DimensionsList != "" {
		params["dimensionsList"] = req.DimensionsList
		fmt.Printf("[请求参数] DimensionsList: %s\n", req.DimensionsList)
	}
	if req.OrderBy != "" {
		params["orderBy"] = req.OrderBy
	}
	if req.SortMode > 0 {
		params["sortMode"] = fmt.Sprintf("%d", req.SortMode)
	}
	if req.Page > 0 {
		params["page"] = fmt.Sprintf("%d", req.Page)
	} else {
		params["page"] = "1"
	}
	if req.PageSize > 0 {
		params["pagesize"] = fmt.Sprintf("%d", req.PageSize)
	} else {
		params["pagesize"] = "20"
	}
	fmt.Printf("[请求参数] Page: %s, PageSize: %s\n", params["page"], params["pagesize"])

	// 计算签名（不包含signId和sign）
	paramsForSign := make(map[string]string)
	for k, v := range params {
		paramsForSign[k] = v
	}
	sign := c.calculateSign(paramsForSign)

	// 构建URL - 先拼接其他参数，signId和sign放在最后
	reqURL := fmt.Sprintf("%s/openapi/v5/report/getData", c.BaseURL)

	// 手动拼接查询参数（不使用编码，直接使用原值）
	var queryParams []string
	for k, v := range params {
		if k != "signId" {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// 拼接完整URL，确保signId和sign在最后
	var fullURL string
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s&signId=%s&sign=%s", reqURL, strings.Join(queryParams, "&"), c.SignId, sign)
	} else {
		fullURL = fmt.Sprintf("%s?signId=%s&sign=%s", reqURL, c.SignId, sign)
	}

	fmt.Printf("\n[HTTP请求] 完整URL: %s\n", fullURL)
	fmt.Printf("[HTTP请求] 开始发送GET请求...\n")

	// 发送请求
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		fmt.Printf("[HTTP请求] ✗ 请求失败: %v\n", err)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("[HTTP响应] ✓ 状态码: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("[HTTP响应] Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[HTTP响应] ✗ 读取响应失败: %v\n", err)
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	fmt.Printf("[HTTP响应] 响应体长度: %d bytes\n", len(body))
	fmt.Printf("[HTTP响应] 响应内容: %s\n", string(body))

	// 解析响应
	var result ReportDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[数据解析] ✗ JSON解析失败: %v\n", err)
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}
	fmt.Printf("[数据解析] ✓ 解析成功，共 %d 条数据\n", len(result.List))
	fmt.Printf("========== API调用完成 ==========")

	return &result, nil
}

// GetMetricsParamExplain 获取指标参数说明
func (c *XiaomiAPIClient) GetMetricsParamExplain(req *MetricsParamExplainRequest) (*MetricsParamExplainResponse, error) {
	fmt.Printf("\n========== 开始调用小米API - 获取指标参数说明 ==========")
	fmt.Printf("[请求参数] CustomerId: %d\n", req.CustomerId)
	fmt.Printf("[请求参数] DataTopicId: %d\n", req.DataTopicId)

	// 构建请求参数
	params := make(map[string]string)
	params["signId"] = c.SignId
	params["customerId"] = fmt.Sprintf("%d", req.CustomerId)
	params["dataTopicId"] = fmt.Sprintf("%d", req.DataTopicId)

	// 计算签名
	paramsForSign := make(map[string]string)
	for k, v := range params {
		paramsForSign[k] = v
	}
	sign := c.calculateSign(paramsForSign)

	// 构建URL - 先拼接其他参数，signId和sign放在最后
	reqURL := fmt.Sprintf("%s/openapi/v5/report/getMetricsParamExplain", c.BaseURL)

	// 手动拼接查询参数（不使用编码，直接使用原值）
	var queryParams []string
	for k, v := range params {
		if k != "signId" {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// 拼接完整URL，确保signId和sign在最后
	var fullURL string
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s&signId=%s&sign=%s", reqURL, strings.Join(queryParams, "&"), c.SignId, sign)
	} else {
		fullURL = fmt.Sprintf("%s?signId=%s&sign=%s", reqURL, c.SignId, sign)
	}

	fmt.Printf("\n[HTTP请求] 完整URL: %s\n", fullURL)
	fmt.Printf("[HTTP请求] 开始发送GET请求...\n")

	// 发送请求
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		fmt.Printf("[HTTP请求] ✗ 请求失败: %v\n", err)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("[HTTP响应] ✓ 状态码: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("[HTTP响应] Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[HTTP响应] ✗ 读取响应失败: %v\n", err)
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	fmt.Printf("[HTTP响应] 响应体长度: %d bytes\n", len(body))
	fmt.Printf("[HTTP响应] 响应内容: %s\n", string(body))

	// 解析响应
	var result MetricsParamExplainResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[数据解析] ✗ JSON解析失败: %v\n", err)
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}
	fmt.Printf("[数据解析] ✓ 解析成功，维度数: %d, 指标数: %d\n", len(result.Dimensions), len(result.Metrics))
	fmt.Printf("========== API调用完成 ==========")

	return &result, nil
}

// TestXiaomiAPI 测试小米API
func TestXiaomiAPI(t *testing.T) {

	fmt.Printf("开始测试小米营销API")

	// TODO: 替换为实际的signId、secretKey和customerId
	signId := "229d841af7a5fc79b66ee6a47e11eee3" // 从小米营销平台获取
	secretKey := "xtLqxGLlsNMhhKLr"              // 从小米营销平台获取
	customerId := int64(1261482)                 // 账户ID

	fmt.Printf("\n[初始化] SignId: %s", signId)
	fmt.Printf("[初始化] SecretKey: %s (隐藏)", strings.Repeat("*", len(secretKey)))
	fmt.Printf("[初始化] CustomerId: %d", customerId)

	client := NewXiaomiAPIClient(signId, secretKey, customerId)
	fmt.Printf("[初始化] ✓ API客户端创建成功")

	// 构建请求参数
	today := time.Now().Format("2006-01-02")
	req := &ReportDataRequest{
		CustomerId:     customerId,
		SDate:          today,
		EDate:          today,
		MetricsList:    "2017,2012,1863,2018",
		DimensionsList: "",
		Page:           1,
		PageSize:       100,
	}

	fmt.Printf("\n[请求配置] customerId=%d, sdate=%s, edate=%s", req.CustomerId, req.SDate, req.EDate)
	fmt.Printf("[请求配置] metricsList=%s", req.MetricsList)
	fmt.Printf("[请求配置] dimensionsList=%s", req.DimensionsList)

	// 调用API
	fmt.Printf("\n[API调用] 开始调用GetReportData...")
	result, err := client.GetReportData(req)
	if err != nil {
		t.Fatalf("[API调用] ✗ 调用失败: %v", err)
	}
	fmt.Printf("[API调用] ✓ 调用成功")

	// 输出结果
	fmt.Printf("\n[结果统计] 获取到 %d 条数据", len(result.List))
	fmt.Printf("[结果统计] 总页数: %d, 当前页: %d", result.TotalCount, result.Page)

	// 输出汇总数据
	fmt.Printf("\n=== 汇总数据 ===")
	fmt.Printf("总消耗: %.2f 元", result.Total.Cost)
	fmt.Printf("总曝光: %d", result.Total.ExposeNum)
	fmt.Printf("总点击: %d", result.Total.ClickNum)
	fmt.Printf("总下载: %d", result.Total.DownloadNum)
	fmt.Printf("CPC: %.2f 元, CTR: %.4f%%", result.Total.CPC, result.Total.CTR*100)

	// 输出前5条详细数据
	fmt.Printf("\n=== 详细数据（前5条）===")
	for i, data := range result.List {
		if i >= 5 {
			break
		}
		fmt.Printf("\n第 %d 条:", i+1)
		fmt.Printf("  日期: %s", data.RecordDate)
		fmt.Printf("  计划: %s (ID: %d)", data.CampaignName, data.CampaignId)
		fmt.Printf("  广告: %s (ID: %d)", data.AdReportName, data.AdReportId)
		fmt.Printf("  消耗: %.2f 元", data.Cost)
		fmt.Printf("  曝光: %d, 点击: %d, 下载: %d", data.ExposeNum, data.ClickNum, data.DownloadNum)
		fmt.Printf("  CPC: %.2f 元, CTR: %.4f%%, CPD: %.2f 元", data.CPC, data.CTR*100, data.CPD)
	}

	fmt.Printf("测试完成")

}

// TestCalculateSign 测试签名计算
func TestCalculateSign(t *testing.T) {

	fmt.Printf("测试签名计算算法")

	client := NewXiaomiAPIClient("12345678", "00000000", 321)

	// 示例参数（来自文档）
	params := map[string]string{
		"customerId": "321",
		"edate":      "2018-08-05",
		"orderby":    "viewSum",
		"sdate":      "2018-08-05",
		"signId":     "12345678",
		"sortMode":   "1",
	}

	fmt.Printf("\n[测试数据] 使用文档示例参数:")
	for k, v := range params {
		fmt.Printf("  %s = %s", k, v)
	}

	sign := client.calculateSign(params)

	// 根据文档示例计算
	// paramStr应该是: customerId321edate2018-08-05orderbyviewSumsdate2018-08-05signId12345678sortMode1
	// 完整字符串: signId12345678customerId321edate2018-08-05orderbyviewSumsdate2018-08-05signId12345678sortMode100000000
	// MD5结果应该是: 8054b81bd0b454a89846e04f340bf0de5

	fmt.Printf("\n[验证结果] 计算的签名: %s", sign)
	expectedSign := "8054b81bd0b454a89846e04f340bf0de5"
	fmt.Printf("[验证结果] 期望的签名: %s", expectedSign)

	if sign == expectedSign {
		fmt.Printf("\n✓ ✓ ✓ 签名计算正确！")
	} else {
		fmt.Printf("\n✗ ✗ ✗ 签名计算有误！")
		fmt.Printf("  期望值: %s", expectedSign)
		fmt.Printf("  实际值: %s", sign)
		fmt.Printf("  差异: 请检查签名算法实现")
	}
}

// TestGetMetricsParamExplain 测试获取指标参数说明
func TestGetMetricsParamExplain(t *testing.T) {

	fmt.Printf("测试获取指标参数说明API")

	// TODO: 替换为实际的signId、secretKey和customerId
	signId := "229d841af7a5fc79b66ee6a47e11eee3" // 从小米营销平台获取
	secretKey := "xtLqxGLlsNMhhKLr"              // 从小米营销平台获取
	customerId := int64(1261482)                 // 账户ID

	fmt.Printf("\n[初始化] SignId: %s", signId)
	fmt.Printf("[初始化] SecretKey: %s (隐藏)", strings.Repeat("*", len(secretKey)))
	fmt.Printf("[初始化] CustomerId: %d", customerId)

	client := NewXiaomiAPIClient(signId, secretKey, customerId)
	fmt.Printf("[初始化] ✓ API客户端创建成功")

	// 测试不同的数据主题
	dataTopicIds := []int{1} // 1-广告创意数据，2-素材/广告数据，4-关键词数据
	dataTopicNames := map[int]string{
		1: "广告创意数据",
		2: "素材/广告数据",
		4: "关键词数据",
	}

	for _, topicId := range dataTopicIds {
		fmt.Printf("\n[测试数据主题] %d - %s", topicId, dataTopicNames[topicId])

		req := &MetricsParamExplainRequest{
			CustomerId:  customerId,
			DataTopicId: topicId,
		}

		// 调用API
		fmt.Printf("[API调用] 开始调用GetMetricsParamExplain...")
		result, err := client.GetMetricsParamExplain(req)
		if err != nil {
			t.Errorf("[API调用] ✗ 调用失败: %v", err)
			continue
		}
		fmt.Printf("[API调用] ✓ 调用成功")

		// 输出结果
		fmt.Printf("\n[结果统计] 数据主题: %s (ID: %d)", result.DataTopicName, result.DataTopicId)
		fmt.Printf("[结果统计] 维度数量: %v", result.Dimensions)
		fmt.Printf("[结果统计] 指标数量: %v", result.Metrics)

	}

	fmt.Printf("测试完成")

}
