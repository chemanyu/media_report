package internal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"
)

// ================================================================
// 京东联盟商品类目查询 —— 两步链路（绕开 bigfield 的 sceneId 权限）
//
//	sceneId=2 权限申请不下来，bigfield.query 走不通，改走：
//	  第一步：jd.union.open.sh.promotion.get
//	          account/taskId 固定，materialId=商品链接 → 拿到「推广链接」(clickURL)
//	  第二步：jd.union.open.goods.query
//	          把「推广链接」放到 keyword 里查询 → 取返回的 categoryInfo
//
//	参考 PHP: JdShController.php::fetchJdShPromotion
// ================================================================

// 固定参数（由业务方指定）
const (
	jdFixedAccount = "1870229189468553" // 固定 account
	jdFixedTaskId  = "11450"            // 固定 taskId
)

// JdUnionAPIClient 京东联盟开放平台API客户端
type JdUnionAPIClient struct {
	BaseURL     string
	AppKey      string
	AppSecret   string
	AccessToken string
	HTTPClient  *http.Client
}

// NewJdUnionAPIClient 创建京东联盟API客户端
func NewJdUnionAPIClient(appKey, appSecret, accessToken string) *JdUnionAPIClient {
	return &JdUnionAPIClient{
		BaseURL:     "https://api.jd.com/routerjson",
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ---------------------------------------------------------------
// 第一步：获取推广链接
// ---------------------------------------------------------------

// getCodeByTaskIdReq 对应 PHP 里的 360buy_param_json.getCodeByTaskIdReq
type getCodeByTaskIdReq struct {
	Account    string `json:"account"`
	MaterialId string `json:"materialId"`
	TaskId     string `json:"taskId"`
}

type shPromotionParam struct {
	GetCodeByTaskIdReq getCodeByTaskIdReq `json:"getCodeByTaskIdReq"`
}

// shPromotionData 推广接口返回的 data 部分（只关心几个链接字段）
type shPromotionData struct {
	ClickURL             string `json:"clickURL"`
	ClickUrlAlt          string `json:"clickUrl"` // 兼容大小写差异
	AppUrl               string `json:"appUrl"`
	ClickMonitorUrl      string `json:"clickMonitorUrl"`
	ImpressionMonitorUrl string `json:"impressionMonitorUrl"`
	ShortURL             string `json:"shortURL"`
}

// GetSHPromotion 调用 jd.union.open.sh.promotion.get，返回推广链接(clickURL)
//
//	materialId 传商品链接，例如 https://item.m.jd.com/product/100221373410.html
func (c *JdUnionAPIClient) GetSHPromotion(account, taskId, materialId string) (string, error) {
	fmt.Printf("\n========== 第一步：获取推广链接 ==========\n")

	bizBytes, err := json.Marshal(shPromotionParam{
		GetCodeByTaskIdReq: getCodeByTaskIdReq{
			Account:    account,
			MaterialId: materialId,
			TaskId:     taskId,
		},
	})
	if err != nil {
		return "", fmt.Errorf("序列化业务参数失败: %w", err)
	}
	paramJson := string(bizBytes)
	fmt.Printf("[请求参数] 360buy_param_json: %s\n", paramJson)

	// 系统参数：严格对齐 PHP fetchJdShPromotion（含空 access_token，无 format/sign_method）
	params := map[string]string{
		"access_token":      "",
		"method":            "jd.union.open.sh.promotion.get",
		"app_key":           c.AppKey,
		"timestamp":         jdTimestampMillis(), // 带毫秒+时区偏移，对齐 PHP
		"v":                 "1.0",
		"360buy_param_json": paramJson,
	}

	body, err := c.doRequest(params)
	if err != nil {
		return "", err
	}

	// 解析外层响应
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return "", fmt.Errorf("响应解析失败: %w", err)
	}
	if errResp, ok := envelope["error_response"]; ok {
		return "", fmt.Errorf("京东API返回错误: %s", string(errResp))
	}

	raw, ok := envelope["jd_union_open_sh_promotion_get_responce"]
	if !ok {
		return "", fmt.Errorf("响应缺少 jd_union_open_sh_promotion_get_responce 字段: %s", body)
	}
	var responce struct {
		Code      string `json:"code"`
		GetResult string `json:"getResult"`
	}
	if err := json.Unmarshal(raw, &responce); err != nil {
		return "", fmt.Errorf("responce 解析失败: %w", err)
	}
	if responce.Code != "0" {
		return "", fmt.Errorf("京东API响应异常: code=%s", responce.Code)
	}

	// 解析内层 getResult
	var getResult struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    shPromotionData `json:"data"`
	}
	if err := json.Unmarshal([]byte(responce.GetResult), &getResult); err != nil {
		return "", fmt.Errorf("getResult 解析失败: %w", err)
	}
	if getResult.Code != 200 {
		return "", fmt.Errorf("推广链接查询失败: code=%d, message=%s", getResult.Code, getResult.Message)
	}

	// 推广链接：优先 clickURL，兜底 clickUrl
	clickURL := getResult.Data.ClickURL
	if clickURL == "" {
		clickURL = getResult.Data.ClickUrlAlt
	}
	if clickURL == "" {
		return "", fmt.Errorf("未获取到推广链接(clickURL)，data=%+v", getResult.Data)
	}

	fmt.Printf("[结果] 推广链接 clickURL: %s\n", clickURL)
	return clickURL, nil
}

