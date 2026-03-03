package types

// ==================== 淘宝转链 ====================

// TaobaoExtractReq 淘宝单链提取请求
type TaobaoExtractReq struct {
	ShortUrl string `json:"short_url"` // 淘宝短链
	Platform string `json:"platform"`  // 平台: ios / android（默认 ios）
}

// TaobaoExtractResp 淘宝单链提取响应
type TaobaoExtractResp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Deeplink string `json:"deeplink"` // tbopen:// 或 taobao:// 深链
	H5Dp     string `json:"h5_dp"`    // iOS 平台的 ace.tb.cn 短链
}

// ==================== 闲鱼转链 ====================

// XianyuExtractReq 闲鱼单链提取请求
type XianyuExtractReq struct {
	ShortUrl string `json:"short_url"` // 闲鱼短链
	Platform string `json:"platform"`  // 平台: android / ios（默认 android）
}

// XianyuExtractResp 闲鱼单链提取响应
type XianyuExtractResp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Deeplink string `json:"deeplink"` // fleamarket:// 或 pages.goofish.com 链接
}

// ==================== 批量操作通用响应 ====================

// BatchResp 批量操作响应（文件上传 → 下载 Excel）
type BatchResp struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	DownloadUrl string `json:"download_url"` // Excel 下载地址，如 /download/xxx.xlsx
	Success     int    `json:"success"`      // 成功条数
	Failed      int    `json:"failed"`       // 失败条数
}

// ==================== 活动短链批量获取 ====================

// ActivityBatchReq 活动短链批量获取请求（文件通过 multipart form 传递）
type ActivityBatchReq struct {
	ActivityMaterialId string `json:"activity_material_id" form:"activity_material_id"` // 活动素材ID
}

// ==================== 活动报表查询 ====================

// ActivityReportReq 活动报表批量查询请求（Excel 文件通过 multipart form 传递）
type ActivityReportReq struct {
	BizDate   string `json:"biz_date" form:"biz_date"`     // 统计日期，格式 yyyyMMdd
	QueryType string `json:"query_type" form:"query_type"` // 查询类型: 1=预估数据, 2=结算数据
	EventId   string `json:"event_id" form:"event_id"`     // 活动ID: 4439384=福利购, 4297311=超级红包
}
