package tanx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"media_report/service/api/internal/model"
	"media_report/service/api/internal/svc"
	"media_report/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchDataLogic {
	return &FetchDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FetchData 抓取淘宝联盟数据
func (l *FetchDataLogic) FetchData(req *types.TanxFetchDataReq) (resp *types.TanxFetchDataResp, err error) {
	// 获取全局配置
	config := GetTanxConfig()

	// 从数据库获取 cookie
	cookie, _, err := model.GetTokensByMedia(l.svcCtx.DB, "tanx_pachong")
	if err != nil {
		return nil, fmt.Errorf("获取 Cookie 失败: %w", err)
	}
	if cookie == "" {
		return nil, fmt.Errorf("Cookie 为空，请先更新 Cookie")
	}

	// 计算查询日期（默认查前1天）
	queryDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	l.Logger.Infof("查询日期: %s", queryDate)

	url := "https://tanx.alimama.com/api/media/debug/report/getReport.htm"

	// 遍历广告位列表抓取数据
	successCount := 0
	failCount := 0
	for _, pid := range config.AdSlots {
		if err := l.fetchDataForPid(url, cookie, queryDate, pid); err != nil {
			l.Logger.Errorf("抓取广告位 %s 数据失败: %v", pid, err)
			failCount++
		} else {
			successCount++
		}
	}

	message := fmt.Sprintf("数据抓取完成！成功: %d, 失败: %d", successCount, failCount)
	l.Logger.Info(message)

	return &types.TanxFetchDataResp{
		Message: message,
	}, nil
}

// fetchDataForPid 抓取单个广告位的数据
func (l *FetchDataLogic) fetchDataForPid(url, cookie, queryDate, pid string) error {
	// 构造请求数据
	requestData := map[string]interface{}{
		"ds":         queryDate,
		"mediaClick": "1",
		"mediaCost":  "1",
		"mediaPV":    "1",
		"pid":        pid,
		"type":       "rtb",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Cookie", cookie)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	//l.Logger.Info("apiResp: %v", string(body))
	var apiResp types.TanxAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	//l.Logger.Info("apiResp: %v", apiResp.Data)

	// 检查是否需要登录
	if apiResp.Info.ErrorCode == "user_not_login" {
		l.Logger.Error("用户未登录，请更新Cookie")
		return fmt.Errorf("用户未登录，请更新Cookie")
	}

	// 检查是否有数据
	if len(apiResp.Data.ClickMonitorParamList) == 0 {
		l.Logger.Infof("广告位 %s 在 %s 没有数据", pid, queryDate)
		return nil
	}

	// 保存数据到数据库
	for _, item := range apiResp.Data.ClickMonitorParamList {
		// 转换日期格式：20260208 -> 2006-01-02
		ds := item.Ds
		if len(ds) == 8 {
			// 格式化为 YYYY-MM-DD
			ds = fmt.Sprintf("%s-%s-%s", ds[0:4], ds[4:6], ds[6:8])
		}

		// 转换 activeRatioDf 为百分比格式
		activeRatioDf := item.ActiveRatioDf
		if activeRatioDf == "" {
			activeRatioDf = "0"
		}

		ratio := 0.0
		fmt.Sscanf(activeRatioDf, "%f", &ratio)
		activeRatioDfPercent := fmt.Sprintf("%.2f%%", ratio*100)

		// 转换字符串到数值类型
		var qingqiupv, tanxEffectPv, tanxClk int64
		var dongfengEf float64

		fmt.Sscanf(item.Qingqiupv, "%d", &qingqiupv)
		fmt.Sscanf(item.TanxEffectPv, "%d", &tanxEffectPv)
		fmt.Sscanf(item.TanxClk, "%d", &tanxClk)
		fmt.Sscanf(item.DongfengEf, "%f", &dongfengEf)

		if err := l.insertData(ds, item.Pid, item.AdzoneName, qingqiupv,
			activeRatioDfPercent, tanxEffectPv, tanxClk, dongfengEf); err != nil {
			l.Logger.Errorf("保存数据失败: %v", err)
		} else {
			l.Logger.Infof("保存广告位 %s 数据成功", pid)
		}
	}

	return nil
}

// insertData 插入或更新数据到数据库
func (l *FetchDataLogic) insertData(ds, pid, adzoneName string, qingqiupv int64,
	activeRatioDf string, tanxEffectPv, tanxClk int64, dongfengEf float64) error {

	// 使用 GORM 模型创建或更新数据
	monitor := &model.TanxMonitor{
		Ds:            ds,
		Pid:           pid,
		AdzoneName:    adzoneName,
		Qingqiupv:     qingqiupv,
		ActiveRatioDf: activeRatioDf,
		TanxEffectPv:  tanxEffectPv,
		TanxClk:       tanxClk,
		DongfengEf:    dongfengEf,
	}

	return model.UpsertTanxMonitor(l.svcCtx.DB, monitor)
}
