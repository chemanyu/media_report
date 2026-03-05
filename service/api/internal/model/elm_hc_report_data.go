package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ElmHcReportData 饿了么汇川报表数据
type ElmHcReportData struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerName      string    `gorm:"column:customer_name;type:varchar(200);not null" json:"customer_name"`             // 客户名称
	CustomerShort     string    `gorm:"column:customer_short;type:varchar(50);not null" json:"customer_short"`            // 客户简称
	AgentName         string    `gorm:"column:agent_name;type:varchar(200);not null" json:"agent_name"`                   // 代理名称
	AgentShort        string    `gorm:"column:agent_short;type:varchar(50);not null" json:"agent_short"`                  // 代理简称
	MediaPlatformName string    `gorm:"column:media_platform_name;type:varchar(100);not null" json:"media_platform_name"` // 媒体平台名称
	MediaAdvId        string    `gorm:"column:media_adv_id;type:varchar(100);not null" json:"media_adv_id"`               // 媒体账户ID
	MediaAdvName      string    `gorm:"column:media_adv_name;type:varchar(200);not null" json:"media_adv_name"`           // 媒体账户名称
	HuichuanAdvId     int64     `gorm:"column:huichuan_adv_id;not null;default:0" json:"huichuan_adv_id"`                 // 汇川账户ID
	Cost              float64   `gorm:"column:cost;default:0" json:"cost"`                                                // 消费
	ShowNum           int64     `gorm:"column:show_num;default:0" json:"show_num"`                                        // 展示数
	ClickNum          int64     `gorm:"column:click_num;default:0" json:"click_num"`                                      // 点击数
	ConvertNum        int64     `gorm:"column:convert_num;default:0" json:"convert_num"`                                  // 转化数
	DeepConvertNum    int       `gorm:"column:deep_convert_num;default:0" json:"deep_convert_num"`                        // 深度转化数
	ConvertType       string    `gorm:"column:convert_type;type:varchar(50);default:''" json:"convert_type"`              // 转化类型
	DeepConvertType   string    `gorm:"column:deep_convert_type;type:varchar(50);default:''" json:"deep_convert_type"`    // 深度转化类型
	RedirectNum       int       `gorm:"column:redirect_num;default:0" json:"redirect_num"`                                // 调起数
	PayNum            int       `gorm:"column:pay_num;default:0" json:"pay_num"`                                          // 付费数
	Dt                string    `gorm:"column:dt;type:varchar(8);not null" json:"dt"`                                     // 日期，格式：yyyyMMdd
	Hh                string    `gorm:"column:hh;type:varchar(2);not null;default:''" json:"hh"`                          // 小时，格式：hh
	CreateTime        time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ElmHcReportData) TableName() string {
	return "elm_hc_report_data"
}

// QueryElmHcReportDataParams 查询参数
type QueryElmHcReportDataParams struct {
	StartDate     string // 开始日期，格式：yyyyMMdd
	EndDate       string // 结束日期，格式：yyyyMMdd
	CustomerShort string // 客户简称，可选
}

// QueryElmHcReportData 查询饿了么汇川报表数据
func QueryElmHcReportData(db *gorm.DB, params QueryElmHcReportDataParams) ([]ElmHcReportData, error) {
	var records []ElmHcReportData
	query := db.Model(&ElmHcReportData{})

	if params.StartDate != "" {
		query = query.Where("dt >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("dt <= ?", params.EndDate)
	}
	if params.CustomerShort != "" {
		query = query.Where("customer_short = ?", params.CustomerShort)
	}

	err := query.Order("dt DESC, customer_short ASC, media_adv_id ASC").Find(&records).Error
	return records, err
}

// InsertOrUpdateElmHcReportData 插入或更新饿了么汇川报表数据（Upsert）
func InsertOrUpdateElmHcReportData(db *gorm.DB, record *ElmHcReportData) error {
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "media_adv_id"},
			{Name: "dt"},
			{Name: "hh"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"customer_name":       record.CustomerName,
			"customer_short":      record.CustomerShort,
			"agent_name":          record.AgentName,
			"agent_short":         record.AgentShort,
			"media_platform_name": record.MediaPlatformName,
			"media_adv_name":      record.MediaAdvName,
			"huichuan_adv_id":     record.HuichuanAdvId,
			"cost":                record.Cost,
			"show_num":            record.ShowNum,
			"click_num":           record.ClickNum,
			"convert_num":         record.ConvertNum,
			"deep_convert_num":    record.DeepConvertNum,
			"convert_type":        record.ConvertType,
			"deep_convert_type":   record.DeepConvertType,
			"redirect_num":        record.RedirectNum,
			"pay_num":             record.PayNum,
			"update_time":         gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(record).Error

	if err != nil {
		return fmt.Errorf("插入或更新饿了么报表数据失败: %w", err)
	}
	return nil
}
