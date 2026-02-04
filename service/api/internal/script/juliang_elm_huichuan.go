package script

import (
	"context"
	"encoding/json"
	"fmt"
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

	logx.Infof("巨量DLS Token 刷新成功，新 AccessToken: %s, 有效期: %d 秒", resp.Data.AccessToken, resp.Data.ExpiresIn)
	logx.Infof("新 RefreshToken: %s, 有效期: %d 秒", resp.Data.RefreshToken, resp.Data.RefreshTokenExpiresIn)
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

	logx.Infof("巨量KH Token 刷新成功，新 AccessToken: %s, 有效期: %d 秒", resp.Data.AccessToken, resp.Data.ExpiresIn)
	logx.Infof("新 RefreshToken: %s, 有效期: %d 秒", resp.Data.RefreshToken, resp.Data.RefreshTokenExpiresIn)
}

// FetchHuichuanElmReports 获取回传饿了么所有账户的报表数据
func FetchHuichuanElmReports(db *gorm.DB, juliangConfig config.JuliangConfig, huichuanElmConfig config.HuichuanElmConfig) {
	logx.Infof("开始获取回传饿了么报表数据 - %s", time.Now().Format("2006-01-02 15:04:05"))

	// 获取今天的开始和结束时间
	now := time.Now()
	startTime := now.Format("2006-01-02") + " 00:00:00"
	endTime := now.Format("2006-01-02") + " 23:59:59"

	logx.Infof("查询时间范围: %s ~ %s", startTime, endTime)

	// 遍历所有广告主ID
	for _, advertiserId := range huichuanElmConfig.AdvertiserIds {
		logx.Infof("正在获取广告主 %d 的报表数据...", advertiserId)

		resp, err := getJuliangReportData(db, juliangConfig, advertiserId, startTime, endTime)
		if err != nil {
			logx.Errorf("获取广告主 %d 的报表数据失败: %v", advertiserId, err)
			continue
		}

		// 处理返回的数据
		logx.Infof("广告主 %d 报表数据获取成功:", advertiserId)
		logx.Infof("  总记录数: %d", len(resp.Data.Rows))
		if len(resp.Data.Rows) > 0 {
			logx.Infof("  总指标: %+v", resp.Data.TotalMetrics)
		}

		// TODO: 这里可以添加保存数据到数据库或其他处理逻辑
	}

	logx.Infof("回传饿了么报表数据获取完成")
}

// getJuliangReportData 获取巨量引擎报表数据
func getJuliangReportData(db *gorm.DB, juliangConfig config.JuliangConfig, advertiserId int64, startTime, endTime string) (*types.JuliangCustomReportResp, error) {
	ctx := context.Background()
	logx.Infof("开始获取巨量引擎报表数据 - advertiser_id: %d, 时间范围: %s ~ %s", advertiserId, startTime, endTime)

	mediaType := "juliang_dls"

	// 从数据库获取 access_token
	mediaToken, err := model.GetByMedia(db, mediaType)
	if err != nil {
		logx.Errorf("从数据库获取 %s token 失败: %v", mediaType, err)
		return nil, err
	}

	// 创建 HTTP 客户端
	client := httpclient.NewClient(juliangConfig.BaseUrl, juliangConfig.Timeout)
	client.SetHeader("Access-Token", mediaToken.Token)

	// 构建查询参数（需要序列化为JSON字符串）
	dimensions := []string{
		"stat_time_day",
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
			Field: "stat_time_day",
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
	err = client.Get(ctx, "/open_api/v3.0/report/custom/get/", params, &resp)
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
