package internal

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// generateSign 按流程生成加密签名
// 1. 拼接 mediaId + originalUrl
// 2. 以 token 为密钥，对拼接字符串做 HMAC-SHA256
// 3. 对 HMAC 结果做 MD5，转大写，得到最终 sign
func generateSign(mediaId, token, originalUrl string) string {
	// Step 1: 拼接 mediaId + originalUrl
	concatStr := mediaId + originalUrl
	fmt.Printf("[签名] 拼接字符串: %s\n", concatStr)

	// Step 2: HMAC-SHA256，以 token 为密钥
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(concatStr))
	hmacResult := mac.Sum(nil)
	fmt.Printf("[签名] HMAC-SHA256 结果(hex): %s\n", hex.EncodeToString(hmacResult))

	// Step 3: MD5(hex(hmacResult)) → 大写
	// PHP hash_hmac 默认返回 hex 字符串，md5 是对 hex 字符串计算，而非原始字节
	hmacHex := hex.EncodeToString(hmacResult)
	md5Hash := md5.Sum([]byte(hmacHex))
	sign := strings.ToUpper(hex.EncodeToString(md5Hash[:]))
	fmt.Printf("[签名] MD5 大写结果: %s\n", sign)

	return sign
}

// TestGenerateSign 测试获取加密 sign
func TestGenerateSign(t *testing.T) {
	// TODO: 替换为实际的 mediaId、token、originalUrl
	mediaId := "100339"
	token := "OhqCUwzRIj1Pwp2D6d1oo0IbWM0OLHX0"
	originalUrl := "http://adn.atd.com/index.php?r=openapi/ocpx/advertiser-control&action=stop&advertiserId=100609&specialType=630"

	fmt.Printf("\n========== 生成加密 Sign ==========\n")
	fmt.Printf("mediaId:     %s\n", mediaId)
	fmt.Printf("token:       %s\n", strings.Repeat("*", len(token)))
	fmt.Printf("originalUrl: %s\n", originalUrl)
	fmt.Printf("------------------------------------\n")

	sign := generateSign(mediaId, token, originalUrl)

	fmt.Printf("====================================\n")
	fmt.Printf("最终 sign: %s\n", sign)
	fmt.Printf("====================================\n\n")
}
