package types

// ZfbDownloadReq ZFB 下载请求参数
type ZfbDownloadReq struct {
	FileType string `form:"file_type"` // 文件类型：如 report, data 等
	Date     string `form:"date"`      // 日期参数
	// 其他需要的参数，根据你的业务需求添加
}

// ZfbDownloadResp ZFB 下载响应
type ZfbDownloadResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	FileUrl string `json:"file_url,omitempty"` // 文件下载URL
	// 或者直接返回文件流，根据你的需求决定
}
