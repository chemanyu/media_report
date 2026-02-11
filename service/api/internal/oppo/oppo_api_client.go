package oppo

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OppoAPIConfig OPPO API配置
type OppoAPIConfig struct {
	OwnerId int    `json:"ownerId"`
	ApiId   string `json:"api_id"`
	ApiKey  string `json:"api_key"`
}

// OppoAPIClient OPPO广告API客户端
type OppoAPIClient struct {
	config OppoAPIConfig
	client *http.Client
}

// AdDataRequest 查询广告数据请求参数
type AdDataRequest struct {
	BeginTime string                 `json:"beginTime"` // 格式: 20260211
	EndTime   string                 `json:"endTime"`   // 格式: 20260211
	TimeLevel string                 `json:"timeLevel"` // DAY, HOUR
	OwnerId   int64                  `json:"ownerId"`   // 广告主ID
	ParaMap   map[string]interface{} `json:"paraMap"`
}

// AdDataResponse 广告数据响应
type AdDataResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// AdDataDetail 广告数据详情（提取需要的字段）
type AdDataDetail struct {
	OwnerId         string  `json:"owner_id"`           // 广告主ID
	OwnerName       string  `json:"owner_name"`         // 广告主名称
	Cost            float64 `json:"cost"`               // 消耗
	ConvertDp       int64   `json:"convert_dp"`         // deeplink拉活数
	DpAppOrderNums  int64   `json:"dp_app_order_nums"`  // 订单数
	Ftime           int64   `json:"ftime"`              // 日期（格式: 20260211）
	Click           int64   `json:"click"`              // 点击数
	Expose          int64   `json:"expose"`             // 曝光数
	ConvertDpPrice  float64 `json:"convert_dp_price"`   // 拉活成本
	DpAppOrderPrice float64 `json:"dp_app_order_price"` // 订单成本
}

// NewOppoAPIClient 创建OPPO API客户端
func NewOppoAPIClient(ownerId int, apiId, apiKey string) *OppoAPIClient {
	return &OppoAPIClient{
		config: OppoAPIConfig{
			OwnerId: ownerId,
			ApiId:   apiId,
			ApiKey:  apiKey,
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateToken 生成API访问Token
// Token生成规则：
// 1. sign = sha1(api_id + api_key + timestamp)
// 2. token_str = owner_id + "," + api_id + "," + timestamp + "," + sign
// 3. token = base64(token_str)
func (c *OppoAPIClient) generateToken() string {
	// 使用北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := time.Now().In(loc)
	timestamp := beijingTime.Unix()

	// 1. 生成签名
	signStr := fmt.Sprintf("%s%s%d", c.config.ApiId, c.config.ApiKey, timestamp)
	h := sha1.New()
	h.Write([]byte(signStr))
	sign := hex.EncodeToString(h.Sum(nil))

	// 2. 拼接token字符串
	tokenStr := fmt.Sprintf("%d,%s,%d,%s", c.config.OwnerId, c.config.ApiId, timestamp, sign)

	// 3. Base64编码
	token := base64.StdEncoding.EncodeToString([]byte(tokenStr))

	return token
}

// QueryAdData 查询广告数据
func (c *OppoAPIClient) QueryAdData(req AdDataRequest) (*AdDataDetail, error) {
	// 生成Token
	token := c.generateToken()

	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	url := "https://sapi.ads.heytapmobi.com/v3/data/common/summary/queryAdData"
	httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	// 发送请求
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// fmt.Printf("========== 响应信息 ==========\n")
	// fmt.Printf("状态码: %d\n", resp.StatusCode)
	// fmt.Printf("响应内容: %s\n", string(body))
	// fmt.Printf("==============================\n")

	// 解析响应
	var apiResp AdDataResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 响应内容: %s", err, string(body))
	}

	// 检查响应码
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API返回错误: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	}

	// 解析数据
	var adData AdDataDetail
	if err := json.Unmarshal(apiResp.Data, &adData); err != nil {
		return nil, fmt.Errorf("解析广告数据失败: %w", err)
	}

	return &adData, nil
}

// QueryAdDataBatch 批量查询广告数据（多个广告主）
func (c *OppoAPIClient) QueryAdDataBatch(ownerIds []int64, beginTime, endTime string) ([]*AdDataDetail, error) {
	var results []*AdDataDetail

	for _, ownerId := range ownerIds {
		req := AdDataRequest{
			BeginTime: beginTime,
			EndTime:   endTime,
			TimeLevel: "DAY",
			OwnerId:   ownerId,
			ParaMap: map[string]interface{}{
				"filter_zero": 0,
			},
		}

		data, err := c.QueryAdData(req)
		if err != nil {
			// 记录错误但继续处理其他账户
			fmt.Printf("查询广告主 %d 数据失败: %v\n", ownerId, err)
			continue
		}

		results = append(results, data)
	}

	return results, nil
}
