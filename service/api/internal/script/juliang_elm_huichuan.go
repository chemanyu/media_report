package script

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"media_report/common/httpclient"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/model"
	"media_report/service/api/internal/types"
)

// ErrJuliangRateLimit 表示巨量接口返回请求频率超限（code=40110），调用方可据此重试。
var ErrJuliangRateLimit = errors.New("巨量请求频率超限")

// isRateLimitErr 判断错误是否为巨量限频错误。
func isRateLimitErr(err error) bool {
	return errors.Is(err, ErrJuliangRateLimit)
}

// elmHcRetryBackoffs 限频重试的退避间隔（递增，给上游让出请求窗口）；超过长度后取最后一档。
var elmHcRetryBackoffs = []time.Duration{
	30 * time.Second,
	60 * time.Second,
	120 * time.Second,
	300 * time.Second,
}

// elmHcTask 一个待抓取的账户任务（客户 + 媒体账户）。
type elmHcTask struct {
	perf  model.ElmHcPerformanceReport
	media model.ElmHcMediaReport
}

// collectElmHcTasks 把所有客户及其媒体账户平铺成任务列表。
func collectElmHcTasks(db *gorm.DB) ([]elmHcTask, error) {
	performances, err := model.GetAllElmHcPerformanceReports(db)
	if err != nil {
		return nil, fmt.Errorf("获取客户列表失败: %w", err)
	}
	if len(performances) == 0 {
		return nil, nil
	}

	var tasks []elmHcTask
	for _, performance := range performances {
		logx.Infof("处理客户: %s (%s)", performance.CustomerName, performance.CustomerShort)

		mediaReports, err := model.GetElmHcMediaReportsByPerformanceId(db, int(performance.ID))
		if err != nil {
			logx.Errorf("获取客户 %s 的媒体账户失败: %v", performance.CustomerShort, err)
			continue
		}
		if len(mediaReports) == 0 {
			logx.Infof("客户 %s 暂无媒体账户配置", performance.CustomerShort)
			continue
		}
		for _, media := range mediaReports {
			tasks = append(tasks, elmHcTask{perf: performance, media: media})
		}
	}
	return tasks, nil
}

// fetchElmHcReportsWithRetry 对任务列表逐个抓取报表，限频失败的账户按退避表多轮重试，
// 直到全部成功、或达到 deadline 时间预算（避免与下一次 cron 重叠）。
// statDate 为巨量维度（stat_time_day / stat_time_hour），buildRow 把单账户响应转换为 ADX 数据。
func fetchElmHcReportsWithRetry(
	juliangConfig config.JuliangConfig,
	token string,
	tasks []elmHcTask,
	startTime, endTime, statDate string,
	deadline time.Time,
	buildRow func(t elmHcTask, resp *types.JuliangCustomReportResp) []types.ADXReportData,
) []types.ADXReportData {
	var allReportData []types.ADXReportData
	remaining := tasks

	for round := 0; len(remaining) > 0; round++ {
		var stillFailed []elmHcTask

		for _, t := range remaining {
			logx.Infof("  正在获取账户 %s (汇川ID: %d) 的报表数据...", t.media.MediaAdvName, t.media.HuichuanAdvId)

			advertiserId, _ := strconv.Atoi(t.media.MediaAdvId)
			resp, err := getJuliangReportData(juliangConfig, token, advertiserId, startTime, endTime, statDate)
			if err != nil {
				if isRateLimitErr(err) {
					logx.Infof("  账户 %s 被限频，加入重试队列: %v", t.media.MediaAdvName, err)
					stillFailed = append(stillFailed, t)
				} else {
					logx.Errorf("获取账户 %s 的报表数据失败(非限频，跳过): %v", t.media.MediaAdvName, err)
				}
				continue
			}

			if len(resp.Data.Rows) > 0 {
				allReportData = append(allReportData, buildRow(t, resp)...)
				logx.Infof("  账户 %s 获取到 %d 条记录", t.media.MediaAdvName, len(resp.Data.Rows))
			} else {
				logx.Infof("  账户 %s 暂无数据", t.media.MediaAdvName)
			}
		}

		if len(stillFailed) == 0 {
			break
		}

		wait := elmHcRetryBackoffs[len(elmHcRetryBackoffs)-1]
		if round < len(elmHcRetryBackoffs) {
			wait = elmHcRetryBackoffs[round]
		}
		if time.Now().Add(wait).After(deadline) {
			logx.Errorf("仍有 %d 个账户因限频未成功，已达时间预算上限，放弃本次重试（避免与下一次任务重叠）", len(stillFailed))
			break
		}

		logx.Infof("本轮第 %d 次重试结束，仍有 %d 个账户被限频，%v 后重试…", round+1, len(stillFailed), wait)
		time.Sleep(wait)
		remaining = stillFailed
	}

	return allReportData
}

