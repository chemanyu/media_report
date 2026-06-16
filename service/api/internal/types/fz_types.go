package types

// =============== 飞猪媒体账户 ===============

// 添加飞猪媒体账户请求
type FzAdvertiserAddReq struct {
	Media        string `json:"media" binding:"required"`          // 媒体类型：oppo, xiaomi, adn
	MediaAdvId   string `json:"media_adv_id" binding:"required"`   // 媒体账户ID
	MediaAdvName string `json:"media_adv_name" binding:"required"` // 媒体账户名称
	ClientID     string `json:"client_id"`                         // OAuth2 Client ID（honor专用）
	ClientSecret string `json:"client_secret"`                     // OAuth2 Client Secret（honor专用）
}

// 更新飞猪媒体账户请求
type FzAdvertiserUpdateReq struct {
	Id           int64  `json:"id" binding:"required"`             // 账户ID
	MediaAdvName string `json:"media_adv_name" binding:"required"` // 媒体账户名称
	ClientID     string `json:"client_id"`                         // OAuth2 Client ID（honor专用）
	ClientSecret string `json:"client_secret"`                     // OAuth2 Client Secret（honor专用）
}

// 删除飞猪媒体账户请求
type FzAdvertiserDeleteReq struct {
	Id int64 `json:"id" binding:"required"` // 账户ID
}

// ADN账户信息
type FzAdnAdvertiserItem struct {
	Id           int64  `json:"id"`
	MediaAdvId   string `json:"media_adv_id"`
	MediaAdvName string `json:"media_adv_name"`
}

// 获取ADN账户列表响应
type FzAdnAdvertiserListResp struct {
	List []*FzAdnAdvertiserItem `json:"list"`
}

// =============== 飞猪小时报 ===============

// 飞猪小时报列表查询请求
type FzHourlyReportListReq struct {
	Media     string `form:"media,optional"`      // 媒体类型（可选）
	StartDate string `form:"start_date,optional"` // 开始日期 格式：20260211
	EndDate   string `form:"end_date,optional"`   // 结束日期 格式：20260211
}

// =============== 飞猪系数配置 ===============

// 飞猪系数配置响应
type FzConfigResp struct {
	Id          uint    `json:"id"`
	Coefficient float64 `json:"coefficient"` // 系数
	BaseNum     float64 `json:"base_num"`    // 基数
	UpdateTime  string  `json:"update_time"`
}

// 更新飞猪系数配置请求
type FzUpdateConfigReq struct {
	Coefficient float64 `json:"coefficient"` // 系数
	BaseNum     float64 `json:"base_num"`    // 基数
}

// ADN数据同步请求
type FzSyncAdnDataReq struct {
	MediaAdvId      string  `json:"media_adv_id"`       // 媒体账户ID
	MediaAdvName    string  `json:"media_adv_name"`     // 媒体账户名称
	ReportDate      string  `json:"report_date"`        // 报表日期，格式: 20260211
	Cost            float64 `json:"cost"`               // 消耗（单位：分）
	ConvertDp       int64   `json:"convert_dp"`         // 拉活数
	DpAppOrderNums  int64   `json:"dp_app_order_nums"`  // 订单数
	Click           int64   `json:"click"`              // 点击数
	Expose          int64   `json:"expose"`             // 曝光数
	ConvertDpPrice  float64 `json:"convert_dp_price"`   // 拉活成本（单位：分）
	DpAppOrderPrice float64 `json:"dp_app_order_price"` // 订单成本（单位：分）
}

// 华为数据同步请求（字段与外部回传一致，cost/成本单位为元，入库时 ×100 转分）
type FzSyncHuaweiDataReq struct {
	MediaAdvId      string  `json:"media_adv_id"`       // 媒体账户ID
	MediaAdvName    string  `json:"media_adv_name"`     // 媒体账户名称
	ReportDate      string  `json:"report_date"`        // 报表日期，格式: 20260211
	Cost            float64 `json:"cost"`               // 消耗（单位：元）
	ConvertDp       int64   `json:"convert_dp"`         // 拉活数
	DpAppOrderNums  int64   `json:"dp_app_order_nums"`  // 订单数
	Click           int64   `json:"click"`              // 点击数
	Expose          int64   `json:"expose"`             // 曝光数
	ConvertDpPrice  float64 `json:"convert_dp_price"`   // 拉活成本（单位：元）
	DpAppOrderPrice float64 `json:"dp_app_order_price"` // 订单成本（单位：元）
}
