package types

// =============== sync_adn_data ===============

type QczjSyncDataReq struct {
	ReportDate string `json:"report_date"`
	View       int64  `json:"view"`
	Click      int64  `json:"click"`
	Action     int64  `json:"action"`
}

// =============== config ===============

type QczjConfigResp struct {
	Id         uint    `json:"id"`
	TotalUv    int64   `json:"total_uv"`
	Ratio      float64 `json:"ratio"`
	UpdateTime string  `json:"update_time"`
}

type QczjUpdateConfigReq struct {
	TotalUv int64   `json:"total_uv"`
	Ratio   float64 `json:"ratio"`
}

// =============== report ===============

type QczjReportItem struct {
	Id         int64  `json:"id"`
	ReportDate int    `json:"report_date"`
	View       int64  `json:"view"`
	Click      int64  `json:"click"`
	Action     int64  `json:"action"`
	UpdateTime string `json:"update_time"`
}

type QczjReportListResp struct {
	Total int64             `json:"total"`
	List  []*QczjReportItem `json:"list"`
}
