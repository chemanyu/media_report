package types

// JDAttributionRequest 查询京东归因数据的请求
type JDAttributionRequest struct {
	Date string `form:"date"` // 日期格式：YYYYMMDD 或 YYYY-MM-DD
}

// JDAttributionResponse 归因数据的响应
type JDAttributionResponse struct {
	Date          string                      `json:"date"`
	TotalRequests map[string]int64            `json:"total_requests"` // 每个账户的总请求数
	ErrorCounts   map[string]map[string]int64 `json:"error_counts"`   // 账户ID -> 错误类型 -> 错误数量
}

// JDAttributionAPIResponse 归因API的原始响应
type JDAttributionAPIResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    JDAttributionAPIData `json:"data"`
}

// JDAttributionAPIData API返回的数据部分
type JDAttributionAPIData struct {
	Date          string                      `json:"date"`
	TotalRequests map[string]int64            `json:"total_requests"`
	ErrorCounts   map[string]map[string]int64 `json:"error_counts"`
	Summary       struct {
		TotalRequestCount int64            `json:"total_request_count"`
		TotalErrorCount   int64            `json:"total_error_count"`
		ErrorTypes        map[string]int64 `json:"error_types"`
	} `json:"summary"`
}

// JDErrorExportRequest 导出错误统计数据的请求
type JDErrorExportRequest struct {
	NumDays int `form:"num_days,default=10"` // 导出过去多少天的数据，默认10天
}

// JDErrorCountsData 错误统计数据
type JDErrorCountsData struct {
	Date         string `json:"date"`
	AdvertiserId string `json:"advertiser_id"`
	ErrorType    string `json:"error_type"`
	Event        string `json:"event"`
	ErrorCount   int64  `json:"error_count"`
}
