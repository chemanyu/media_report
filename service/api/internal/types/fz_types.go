package types

// =============== 飞猪媒体账户 ===============

// 添加飞猪媒体账户请求
type FzAdvertiserAddReq struct {
	Media        string `json:"media" binding:"required"`          // 媒体类型：oppo, xiaomi, adn
	MediaAdvId   string `json:"media_adv_id" binding:"required"`   // 媒体账户ID
	MediaAdvName string `json:"media_adv_name" binding:"required"` // 媒体账户名称
}

// 更新飞猪媒体账户请求
type FzAdvertiserUpdateReq struct {
	Id           int64  `json:"id" binding:"required"`             // 账户ID
	MediaAdvName string `json:"media_adv_name" binding:"required"` // 媒体账户名称
}

// 删除飞猪媒体账户请求
type FzAdvertiserDeleteReq struct {
	Id int64 `json:"id" binding:"required"` // 账户ID
}

// =============== 飞猪小时报 ===============

// 飞猪小时报列表查询请求
type FzHourlyReportListReq struct {
	Media     string `form:"media,optional"`      // 媒体类型（可选）
	StartDate string `form:"start_date,optional"` // 开始日期 格式：20260211
	EndDate   string `form:"end_date,optional"`   // 结束日期 格式：20260211
}
