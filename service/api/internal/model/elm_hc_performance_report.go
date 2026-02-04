package model

import (
	"time"

	"gorm.io/gorm"
)

// ElmHcPerformanceReport 汇川饿了么数据报表
type ElmHcPerformanceReport struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerName      string    `gorm:"column:customer_name;type:varchar(200);not null" json:"customer_name"`             // 客户名称（如：拉扎斯网络科技（上海）有限公司）
	CustomerShort     string    `gorm:"column:customer_short;type:varchar(50);not null" json:"customer_short"`            // 客户简称（如：淘宝闪购）
	AgentName         string    `gorm:"column:agent_name;type:varchar(200);not null" json:"agent_name"`                   // 代理名称（如：北京美数信息科技有限公司）
	AgentShort        string    `gorm:"column:agent_short;type:varchar(50);not null" json:"agent_short"`                  // 代理简称（如：美数）
	MediaPlatformName string    `gorm:"column:media_platform_name;type:varchar(100);not null" json:"media_platform_name"` // 媒体平台名称（如：巨量引擎）
	RedirectNum       int       `gorm:"column:redirect_num;not null;default:0" json:"redirect_num"`                       // 调起数
	PayNum            int       `gorm:"column:pay_num;not null;default:0" json:"pay_num"`                                 // 付费数
	CreateTime        time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`                             // 创建时间
	UpdateTime        time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`                             // 更新时间
}

// TableName 指定表名
func (ElmHcPerformanceReport) TableName() string {
	return "elm_hc_performance_report"
}

// GetAll 获取所有汇川饿了么数据报表
func GetAllElmHcPerformanceReports(db *gorm.DB) ([]ElmHcPerformanceReport, error) {
	var reports []ElmHcPerformanceReport
	err := db.Order("id DESC").Find(&reports).Error
	return reports, err
}

// GetByID 根据ID获取汇川饿了么数据报表
func GetElmHcPerformanceReportByID(db *gorm.DB, id int64) (*ElmHcPerformanceReport, error) {
	var report ElmHcPerformanceReport
	err := db.Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// Create 创建汇川饿了么数据报表
func CreateElmHcPerformanceReport(db *gorm.DB, report *ElmHcPerformanceReport) error {
	return db.Create(report).Error
}

// Update 更新汇川饿了么数据报表
func UpdateElmHcPerformanceReport(db *gorm.DB, report *ElmHcPerformanceReport) error {
	return db.Model(&ElmHcPerformanceReport{}).Where("id = ?", report.ID).Updates(map[string]interface{}{
		"customer_name":       report.CustomerName,
		"customer_short":      report.CustomerShort,
		"agent_name":          report.AgentName,
		"agent_short":         report.AgentShort,
		"media_platform_name": report.MediaPlatformName,
		"redirect_num":        report.RedirectNum,
		"pay_num":             report.PayNum,
	}).Error
}

// Delete 删除汇川饿了么数据报表
func DeleteElmHcPerformanceReport(db *gorm.DB, id int64) error {
	return db.Delete(&ElmHcPerformanceReport{}, id).Error
}
