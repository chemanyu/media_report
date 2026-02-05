package script

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

// refreshJuliangDLSAccessToken 刷新巨量DLS的 access token
func refreshJuliangDLSAccessToken(db *gorm.DB, juliangConfig config.JuliangConfig) {
	ctx := context.Background()
	logx.Infof("开始刷新巨量DLS access token - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 从数据库获取当前的 refresh_token
	mediaToken, err := model.GetByMedia(db, "juliang_dls")
	if err != nil {
		logx.Errorf("从数据库获取 juliang_dls token 失败: %v", err)
		return
	}

	// 创建 HTTP 客户端
	client := httpclient.NewClient(juliangConfig.BaseUrl, juliangConfig.Timeout)
	client.SetHeader("Content-Type", "application/json")

	// 构建刷新请求
	req := map[string]interface{}{
		"app_id":        juliangConfig.AppId,
		"secret":        juliangConfig.Secret,
		"refresh_token": mediaToken.RefreshToken,
	}

	// 调用刷新 token API
	var resp types.TokenRefreshResponse
	err = client.Post(ctx, "/open_api/oauth2/refresh_token/", req, &resp)
	if err != nil {
		logx.Errorf("调用刷新 juliang_dls token API 失败: %v", err)
		return
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("刷新 juliang_dls token 失败: code=%d, message=%s", resp.Code, resp.Message)
		return
	}

	// 更新数据库中的 token
	mediaToken.Token = resp.Data.AccessToken
	mediaToken.RefreshToken = resp.Data.RefreshToken
	err = db.Save(mediaToken).Error
	if err != nil {
		logx.Errorf("更新数据库 juliang_dls token 失败: %v", err)
		return
	}

	//logx.Infof("巨量DLS Token 刷新成功，新 AccessToken: %s, 有效期: %d 秒", resp.Data.AccessToken, resp.Data.ExpiresIn)
	//logx.Infof("新 RefreshToken: %s, 有效期: %d 秒", resp.Data.RefreshToken, resp.Data.RefreshTokenExpiresIn)
}

// refreshJuliangKHAccessToken 刷新巨量KH的 access token
func refreshJuliangKHAccessToken(db *gorm.DB, juliangConfig config.JuliangConfig) {
	ctx := context.Background()
	logx.Infof("开始刷新巨量KH access token - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 从数据库获取当前的 refresh_token
	mediaToken, err := model.GetByMedia(db, "juliang_kh")
	if err != nil {
		logx.Errorf("从数据库获取 juliang_kh token 失败: %v", err)
		return
	}

	// 创建 HTTP 客户端
	client := httpclient.NewClient(juliangConfig.BaseUrl, juliangConfig.Timeout)
	client.SetHeader("Content-Type", "application/json")

	// 构建刷新请求
	req := map[string]interface{}{
		"app_id":        juliangConfig.AppId,
		"secret":        juliangConfig.Secret,
		"refresh_token": mediaToken.RefreshToken,
	}

	// 调用刷新 token API
	var resp types.TokenRefreshResponse
	err = client.Post(ctx, "/open_api/oauth2/refresh_token/", req, &resp)
	if err != nil {
		logx.Errorf("调用刷新 juliang_kh token API 失败: %v", err)
		return
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("刷新 juliang_kh token 失败: code=%d, message=%s", resp.Code, resp.Message)
		return
	}

	// 更新数据库中的 token
	mediaToken.Token = resp.Data.AccessToken
	mediaToken.RefreshToken = resp.Data.RefreshToken
	err = db.Save(mediaToken).Error
	if err != nil {
		logx.Errorf("更新数据库 juliang_kh token 失败: %v", err)
		return
	}

	//logx.Infof("巨量KH Token 刷新成功，新 AccessToken: %s, 有效期: %d 秒", resp.Data.AccessToken, resp.Data.ExpiresIn)
	//logx.Infof("新 RefreshToken: %s, 有效期: %d 秒", resp.Data.RefreshToken, resp.Data.RefreshTokenExpiresIn)
}

// FetchHuichuanElmReports 获取回传饿了么所有账户的报表数据并发送到ADX
func FetchHuichuanElmReports(db *gorm.DB, juliangConfig config.JuliangConfig, adxConfig config.ADXConfig) {
	logx.Infof("开始获取回传饿了么报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 获取昨天的日期（因为数据录入时间要求：每日7点前完成数据录入）
	yesterday := time.Now().AddDate(0, 0, -1)
	dt := yesterday.Format("20060102")
	startTime := yesterday.Format("2006-01-02") + " 00:00:00"
	endTime := yesterday.Format("2006-01-02") + " 23:59:59"

	logx.Infof("查询日期: %s, 时间范围: %s ~ %s", dt, startTime, endTime)

	// 从数据库获取 access_token（在循环外查询一次，避免频繁查询数据库）
	mediaToken, err := model.GetByMedia(db, "juliang_dls")
	if err != nil {
		logx.Errorf("从数据库获取 juliang_dls token 失败: %v", err)
		return
	}

	// 从数据库获取所有客户及其媒体账户
	performances, err := model.GetAllElmHcPerformanceReports(db)
	if err != nil {
		logx.Errorf("获取客户列表失败: %v", err)
		return
	}

	if len(performances) == 0 {
		logx.Info("暂无客户配置")
		return
	}

	// 收集所有需要发送的数据
	var allReportData []types.ADXReportData

	// 遍历所有客户
	for _, performance := range performances {
		logx.Infof("处理客户: %s (%s)", performance.CustomerName, performance.CustomerShort)

		// 获取该客户的所有媒体账户
		mediaReports, err := model.GetElmHcMediaReportsByPerformanceId(db, int(performance.ID))
		if err != nil {
			logx.Errorf("获取客户 %s 的媒体账户失败: %v", performance.CustomerShort, err)
			continue
		}

		if len(mediaReports) == 0 {
			logx.Infof("客户 %s 暂无媒体账户配置", performance.CustomerShort)
			continue
		}

		// 遍历该客户的所有媒体账户
		for _, media := range mediaReports {
			logx.Infof("  正在获取账户 %s (汇川ID: %d) 的报表数据...", media.MediaAdvName, media.HuichuanAdvId)

			// 调用巨量引擎API获取报表数据
			advertiser_id, _ := strconv.Atoi(media.MediaAdvId)
			resp, err := getJuliangReportData(juliangConfig, mediaToken.Token, advertiser_id, startTime, endTime, "stat_time_day")
			if err != nil {
				logx.Errorf("获取账户 %s 的报表数据失败: %v", media.MediaAdvName, err)
				continue
			}

			// 处理报表数据并转换为ADX格式
			if len(resp.Data.Rows) > 0 {
				for _, row := range resp.Data.Rows {
					// 从 map 中提取数据（巨量接口返回的是字符串类型）
					var cost float64
					if v, ok := row.Metrics["stat_cost"]; ok {
						if val, ok := v.(string); ok {
							cost, _ = strconv.ParseFloat(val, 64)
						} else if val, ok := v.(float64); ok {
							cost = val
						}
					}

					var showNum, clickNum, convertNum int64
					if v, ok := row.Metrics["show_cnt"]; ok {
						if val, ok := v.(string); ok {
							showNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							showNum = int64(val)
						}
					}
					if v, ok := row.Metrics["click_cnt"]; ok {
						if val, ok := v.(string); ok {
							clickNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							clickNum = int64(val)
						}
					}
					if v, ok := row.Metrics["convert_cnt"]; ok {
						if val, ok := v.(string); ok {
							convertNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							convertNum = int64(val)
						}
					}

					reportData := types.ADXReportData{
						CustomerName:      performance.CustomerName,
						CustomerShort:     performance.CustomerShort,
						AgentName:         performance.AgentName,
						AgentShort:        performance.AgentShort,
						MediaPlatformName: performance.MediaPlatformName,
						MediaAdvId:        media.MediaAdvId,
						MediaAdvName:      media.MediaAdvName,
						HuichuanAdvId:     media.HuichuanAdvId,
						Cost:              cost,
						ShowNum:           showNum,
						ClickNum:          clickNum,
						ConvertNum:        convertNum,
						DeepConvertNum:    media.PayNum,
						ConvertType:       "调起",
						DeepConvertType:   "付费",
						RedirectNum:       media.RedirectNum,
						PayNum:            media.PayNum,
						Dt:                dt,
					}
					allReportData = append(allReportData, reportData)
				}
				logx.Infof("  账户 %s 获取到 %d 条记录", media.MediaAdvName, len(resp.Data.Rows))
			} else {
				logx.Infof("  账户 %s 暂无数据", media.MediaAdvName)
			}
		}
	}

	// 发送数据到ADX接口
	if len(allReportData) > 0 {
		logx.Infof("准备发送 %d 条数据到ADX接口", len(allReportData))
		err := sendDataToADX(adxConfig, allReportData)
		if err != nil {
			logx.Errorf("发送数据到ADX失败: %v", err)
		} else {
			logx.Infof("数据发送成功")
		}
	} else {
		logx.Info("暂无数据需要发送")
	}

	logx.Infof("回传饿了么报表数据获取完成")
}

// getJuliangReportData 获取巨量引擎报表数据
func getJuliangReportData(juliangConfig config.JuliangConfig, accessToken string, advertiserId int, startTime, endTime, stat_date string) (*types.JuliangCustomReportResp, error) {
	ctx := context.Background()
	logx.Infof("开始获取巨量引擎报表数据 - advertiser_id: %d, 时间范围: %s ~ %s", advertiserId, startTime, endTime)

	// 创建 HTTP 客户端
	client := httpclient.NewClient(juliangConfig.BaseUrl, juliangConfig.Timeout)
	client.SetHeader("Access-Token", accessToken)

	// 构建查询参数（需要序列化为JSON字符串）
	dimensions := []string{
		stat_date,
		"external_action",
		"deep_external_action",
	}
	metrics := []string{
		"stat_cost",
		"show_cnt",
		"click_cnt",
		"convert_cnt",
		"in_app_pay",
	}
	filters := []interface{}{}
	orderBy := []types.JuliangOrderBy{
		{
			Field: stat_date,
			Type:  "DESC",
		},
	}

	// 序列化为JSON字符串
	dimensionsJSON, _ := json.Marshal(dimensions)
	metricsJSON, _ := json.Marshal(metrics)
	filtersJSON, _ := json.Marshal(filters)
	orderByJSON, _ := json.Marshal(orderBy)

	// 构建URL查询参数
	params := map[string]string{
		"advertiser_id": fmt.Sprintf("%d", advertiserId),
		"dimensions":    string(dimensionsJSON),
		"metrics":       string(metricsJSON),
		"filters":       string(filtersJSON),
		"start_time":    startTime,
		"end_time":      endTime,
		"order_by":      string(orderByJSON),
	}

	// 调用报表 API (GET请求，参数通过query string传递)
	var resp types.JuliangCustomReportResp
	err := client.Get(ctx, "/open_api/v3.0/report/custom/get/", params, &resp)
	if err != nil {
		logx.Errorf("调用巨量引擎报表 API 失败: %v", err)
		return nil, err
	}

	// 检查响应
	if resp.Code != 0 {
		logx.Errorf("获取巨量引擎报表失败: code=%d, message=%s", resp.Code, resp.Message)
		return nil, fmt.Errorf("获取报表失败: %s", resp.Message)
	}

	logx.Infof("成功获取巨量引擎报表数据：%v", resp.Data.Rows)
	return &resp, nil
}

// FetchHuichuanElmReportsByHour 获取回传饿了么所有账户的小时级报表数据并发送到ADX
func FetchHuichuanElmReportsByHour(db *gorm.DB, juliangConfig config.JuliangConfig, adxConfig config.ADXConfig) {
	logx.Infof("开始获取回传饿了么小时级报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 获取当前时间，计算上一个小时的时间范围
	now := time.Now()
	// 上一个小时的开始时间：当前小时-1，分钟和秒为00
	lastHour := now.Add(-1 * time.Hour)
	startTime := lastHour.Format("2006-01-02 15") + ":00:00"
	// 上一个小时的结束时间：当前小时，分钟为00，秒为00
	endTime := lastHour.Format("2006-01-02 15") + ":59:59"
	// 日期和小时
	dt := lastHour.Format("20060102")
	hh := lastHour.Format("15") // 24小时制的小时，如：01, 02, 15, 23

	logx.Infof("查询日期: %s, 小时: %s, 时间范围: %s ~ %s", dt, hh, startTime, endTime)

	// 从数据库获取 access_token（在循环外查询一次，避免频繁查询数据库）
	mediaToken, err := model.GetByMedia(db, "juliang_dls")
	if err != nil {
		logx.Errorf("从数据库获取 juliang_dls token 失败: %v", err)
		return
	}

	// 从数据库获取所有客户及其媒体账户
	performances, err := model.GetAllElmHcPerformanceReports(db)
	if err != nil {
		logx.Errorf("获取客户列表失败: %v", err)
		return
	}

	if len(performances) == 0 {
		logx.Info("暂无客户配置")
		return
	}

	// 收集所有需要发送的数据
	var allReportData []types.ADXReportData

	// 遍历所有客户
	for _, performance := range performances {
		logx.Infof("处理客户: %s (%s)", performance.CustomerName, performance.CustomerShort)

		// 获取该客户的所有媒体账户
		mediaReports, err := model.GetElmHcMediaReportsByPerformanceId(db, int(performance.ID))
		if err != nil {
			logx.Errorf("获取客户 %s 的媒体账户失败: %v", performance.CustomerShort, err)
			continue
		}

		if len(mediaReports) == 0 {
			logx.Infof("客户 %s 暂无媒体账户配置", performance.CustomerShort)
			continue
		}

		// 遍历该客户的所有媒体账户
		for _, media := range mediaReports {
			logx.Infof("  正在获取账户 %s (汇川ID: %d) 的小时级报表数据...", media.MediaAdvName, media.HuichuanAdvId)

			// 调用巨量引擎API获取报表数据
			advertiser_id, _ := strconv.Atoi(media.MediaAdvId)
			resp, err := getJuliangReportData(juliangConfig, mediaToken.Token, advertiser_id, startTime, endTime, "stat_time_hour")
			if err != nil {
				logx.Errorf("获取账户 %s 的小时级报表数据失败: %v", media.MediaAdvName, err)
				continue
			}

			// 处理报表数据并转换为ADX格式
			if len(resp.Data.Rows) > 0 {
				for _, row := range resp.Data.Rows {
					// 从 map 中提取数据（巨量接口返回的是字符串类型）
					var cost float64
					if v, ok := row.Metrics["stat_cost"]; ok {
						if val, ok := v.(string); ok {
							cost, _ = strconv.ParseFloat(val, 64)
						} else if val, ok := v.(float64); ok {
							cost = val
						}
					}

					var showNum, clickNum, convertNum int64
					if v, ok := row.Metrics["show_cnt"]; ok {
						if val, ok := v.(string); ok {
							showNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							showNum = int64(val)
						}
					}
					if v, ok := row.Metrics["click_cnt"]; ok {
						if val, ok := v.(string); ok {
							clickNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							clickNum = int64(val)
						}
					}
					if v, ok := row.Metrics["convert_cnt"]; ok {
						if val, ok := v.(string); ok {
							convertNum, _ = strconv.ParseInt(val, 10, 64)
						} else if val, ok := v.(float64); ok {
							convertNum = int64(val)
						}
					}

					reportData := types.ADXReportData{
						CustomerName:      performance.CustomerName,
						CustomerShort:     performance.CustomerShort,
						AgentName:         performance.AgentName,
						AgentShort:        performance.AgentShort,
						MediaPlatformName: performance.MediaPlatformName,
						MediaAdvId:        media.MediaAdvId,
						MediaAdvName:      media.MediaAdvName,
						HuichuanAdvId:     media.HuichuanAdvId,
						Cost:              cost,
						ShowNum:           showNum,
						ClickNum:          clickNum,
						ConvertNum:        convertNum,
						DeepConvertNum:    media.PayNum,
						ConvertType:       "调起",
						DeepConvertType:   "付费",
						RedirectNum:       media.RedirectNum,
						PayNum:            media.PayNum,
						Dt:                dt,
						Hh:                hh,
					}
					allReportData = append(allReportData, reportData)
				}
				logx.Infof("  账户 %s 获取到 %d 条记录", media.MediaAdvName, len(resp.Data.Rows))
			} else {
				logx.Infof("  账户 %s 暂无数据", media.MediaAdvName)
			}
		}
	}

	// 发送数据到ADX小时接口
	if len(allReportData) > 0 {
		logx.Infof("准备发送 %d 条小时数据到ADX接口", len(allReportData))
		err := sendHourDataToADX(adxConfig, allReportData)
		if err != nil {
			logx.Errorf("发送小时数据到ADX失败: %v", err)
		} else {
			logx.Infof("小时数据发送成功")
		}
	} else {
		logx.Info("暂无数据需要发送")
	}

	logx.Infof("回传饿了么小时级报表数据获取完成")
}

// generateSignature 生成 HMAC-SHA256 签名
func generateSignature(secret string, path string, timestamp string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(path + timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

// sendDataToADX 发送数据到ADX接口（支持批量，最大100条）
func sendDataToADX(adxConfig config.ADXConfig, data []types.ADXReportData) error {
	ctx := context.Background()

	// 检查数据量，如果超过100条，需要分批发送
	batchSize := 100
	totalBatches := (len(data) + batchSize - 1) / batchSize

	logx.Infof("开始发送数据到ADX，总数据量: %d, 分 %d 批发送", len(data), totalBatches)

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		batchNum := i/batchSize + 1

		logx.Infof("发送第 %d/%d 批数据，本批数量: %d", batchNum, totalBatches, len(batch))

		// 生成时间戳（过期时间，例如：1747211953374，13位毫秒级时间戳）
		timestamp := strconv.FormatInt(time.Now().Add(5*time.Minute).UnixMilli(), 10)

		// 调用 ADX API
		url := "/adx/agent/customer/media/data/day/input"

		// 生成签名
		signature := generateSignature(adxConfig.Secret, "/assistant-external"+url, timestamp)

		// 创建 HTTP 客户端
		client := httpclient.NewClient(adxConfig.BaseURL, adxConfig.Timeout)
		client.SetHeader("Content-Type", "application/json")
		client.SetHeader("X-API-KEY", adxConfig.APIKey)
		client.SetHeader("X-Timestamp", timestamp)
		client.SetHeader("X-Signature", signature)

		var resp types.ADXResponse
		err := client.Post(ctx, url, batch, &resp)
		batchJSON, _ := json.Marshal(batch)
		headers := map[string]string{
			"Content-Type": "application/json",
			"X-API-KEY":    adxConfig.APIKey,
			"X-Timestamp":  timestamp,
			"X-Signature":  signature,
		}
		if err != nil {
			headersJSON, _ := json.Marshal(headers)
			logx.Errorf("第 %d 批数据发送失败: %v, URL: %s, Headers: %s, 数据: %s", batchNum, err, adxConfig.BaseURL+url, string(headersJSON), string(batchJSON))
			return fmt.Errorf("第 %d 批数据发送失败: %v", batchNum, err)
		}
		//logx.Infof("day input: %v", string(batchJSON))

		// 检查响应
		if !resp.Data {
			logx.Errorf("第 %d 批数据ADX接口返回失败", batchNum)
			return fmt.Errorf("第 %d 批数据ADX接口返回失败", batchNum)
		}

		logx.Infof("第 %d/%d 批数据发送成功", batchNum, totalBatches)
	}

	logx.Infof("所有数据发送完成，共 %d 批", totalBatches)
	return nil
}

// sendHourDataToADX 发送小时数据到ADX接口（支持批量，最大100条）
func sendHourDataToADX(adxConfig config.ADXConfig, data []types.ADXReportData) error {
	ctx := context.Background()

	// 检查数据量，如果超过100条，需要分批发送
	batchSize := 100
	totalBatches := (len(data) + batchSize - 1) / batchSize

	logx.Infof("开始发送小时数据到ADX，总数据量: %d, 分 %d 批发送", len(data), totalBatches)

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		batchNum := i/batchSize + 1

		logx.Infof("发送第 %d/%d 批小时数据，本批数量: %d", batchNum, totalBatches, len(batch))

		// 生成时间戳（过期时间，例如：1747211953374，13位毫秒级时间戳）
		timestamp := strconv.FormatInt(time.Now().Add(5*time.Minute).UnixMilli(), 10)

		// 调用 ADX 小时数据接口
		url := "/adx/agent/customer/media/data/hour/input"

		// 生成签名
		signature := generateSignature(adxConfig.Secret, "/assistant-external"+url, timestamp)

		// 创建 HTTP 客户端
		client := httpclient.NewClient(adxConfig.BaseURL, adxConfig.Timeout)
		client.SetHeader("Content-Type", "application/json")
		client.SetHeader("X-API-KEY", adxConfig.APIKey)
		client.SetHeader("X-Timestamp", timestamp)
		client.SetHeader("X-Signature", signature)

		var resp types.ADXResponse
		err := client.Post(ctx, url, batch, &resp)
		batchJSON, _ := json.Marshal(batch)
		headers := map[string]string{
			"Content-Type": "application/json",
			"X-API-KEY":    adxConfig.APIKey,
			"X-Timestamp":  timestamp,
			"X-Signature":  signature,
		}
		headersJSON, _ := json.Marshal(headers)
		if err != nil {
			logx.Errorf("第 %d 批小时数据发送失败: %v, URL: %s, Headers: %s, 数据: %s", batchNum, err, adxConfig.BaseURL+url, string(headersJSON), string(batchJSON))
			return fmt.Errorf("第 %d 批小时数据发送失败: %v", batchNum, err)
		}
		//logx.Infof("hour input: %v", string(batchJSON))

		// 检查响应
		if !resp.Data {
			logx.Errorf("第 %d 批小时数据ADX接口返回失败", batchNum)
			return fmt.Errorf("第 %d 批小时数据ADX接口返回失败", batchNum)
		}

		logx.Infof("第 %d/%d 批小时数据发送成功", batchNum, totalBatches)
	}

	logx.Infof("所有小时数据发送完成，共 %d 批", totalBatches)
	return nil
}
