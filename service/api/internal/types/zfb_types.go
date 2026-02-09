package types

// ZfbDownloadReq ZFB 下载请求参数
type ZfbDownloadReq struct {
	Uid       string `form:"uid"`                 // 用户ID (2088开头的uid)
	StartDate string `form:"start_date,optional"` // 开始日期 YYYY-MM-DD (可选，默认当天)
	EndDate   string `form:"end_date,optional"`   // 结束日期 YYYY-MM-DD (可选，默认当天)
}

// ZfbDownloadResp ZFB 下载响应
type ZfbDownloadResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	FileUrl string `json:"file_url,omitempty"` // 文件下载URL
}

// 支付宝接口相关数据结构
type AlipayDataMeta struct {
	ColumnID    string `json:"columnId"`
	ColumnName  string `json:"columnName"`
	DisplayName string `json:"displayName"`
	DataType    string `json:"dataType"`
}

type AlipaySyncResult struct {
	DataMetaMap map[string]AlipayDataMeta `json:"dataMetaMap"`
	DataValue   []map[string]interface{}  `json:"dataValue"`
}

type AlipayApiResponse struct {
	Success bool `json:"success"`
	Data    struct {
		SyncResult []AlipaySyncResult `json:"syncResult"`
	} `json:"data"`
	ErrorDesc string `json:"errorDesc"`
}

// 更新支付宝Cookie请求
type UpdateZfbCookieReq struct {
	Token        string `json:"cookie"`
	RefreshToken string `json:"csrfToken"`
}

// 更新支付宝Cookie响应
type UpdateZfbCookieResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
