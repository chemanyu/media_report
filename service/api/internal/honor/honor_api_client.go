package honor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	tokenURL  = "https://iam.developer.hihonor.com/auth/realms/advertisement/protocol/openid-connect/token"
	reportURL = "https://ads.cloud.honor.com/openapi/v2_2/support/ad-report/advertiser"
)

// HonorAPIClient Honor广告API客户端
type HonorAPIClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

// NewHonorAPIClient 创建Honor API客户端
func NewHonorAPIClient(clientID, clientSecret string) *HonorAPIClient {
	return &HonorAPIClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

// getAccessToken 获取或刷新 access token（有缓存）
func (c *HonorAPIClient) getAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)

	resp, err := c.httpClient.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("获取token请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取token响应失败: %w", err)
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("解析token响应失败: %w, body: %s", err, string(body))
	}
	if tr.Error != "" {
		return "", fmt.Errorf("获取token失败: %s - %s", tr.Error, tr.ErrorDesc)
	}

	c.accessToken = tr.AccessToken
	// 提前60秒过期，防止边界问题
	c.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn-60) * time.Second)

	return c.accessToken, nil
}

// ReportRequest 广告主报表查询请求
type ReportRequest struct {
	StartTime       string   `json:"startTime"`                 // yyyy-MM-dd
	EndTime         string   `json:"endTime"`                   // yyyy-MM-dd
	TimeDimension   int      `json:"timeDimension,omitempty"`   // 0=天
	PageIndex       int      `json:"pageIndex,omitempty"`
	PageSize        int      `json:"pageSize,omitempty"`
	IndexScreenList []string `json:"indexScreenList,omitempty"`
}

// ReportMetrics 报表指标
type ReportMetrics struct {
	AdBilling      any `json:"adBilling"`     // 广告消耗（微，1元=1000000微）
	Impression     any `json:"impression"`    // 展示数
	Click          any `json:"click"`         // 点击数
	HonorPull      any `json:"honorPull"`     // 全网首唤数
	HonorPullCost  any `json:"honorPullCost"` // 全网首唤成本（微）
	Payment        any `json:"payment"`       // 付费数
	PaymentCost    any `json:"paymentCost"`   // 付费成本（微）
}

// ReportItem 报表数据项
type ReportItem struct {
	Metrics    ReportMetrics `json:"metrics"`
	StaticTime string        `json:"staticTime"` // yyyy-MM-dd
}

type reportPageResponse struct {
	TotalCount int          `json:"totalCount"`
	PageSize   int          `json:"pageSize"`
	PageIndex  int          `json:"pageIndex"`
	Data       []ReportItem `json:"data"`
}

type reportData struct {
	AdPageResponse      reportPageResponse `json:"adPageResponse"`
	AdAdvertiserStatcVo ReportMetrics      `json:"adAdvertiserStatcVo"`
}

// ReportResponse 广告主报表响应
type ReportResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    reportData `json:"data"`
}

// QueryAdvertiserReport 查询广告主报表（自动分页拉取全量数据）
func (c *HonorAPIClient) QueryAdvertiserReport(req ReportRequest, subAdvertiserID string) ([]ReportItem, error) {
	if req.PageSize <= 0 {
		req.PageSize = 100
	}
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.TimeDimension == 0 {
		req.TimeDimension = 0
	}

	var allItems []ReportItem
	for {
		resp, err := c.queryPage(req, subAdvertiserID)
		if err != nil {
			return nil, err
		}
		page := resp.Data.AdPageResponse
		allItems = append(allItems, page.Data...)
		if len(allItems) >= page.TotalCount || len(page.Data) < req.PageSize {
			break
		}
		req.PageIndex++
	}
	return allItems, nil
}

func (c *HonorAPIClient) queryPage(req ReportRequest, subAdvertiserID string) (*ReportResponse, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", reportURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("x-request-id", fmt.Sprintf("media-report-%d", time.Now().UnixNano()))
	if subAdvertiserID != "" {
		httpReq.Header.Set("sub-advertiser-id", subAdvertiserID)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var reportResp ReportResponse
	if err := json.Unmarshal(respBody, &reportResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(respBody))
	}
	if reportResp.Code != 0 {
		return nil, fmt.Errorf("API返回错误: code=%d, message=%s", reportResp.Code, reportResp.Message)
	}

	return &reportResp, nil
}

// ToFloat64 将数值转换为 float64（Honor API 金额单位：分）
func ToFloat64(v any) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case json.Number:
		f, _ := val.Float64()
		return f
	}
	return 0
}

// ToInt64 将数值转换为 int64
func ToInt64(v any) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	case json.Number:
		i, _ := val.Int64()
		return i
	}
	return 0
}

// ToString 将数值转为 string
func ToString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case json.Number:
		return val.String()
	}
	return fmt.Sprintf("%v", v)
}
