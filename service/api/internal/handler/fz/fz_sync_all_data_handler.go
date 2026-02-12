package fz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	script "media_report/service/api/internal/script"
	"media_report/service/api/internal/svc"
)

// FzSyncAllDataHandler 同步飞猪所有媒体数据
func FzSyncAllDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取日期参数（可选）
		reportDate := r.URL.Query().Get("date") // 格式: 20260211

		l := fz.NewFzHourlyReportLogic(r.Context(), svcCtx)

		var oppoCount, xiaomiCount int
		var oppoErr, xiaomiErr error

		// 同步OPPO数据
		if reportDate != "" {
			oppoCount, oppoErr = l.SyncOppoData(reportDate)
		} else {
			oppoCount, oppoErr = l.SyncTodayOppoData()
		}

		// 同步小米数据
		if reportDate != "" {
			xiaomiCount, xiaomiErr = l.SyncXiaomiData(reportDate)
		} else {
			xiaomiCount, xiaomiErr = l.SyncTodayXiaomiData()
		}

		// 构建响应
		result := map[string]interface{}{
			"success": true,
			"message": "同步完成",
			"data": map[string]interface{}{
				"oppo": map[string]interface{}{
					"count":   oppoCount,
					"success": oppoErr == nil,
					"error":   "",
				},
				"xiaomi": map[string]interface{}{
					"count":   xiaomiCount,
					"success": xiaomiErr == nil,
					"error":   "",
				},
			},
		}

		// 添加错误信息
		if oppoErr != nil {
			result["data"].(map[string]interface{})["oppo"].(map[string]interface{})["error"] = oppoErr.Error()
		}
		if xiaomiErr != nil {
			result["data"].(map[string]interface{})["xiaomi"].(map[string]interface{})["error"] = xiaomiErr.Error()
		}

		// 如果都失败了，返回错误
		if oppoErr != nil && xiaomiErr != nil {
			result["success"] = false
			result["message"] = "同步失败"
		} else if oppoErr != nil || xiaomiErr != nil {
			result["message"] = "部分同步成功"
		}

		httpx.OkJsonCtx(r.Context(), w, result)

		// 实际调用（请确保顶部已import script包）
		script.SendFzDingTalkNotification(r.Context(), svcCtx.DB, svcCtx.Config.DingTalk)
	}
}
