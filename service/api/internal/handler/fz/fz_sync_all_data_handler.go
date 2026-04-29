package fz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	fz "media_report/service/api/internal/logic/fz"
	"media_report/service/api/internal/svc"
)

// FzSyncAllDataHandler 同步飞猪所有媒体数据
func FzSyncAllDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取日期参数（可选）
		reportDate := r.URL.Query().Get("date") // 格式: 20260211

		l := fz.NewFzHourlyReportLogic(r.Context(), svcCtx)

		var oppoCount, xiaomiCount, honorCount int
		var oppoErr, xiaomiErr, honorErr error

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

		// 同步荣耀数据
		if reportDate != "" {
			honorCount, honorErr = l.SyncHonorData(reportDate)
		} else {
			honorCount, honorErr = l.SyncTodayHonorData()
		}

		oppoErrMsg := ""
		if oppoErr != nil {
			oppoErrMsg = oppoErr.Error()
		}
		xiaomiErrMsg := ""
		if xiaomiErr != nil {
			xiaomiErrMsg = xiaomiErr.Error()
		}
		honorErrMsg := ""
		if honorErr != nil {
			honorErrMsg = honorErr.Error()
		}

		// 构建响应
		result := map[string]interface{}{
			"success": true,
			"message": "同步完成",
			"data": map[string]interface{}{
				"oppo": map[string]interface{}{
					"count":   oppoCount,
					"success": oppoErr == nil,
					"error":   oppoErrMsg,
				},
				"xiaomi": map[string]interface{}{
					"count":   xiaomiCount,
					"success": xiaomiErr == nil,
					"error":   xiaomiErrMsg,
				},
				"honor": map[string]interface{}{
					"count":   honorCount,
					"success": honorErr == nil,
					"error":   honorErrMsg,
				},
			},
		}

		if oppoErr != nil && xiaomiErr != nil && honorErr != nil {
			result["success"] = false
			result["message"] = "同步失败"
		} else if oppoErr != nil || xiaomiErr != nil || honorErr != nil {
			result["message"] = "部分同步成功"
		}

		httpx.OkJsonCtx(r.Context(), w, result)

		// 实际调用（请确保顶部已import script包）
		//script.SendFzDingTalkNotification(r.Context(), svcCtx.DB, svcCtx.Config.DingTalk)
		//script.SendFzDailyReport(r.Context(), svcCtx.DB, svcCtx.Config.DingTalk)
	}
}