// buildElmHcReportData 把单个账户的巨量响应行转换为 ADX 回传数据。
// isHourly=true 为小时级（RedirectNum/PayNum 固定 0，带 hh）；
// isHourly=false 为日级（按 update_time 是否为今天决定 RedirectNum/PayNum）。
func buildElmHcReportData(
	perf model.ElmHcPerformanceReport,
	media model.ElmHcMediaReport,
	resp *types.JuliangCustomReportResp,
	dt, hh string,
	isHourly bool,
) []types.ADXReportData {
	// 日级：检查 update_time 是否为今天，不是今天则 RedirectNum 和 PayNum 设置为 0
	redirectNum := media.RedirectNum
	payNum := media.PayNum
	if isHourly {
		redirectNum = 0
		payNum = 0
	} else {
		today := time.Now().Format("2006-01-02")
		updateDate := media.UpdateTime.Format("2006-01-02")
		if updateDate != today {
			redirectNum = 0
			payNum = 0
			logx.Infof("  账户 %s 的 update_time (%s) 不是今天，RedirectNum 和 PayNum 设置为 0", media.MediaAdvName, updateDate)
		}
	}

	var result []types.ADXReportData
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

		result = append(result, types.ADXReportData{
			CustomerName:      perf.CustomerName,
			CustomerShort:     perf.CustomerShort,
			AgentName:         perf.AgentName,
			AgentShort:        perf.AgentShort,
			MediaPlatformName: perf.MediaPlatformName,
			HuichuanAdvId:     media.HuichuanAdvId,
			Cost:              cost,
			ShowNum:           showNum,
			ClickNum:          clickNum,
			ConvertNum:        convertNum,
			DeepConvertNum:    media.PayNum,
			ConvertType:       "调起",
			DeepConvertType:   "付费",
			RedirectNum:       redirectNum,
			PayNum:            payNum,
			Dt:                dt,
			Hh:                hh,
		})
	}
	return result
}

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
// reportDate: 可选参数，格式为 20060102，如果为空则使用昨天的日期
func FetchHuichuanElmReports(db *gorm.DB, juliangConfig config.JuliangConfig, adxConfig config.ADXConfig, reportDate string) {
	logx.Infof("开始获取回传饿了么报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 获取目标日期
	var targetDate time.Time
	if reportDate != "" {
		// 使用指定日期
		var err error
		targetDate, err = time.Parse("20060102", reportDate)
		if err != nil {
			logx.Errorf("日期格式错误: %v，使用昨天的日期", err)
			targetDate = time.Now().AddDate(0, 0, -1)
		}
	} else {
		// 默认使用昨天的日期（因为数据录入时间要求：每日7点前完成数据录入）
		targetDate = time.Now().AddDate(0, 0, -1)
	}

	dt := targetDate.Format("20060102")
	startTime := targetDate.Format("2006-01-02") + " 00:00:00"
	endTime := targetDate.Format("2006-01-02") + " 23:59:59"

	logx.Infof("查询日期: %s, 时间范围: %s ~ %s", dt, startTime, endTime)

	// 从数据库获取 access_token（在循环外查询一次，避免频繁查询数据库）
	mediaToken, err := model.GetByMedia(db, "juliang_dls")
	if err != nil {
		logx.Errorf("从数据库获取 juliang_dls token 失败: %v", err)
		return
	}

	// 从数据库获取所有客户及其媒体账户，平铺成任务列表
	tasks, err := collectElmHcTasks(db)
	if err != nil {
		logx.Errorf("%v", err)
		return
	}
	if len(tasks) == 0 {
		logx.Info("暂无客户/媒体账户配置")
		return
	}

	// 限频自动重试：失败账户按退避表重试，预算 40 分钟（11:00 启动，避开 ~11:20 拥塞窗口尾部且不与次日重叠）
	deadline := time.Now().Add(40 * time.Minute)
	allReportData := fetchElmHcReportsWithRetry(
		juliangConfig, mediaToken.Token, tasks, startTime, endTime, "stat_time_day", deadline,
		func(t elmHcTask, resp *types.JuliangCustomReportResp) []types.ADXReportData {
			return buildElmHcReportData(t.perf, t.media, resp, dt, "", false)
		},
	)

	// 发送数据到ADX接口（暂时注释，改为保存到数据库）
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

	// 保存数据到数据库
	if len(allReportData) > 0 {
		logx.Infof("准备保存 %d 条数据到数据库", len(allReportData))
		successCount := 0
		for _, data := range allReportData {
			record := &model.ElmHcReportData{
				CustomerName:      data.CustomerName,
				CustomerShort:     data.CustomerShort,
				AgentName:         data.AgentName,
				AgentShort:        data.AgentShort,
				MediaPlatformName: data.MediaPlatformName,
				// MediaAdvId:        data.MediaAdvId,
				// MediaAdvName:      data.MediaAdvName,
				HuichuanAdvId:   data.HuichuanAdvId,
				Cost:            data.Cost,
				ShowNum:         data.ShowNum,
				ClickNum:        data.ClickNum,
				ConvertNum:      data.ConvertNum,
				DeepConvertNum:  data.DeepConvertNum,
				ConvertType:     data.ConvertType,
				DeepConvertType: data.DeepConvertType,
				RedirectNum:     data.RedirectNum,
				PayNum:          data.PayNum,
				Dt:              data.Dt,
				Hh:              data.Hh,
			}
			if err := model.InsertOrUpdateElmHcReportData(db, record); err != nil {
				logx.Errorf("保存数据失败 (账户:%d, 日期:%s): %v", data.HuichuanAdvId, data.Dt, err)
			} else {
				successCount++
			}
		}
		logx.Infof("数据保存完成，成功 %d/%d 条", successCount, len(allReportData))
	} else {
		logx.Info("暂无数据需要保存")
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
		// 限频错误（code=40110，或 message 含“频率超限”兜底）包装为可重试的哨兵错误
		if resp.Code == 40110 || strings.Contains(resp.Message, "频率超限") {
			return nil, fmt.Errorf("%w: code=%d, message=%s", ErrJuliangRateLimit, resp.Code, resp.Message)
		}
		return nil, fmt.Errorf("获取报表失败: %s", resp.Message)
	}

	logx.Infof("成功获取巨量引擎报表数据：%v", resp.Data.Rows)
	return &resp, nil
}

// FetchHuichuanElmReportsByHour 获取回传饿了么所有账户的小时级报表数据并发送到ADX
// reportHour: 指定小时，格式 "2006010215"（如 "2026032614"），为空则取上一个小时
func FetchHuichuanElmReportsByHour(db *gorm.DB, juliangConfig config.JuliangConfig, adxConfig config.ADXConfig, reportHour string) {
	logx.Infof("开始获取回传饿了么小时级报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	var lastHour time.Time
	if reportHour != "" {
		t, err := time.ParseInLocation("2006010215", reportHour, time.Local)
		if err != nil {
			logx.Errorf("解析小时参数失败: %v", err)
			return
		}
		lastHour = t
	} else {
		// 获取当前时间，计算上一个小时的时间范围
		lastHour = time.Now().Add(-1 * time.Hour)
	}
	startTime := lastHour.Format("2006-01-02 15") + ":00:00"
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

	// 从数据库获取所有客户及其媒体账户，平铺成任务列表
	tasks, err := collectElmHcTasks(db)
	if err != nil {
		logx.Errorf("%v", err)
		return
	}
	if len(tasks) == 0 {
		logx.Info("暂无客户/媒体账户配置")
		return
	}

	// 限频自动重试：失败账户按退避表重试，预算 45 分钟（时报每小时跑一次，留足余量给下次 :02 启动，避免重叠）
	deadline := time.Now().Add(45 * time.Minute)
	allReportData := fetchElmHcReportsWithRetry(
		juliangConfig, mediaToken.Token, tasks, startTime, endTime, "stat_time_hour", deadline,
		func(t elmHcTask, resp *types.JuliangCustomReportResp) []types.ADXReportData {
			return buildElmHcReportData(t.perf, t.media, resp, dt, hh, true)
		},
	)

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

// FetchHuichuanElmReportsByDayHours 按天回传饿了么小时级报表：依次回传指定日期当天每个小时。
// reportDate: 指定日期，格式 "20060102"（如 "20260326"），为空则使用今天。
// 若 reportDate 为今天，则只回传 00 时到当前小时（未来小时无数据，跳过）。
func FetchHuichuanElmReportsByDayHours(db *gorm.DB, juliangConfig config.JuliangConfig, adxConfig config.ADXConfig, reportDate string) {
	logx.Infof("开始按天回传饿了么小时级报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	now := time.Now()
	var day time.Time
	if reportDate != "" {
		t, err := time.ParseInLocation("20060102", reportDate, time.Local)
		if err != nil {
			logx.Errorf("解析日期参数失败: %v", err)
			return
		}
		day = t
	} else {
		day = now
	}

	// 确定该天需要回传的最后一个小时：如果是今天，只到当前小时；否则到 23 时
	lastHour := 23
	if day.Format("20060102") == now.Format("20060102") {
		lastHour = now.Hour()
	}

	dt := day.Format("20060102")
	logx.Infof("按天回传日期: %s, 小时范围: 00 ~ %02d", dt, lastHour)

	for h := 0; h <= lastHour; h++ {
		reportHour := fmt.Sprintf("%s%02d", dt, h)
		logx.Infof("=== 按天回传：开始第 %02d 小时 (%s) ===", h, reportHour)
		FetchHuichuanElmReportsByHour(db, juliangConfig, adxConfig, reportHour)
	}

	logx.Infof("按天回传饿了么小时级报表数据全部完成，日期: %s, 共 %d 个小时", dt, lastHour+1)
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
		logx.Infof("hour input: %v", string(batchJSON))

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