// ---------------------------------------------------------------
// 第二步：用推广链接作为 keyword 查询商品，取 categoryInfo
// ---------------------------------------------------------------

// CategoryInfo 类目信息（第二步唯一关心的返回）
type CategoryInfo struct {
	Cid1     int64  `json:"cid1"`
	Cid1Name string `json:"cid1Name"`
	Cid2     int64  `json:"cid2"`
	Cid2Name string `json:"cid2Name"`
	Cid3     int64  `json:"cid3"`
	Cid3Name string `json:"cid3Name"`
}

// goodsQueryParam 对应 360buy_param_json.goodsReqDTO
type goodsReqDTO struct {
	Keyword   string `json:"keyword"`
	PageIndex int    `json:"pageIndex,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
}

type goodsQueryParam struct {
	GoodsReqDTO goodsReqDTO `json:"goodsReqDTO"`
}

// QueryGoodsCategory 调用 jd.union.open.goods.query，keyword=推广链接，返回 categoryInfo
func (c *JdUnionAPIClient) QueryGoodsCategory(keyword string) (*CategoryInfo, error) {
	fmt.Printf("\n========== 第二步：按推广链接查询商品类目 ==========\n")

	// 文档参数虽多，业务方明确只需把推广链接放到 keyword
	bizBytes, err := json.Marshal(goodsQueryParam{
		GoodsReqDTO: goodsReqDTO{
			Keyword:   keyword,
			PageIndex: 1,
			PageSize:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("序列化业务参数失败: %w", err)
	}
	paramJson := string(bizBytes)
	fmt.Printf("[请求参数] 360buy_param_json: %s\n", paramJson)

	params := map[string]string{
		"method":            "jd.union.open.goods.query",
		"app_key":           c.AppKey,
		"timestamp":         time.Now().Format("2006-01-02 15:04:05"),
		"format":            "json",
		"v":                 "1.0",
		"sign_method":       "md5",
		"360buy_param_json": paramJson,
	}
	if c.AccessToken != "" {
		params["access_token"] = c.AccessToken
	}

	body, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	// 解析外层响应
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}
	if errResp, ok := envelope["error_response"]; ok {
		return nil, fmt.Errorf("京东API返回错误: %s", string(errResp))
	}

	raw, ok := envelope["jd_union_open_goods_query_responce"]
	if !ok {
		return nil, fmt.Errorf("响应缺少 jd_union_open_goods_query_responce 字段: %s", body)
	}
	var responce struct {
		Code        string `json:"code"`
		QueryResult string `json:"queryResult"`
	}
	if err := json.Unmarshal(raw, &responce); err != nil {
		return nil, fmt.Errorf("responce 解析失败: %w", err)
	}
	if responce.Code != "0" {
		return nil, fmt.Errorf("京东API响应异常: code=%s", responce.Code)
	}

	// 解析内层 queryResult，只取 data[].categoryInfo
	var queryResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			CategoryInfo CategoryInfo `json:"categoryInfo"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(responce.QueryResult), &queryResult); err != nil {
		return nil, fmt.Errorf("queryResult 解析失败: %w", err)
	}
	if queryResult.Code != 200 {
		return nil, fmt.Errorf("商品查询失败: code=%d, message=%s", queryResult.Code, queryResult.Message)
	}
	if len(queryResult.Data) == 0 {
		return nil, fmt.Errorf("商品查询无结果（keyword 可能无法匹配到商品）")
	}

	cat := queryResult.Data[0].CategoryInfo
	fmt.Printf("[结果] categoryInfo: cid1=%d(%s) cid2=%d(%s) cid3=%d(%s)\n",
		cat.Cid1, cat.Cid1Name, cat.Cid2, cat.Cid2Name, cat.Cid3, cat.Cid3Name)
	return &cat, nil
}

