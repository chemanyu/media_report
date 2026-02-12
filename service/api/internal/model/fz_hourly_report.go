package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FzHourlyReport 飞猪时报数据表
type FzHourlyReport struct {
	Id              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Media           string    `gorm:"column:media" json:"media"`                           // 媒体简称
	MediaAdvId      string    `gorm:"column:media_adv_id" json:"media_adv_id"`             // 媒体账户ID
	MediaAdvName    string    `gorm:"column:media_adv_name" json:"media_adv_name"`         // 媒体账户名称
	ReportDate      int       `gorm:"column:report_date" json:"report_date"`               // 报表日期
	Cost            float64   `gorm:"column:cost" json:"cost"`                             // 消耗
	ConvertDp       int64     `gorm:"column:convert_dp" json:"convert_dp"`                 // 拉活数
	DpAppOrderNums  int64     `gorm:"column:dp_app_order_nums" json:"dp_app_order_nums"`   // 订单数
	Click           int64     `gorm:"column:click" json:"click"`                           // 点击数
	Expose          int64     `gorm:"column:expose" json:"expose"`                         // 曝光数
	ConvertDpPrice  float64   `gorm:"column:convert_dp_price" json:"convert_dp_price"`     // 拉活成本
	DpAppOrderPrice float64   `gorm:"column:dp_app_order_price" json:"dp_app_order_price"` // 订单成本
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (FzHourlyReport) TableName() string {
	return "fz_hourly_report"
}

// FzHourlyReportModel 飞猪时报数据访问接口
type FzHourlyReportModel interface {
	// InsertOrUpdate 插入或更新报表数据
	InsertOrUpdate(report *FzHourlyReport) error
	// FindByMediaAndDate 根据媒体类型和日期查询
	FindByMediaAndDate(media string, reportDate int) ([]*FzHourlyReport, error)
	// FindByMediaAccountAndDate 根据媒体账户和日期查询
	FindByMediaAccountAndDate(media, mediaAdvId string, reportDate int) (*FzHourlyReport, error)
	// FindByDateRange 根据日期范围查询
	FindByDateRange(media string, startDate, endDate int) ([]*FzHourlyReport, error)
	// FindAll 查询所有数据
	FindAll() ([]*FzHourlyReport, error)
}

type defaultFzHourlyReportModel struct {
	db *gorm.DB
}

// NewFzHourlyReportModel 创建时报模型实例
func NewFzHourlyReportModel(db *gorm.DB) FzHourlyReportModel {
	return &defaultFzHourlyReportModel{
		db: db,
	}
}

// InsertOrUpdate 插入或更新报表数据（使用Upsert）
func (m *defaultFzHourlyReportModel) InsertOrUpdate(report *FzHourlyReport) error {
	// 使用Clauses实现ON DUPLICATE KEY UPDATE
	err := m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "media"}, {Name: "media_adv_id"}, {Name: "report_date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"media_adv_name":     report.MediaAdvName,
			"cost":               report.Cost,
			"convert_dp":         report.ConvertDp,
			"dp_app_order_nums":  report.DpAppOrderNums,
			"click":              report.Click,
			"expose":             report.Expose,
			"convert_dp_price":   report.ConvertDpPrice,
			"dp_app_order_price": report.DpAppOrderPrice,
			"update_time":        gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(report).Error

	if err != nil {
		return fmt.Errorf("插入或更新时报数据失败: %w", err)
	}

	return nil
}

// FindByMediaAndDate 根据媒体类型和日期查询
func (m *defaultFzHourlyReportModel) FindByMediaAndDate(media string, reportDate int) ([]*FzHourlyReport, error) {
	var reports []*FzHourlyReport
	err := m.db.Where("media = ? AND report_date = ?", media, reportDate).Find(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("查询时报数据失败: %w", err)
	}
	return reports, nil
}

// FindByMediaAccountAndDate 根据媒体账户和日期查询
func (m *defaultFzHourlyReportModel) FindByMediaAccountAndDate(media, mediaAdvId string, reportDate int) (*FzHourlyReport, error) {
	var report FzHourlyReport
	err := m.db.Where("media = ? AND media_adv_id = ? AND report_date = ?", media, mediaAdvId, reportDate).First(&report).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("时报数据不存在")
		}
		return nil, fmt.Errorf("查询时报数据失败: %w", err)
	}

	return &report, nil
}

// FindByDateRange 根据日期范围查询
func (m *defaultFzHourlyReportModel) FindByDateRange(media string, startDate, endDate int) ([]*FzHourlyReport, error) {
	var reports []*FzHourlyReport
	query := m.db.Order("report_date DESC, id DESC")

	if media != "" {
		query = query.Where("media = ?", media)
	}

	if startDate > 0 && endDate > 0 {
		query = query.Where("report_date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate > 0 {
		query = query.Where("report_date >= ?", startDate)
	} else if endDate > 0 {
		query = query.Where("report_date <= ?", endDate)
	}

	err := query.Find(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("查询时报数据失败: %w", err)
	}
	return reports, nil
}

// FindAll 查询所有数据
func (m *defaultFzHourlyReportModel) FindAll() ([]*FzHourlyReport, error) {
	var reports []*FzHourlyReport
	err := m.db.Order("report_date DESC, id DESC").Find(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("查询所有时报数据失败: %w", err)
	}
	return reports, nil
}
