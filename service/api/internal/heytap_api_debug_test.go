package internal

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// 测试Token生成 - 详细调试版本
func TestHeytapAPICallWithDebug(t *testing.T) {
	// 广告主信息
	ownerId := int64(1000350051)
	apiId := "940ccf6f3b134d388681f16cf6e00234"
	apiKey := "3eb9c90101ec4a4a97dc228d9842c827"

	// 生成时间戳（使用UTC+08:00中国北京时间）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := time.Now().In(loc)
	timestamp := beijingTime.Unix()

	fmt.Printf("\n========== Token生成调试信息 ==========\n")
	fmt.Printf("当前北京时间: %s\n", beijingTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("广告主ID (owner_id): %d\n", ownerId)
	fmt.Printf("API-ID (api_id): %s\n", apiId)
	fmt.Printf("API-KEY (api_key): %s\n", apiKey)
	fmt.Printf("时间戳 (timestamp): %d\n", timestamp)

	// 1. 生成签名：sign=sha1(api_id+api_key+time_stamp)
	signStr := fmt.Sprintf("%s%s%d", apiId, apiKey, timestamp)
	fmt.Printf("\n签名前的字符串: %s\n", signStr)

	h := sha1.New()
	h.Write([]byte(signStr))
	sign := hex.EncodeToString(h.Sum(nil))
	fmt.Printf("SHA1签名结果: %s\n", sign)

	// 2. 拼接token字符串：owner_id+","+api_id+","+time_stamp+","+sign
	tokenStr := fmt.Sprintf("%d,%s,%d,%s", ownerId, apiId, timestamp, sign)
	fmt.Printf("\nToken编码前的字符串: %s\n", tokenStr)

	// 3. Base64编码
	token := base64.StdEncoding.EncodeToString([]byte(tokenStr))
	fmt.Printf("Base64编码后的Token: %s\n", token)
	fmt.Printf("========================================\n\n")

	// 构建请求体（使用JSON格式）
	requestBody := map[string]interface{}{
		"page":           1,
		"pageCount":      10,
		"beginTime":      "20260210", // 使用更近的日期
		"endTime":        "2026021-",
		"timeLevel":      "DAY",
		"orderByColumns": "dt",
		"ascDesc":        "ASC",
		"paraMap": map[string]interface{}{
			"owner_type":    0,
			"groupByColumn": "",
		},
	}

	// 序列化请求体
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("序列化请求体失败: %v", err)
	}
	fmt.Printf("请求体JSON: %s\n\n", string(jsonData))

	// 创建HTTP请求
	url := "https://sapi.ads.heytapmobi.com/v3/data/common/agency/query/queryAgencyEffect"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Printf("========== 请求信息 ==========\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Method: POST\n")
	fmt.Printf("Content-Type: application/json\n")
	fmt.Printf("Authorization: %s\n", token)
	fmt.Printf("==============================\n\n")

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}

	fmt.Printf("========== 响应信息 ==========\n")
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))
	fmt.Printf("==============================\n")

	// 解析响应以获取更详细的错误信息
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if code, ok := result["code"].(float64); ok && code != 0 {
			fmt.Printf("\n❌ 错误代码: %.0f\n", code)
			if msg, ok := result["msg"].(string); ok {
				fmt.Printf("❌ 错误信息: %s\n", msg)
			}
		} else {
			fmt.Printf("\n✅ 请求成功！\n")
		}
	}
}