// ---------------------------------------------------------------
// 通用：签名 + HTTP 请求
// ---------------------------------------------------------------

// doRequest 计算签名并以 POST 表单方式发送，返回原始响应体
func (c *JdUnionAPIClient) doRequest(params map[string]string) (string, error) {
	params["sign"] = c.calculateSign(params)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	fmt.Printf("[HTTP请求] URL=%s method=%s timestamp=%s\n", c.BaseURL, params["method"], params["timestamp"])
	resp, err := c.HTTPClient.PostForm(c.BaseURL, form)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	fmt.Printf("[HTTP响应] 状态码=%d 长度=%d\n[HTTP响应] 内容: %s\n", resp.StatusCode, len(body), string(body))
	return string(body), nil
}

// calculateSign 计算京东签名
//
//	对应 PHP:
//	  ksort($params);
//	  $str = $secretkey;
//	  foreach ($params as $key=>$value) { $str .= $key . $value; }
//	  $str .= $secretkey;
//	  return strtoupper(md5($str));
func (c *JdUnionAPIClient) calculateSign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(c.AppSecret)
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(params[k])
	}
	sb.WriteString(c.AppSecret)

	hash := md5.Sum([]byte(sb.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// jdTimestampMillis 生成 PHP 风格时间戳：2026-03-04 11:32:26.123+0800
func jdTimestampMillis() string {
	cst := time.FixedZone("CST", 8*3600)
	return time.Now().In(cst).Format("2006-01-02 15:04:05.000-0700")
}

// ================================================================
// 测试
// ================================================================

// TestJdLinkToCategory 完整链路：商品链接 → 推广链接 → categoryInfo
func TestJdLinkToCategory(t *testing.T) {
	appKey := "513089760d51b7e98ca01c473f3b4741"
	appSecret := "8a677e1692cd49a790d3c76c72f53484"
	accessToken := ""

	// 待查询的商品链接（materialId）
	materialId := "https://item.m.jd.com/product/100221373410.html"

	client := NewJdUnionAPIClient(appKey, appSecret, accessToken)

	// 第一步：商品链接 → 推广链接
	promotionURL, err := client.GetSHPromotion(jdFixedAccount, jdFixedTaskId, materialId)
	if err != nil {
		t.Fatalf("获取推广链接失败: %v", err)
	}

	// 第二步：推广链接 → categoryInfo
	cat, err := client.QueryGoodsCategory(promotionURL)
	if err != nil {
		t.Fatalf("查询商品类目失败: %v", err)
	}

	fmt.Printf("\n========== 链路完成 ==========\n")
	fmt.Printf("商品链接: %s\n", materialId)
	fmt.Printf("推广链接: %s\n", promotionURL)
	fmt.Printf("一级类目: %d %s\n", cat.Cid1, cat.Cid1Name)
	fmt.Printf("二级类目: %d %s\n", cat.Cid2, cat.Cid2Name)
	fmt.Printf("三级类目: %d %s\n", cat.Cid3, cat.Cid3Name)
}

// TestJdCalculateSign 验证京东签名算法与 PHP 一致
func TestJdCalculateSign(t *testing.T) {
	client := NewJdUnionAPIClient("appkey", "secret", "")

	params := map[string]string{
		"method":      "jd.union.open.goods.query",
		"app_key":     "appkey",
		"timestamp":   "2018-01-01 12:00:00",
		"format":      "json",
		"v":           "1.0",
		"sign_method": "md5",
	}

	sign := client.calculateSign(params)

	var sb strings.Builder
	sb.WriteString("secret")
	keys := []string{"app_key", "format", "method", "sign_method", "timestamp", "v"}
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(params[k])
	}
	sb.WriteString("secret")
	h := md5.Sum([]byte(sb.String()))
	expected := strings.ToUpper(hex.EncodeToString(h[:]))

	if sign != expected {
		t.Fatalf("签名计算有误，期望 %s，实际 %s", expected, sign)
	}
	fmt.Printf("✓ 签名计算正确: %s\n", sign)
}
