package types

// TanxFetchDataReq 抓取淘宝联盟数据请求
type TanxFetchDataReq struct {
}

// TanxFetchDataResp 抓取淘宝联盟数据响应
type TanxFetchDataResp struct {
	Message string `json:"message"`
}

// TanxExportDataReq 导出数据请求
type TanxExportDataReq struct {
}

// TanxExportDataResp 导出数据响应
type TanxExportDataResp struct {
	Message string `json:"message"`
}

// TanxUpdateCookieReq 更新Cookie请求
type TanxUpdateCookieReq struct {
	Cookie string `json:"cookie" binding:"required"`
}

// TanxUpdateCookieResp 更新Cookie响应
type TanxUpdateCookieResp struct {
	Message    string `json:"message"`
	UpdateTime string `json:"update_time"`
}

// TanxReportItem 淘宝联盟API返回的数据项
type TanxReportItem struct {
	Ds            string `json:"ds"`            // 日期
	Pid           string `json:"pid"`           // 广告位ID
	AdzoneName    string `json:"adzoneName"`    // 广告位名称
	Qingqiupv     string `json:"qingqiupv"`     // tanx有效请求（API返回字符串）
	ActiveRatioDf string `json:"activeRatioDf"` // 东风手淘换端率-同步点击
	TanxEffectPv  string `json:"tanxEffectPv"`  // TANX曝光数（API返回字符串）
	TanxClk       string `json:"tanxClk"`       // TANX点击数（API返回字符串）
	DongfengEf    string `json:"dongfengEf"`    // TANX预估收益（API返回字符串）
}

// TanxAPIResponse 淘宝联盟API响应结构
type TanxAPIResponse struct {
	Info struct {
		Ok        bool   `json:"ok"`
		ErrorCode string `json:"errorCode"`
		Msg       string `json:"msg"`
	} `json:"info"`
	Data struct {
		ClickMonitor          bool             `json:"clickMonitor"`
		ClickMonitorParamList []TanxReportItem `json:"clickMonitorParamList"`
		ClickParamList        []TanxReportItem `json:"clickParamList"`
		PvParamList           []TanxReportItem `json:"pvParamList"`
		ReportList            []TanxReportItem `json:"reportList"`
		ParamCheckResult      string           `json:"paramCheckResult"`
		Pids                  []string         `json:"pids"`
		RequestId             string           `json:"requestId"`
		ResourceId            []interface{}    `json:"resourceId"`
	} `json:"data"`
}
