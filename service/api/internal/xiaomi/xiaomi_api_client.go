package xiaomi

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
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
	CustomerId     int64  `json:"customerId"`     // 账户id（必填）
	SDate          string `json:"sdate"`          // 开始日期 yyyy-MM-dd（必填）
	EDate          string `json:"edate"`          // 结束日期 yyyy-MM-dd（必填）
	MetricsList    string `json:"metricsList"`    // 指标列表，多个用英文逗号隔开
	DimensionsList string `json:"dimensionsList"` // 维度列表，多个用英文逗号隔开
	Page           int    `json:"page"`           // 页数，默认1
	PageSize       int    `json:"pagesize"`       // 页面大小，默认20，最大1000
}

// ReportDataResponse 报表数据响应
type ReportDataResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Result  Result `json:"result"`
}

// Result 结果数据
type Result struct {
	Conf  ConfData   `json:"conf"`
	List  []ListData `json:"list"`
	Total TotalData  `json:"total"`
}

// ConfData 配置数据
type ConfData struct {
	Page     int `json:"page"`
	PageSize int `json:"pagesize"`
	Total    int `json:"total"`
}

// ListData 列表数据
type ListData struct {
	RecordDate             string `json:"recordDate"`             // 日期
	CustomerId             int64  `json:"customerId"`             // 账户id
	Cost                   string `json:"cost"`                   // 消耗（元）
	ExposeNum              int64  `json:"exposeNum"`              // 曝光量
	ClickNum               int64  `json:"clickNum"`               // 点击量
	DownloadNum            int64  `json:"downloadNum"`            // 下载量
	CPC                    string `json:"cpc"`                    // 点击均价（元）
	CTR                    string `json:"ctr"`                    // 点击率
	ECPM                   string `json:"ecpm"`                   // ECPM
	CPD                    string `json:"cpd"`                    // 下载均价（元）
	DTR                    string `json:"dtr"`                    // 下载率
	CashCost               string `json:"cashCost"`               // 现金消耗（元）
	VirtualCost            string `json:"virtualCost"`            // 虚拟金消耗（元）
	PaySumjFormat          string `json:"paySumjFormat"`          // 付费数
	CostPerPayzFormat      string `json:"costPerPayzFormat"`      // 付费均价
	CostPerReActivejFormat string `json:"costPerReActivejFormat"` // 首次拉活均价
	ReActiveSumjFormat     string `json:"reActiveSumjFormat"`     // 首次拉活
	CampaignId             int64  `json:"campaignId"`             // 计划id
	AdReportId             int64  `json:"adReportId"`             // 广告id
}

// TotalData 汇总数据
type TotalData struct {
	RecordDate             string `json:"recordDate"`             // 日期范围
	CustomerId             int64  `json:"customerId"`             // 账户id
	Cost                   string `json:"cost"`                   // 消耗（元）
	ExposeNum              int64  `json:"exposeNum"`              // 曝光量
	ClickNum               int64  `json:"clickNum"`               // 点击量
	DownloadNum            int64  `json:"downloadNum"`            // 下载量
	CPC                    string `json:"cpc"`                    // 点击均价（元）
	CTR                    string `json:"ctr"`                    // 点击率
	ECPM                   string `json:"ecpm"`                   // ECPM
	CPD                    string `json:"cpd"`                    // 下载均价（元）
	DTR                    string `json:"dtr"`                    // 下载率
	CashCost               string `json:"cashCost"`               // 现金消耗（元）
	VirtualCost            string `json:"virtualCost"`            // 虚拟金消耗（元）
	PaySumjFormat          string `json:"paySumjFormat"`          // 付费数
	ReActiveSumjFormat     string `json:"reActiveSumjFormat"`     // 首次拉活数
	CostPerPayzFormat      string `json:"costPerPayzFormat"`      // 付费均价
	CostPerReActivejFormat string `json:"costPerReActivejFormat"` // 首次拉活均价
}

// calculateSign 计算签名
func (c *XiaomiAPIClient) calculateSign(params map[string]string) string {
	// 1. 将参数按key排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 拼接成字符串 key1value1key2value2...
	var paramStr strings.Builder
	for _, k := range keys {
		paramStr.WriteString(k)
		paramStr.WriteString(params[k])
	}
	paramStrValue := paramStr.String()

	// 3. 拼接完整签名字符串并加密
	signStr := paramStrValue + c.SecretKey

	// 4. MD5加密
	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])
	return sign
}

// GetReportData 获取报表数据
func (c *XiaomiAPIClient) GetReportData(req *ReportDataRequest) (*ReportDataResponse, error) {
	// 构建请求参数
	params := make(map[string]string)
	params["signId"] = c.SignId
	params["customerId"] = fmt.Sprintf("%d", req.CustomerId)
	params["sdate"] = req.SDate
	params["edate"] = req.EDate

	// 可选参数
	if req.MetricsList != "" {
		params["metricsList"] = req.MetricsList
	}
	if req.DimensionsList != "" {
		params["dimensionsList"] = req.DimensionsList
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

	// 计算签名
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

	// 发送请求
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result ReportDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}

	return &result, nil
}
