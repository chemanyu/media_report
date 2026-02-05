package model

import (
	"time"

	"gorm.io/gorm"
)

// ElmHcMediaReport 汇川饿了么账户报表
type ElmHcMediaReport struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PerformanceID int64     `gorm:"column:performance_id;not null" json:"performance_id"`                   // 客户表id
	MediaAdvId    string    `gorm:"column:media_adv_id;type:varchar(100);not null" json:"media_adv_id"`     // 媒体账户ID
	MediaAdvName  string    `gorm:"column:media_adv_name;type:varchar(200);not null" json:"media_adv_name"` // 媒体账户名称
	HuichuanAdvId int64     `gorm:"column:huichuan_adv_id;not null" json:"huichuan_adv_id"`                 // 汇川账户ID
	RedirectNum   int       `gorm:"column:redirect_num;not null;default:0" json:"redirect_num"`             // 调起数
	PayNum        int       `gorm:"column:pay_num;not null;default:0" json:"pay_num"`                       // 付费数
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`                   // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`                   // 更新时间
}

// TableName 指定表名
func (ElmHcMediaReport) TableName() string {
	return "elm_hc_media_report"
}

// GetAll 获取所有汇川饿了么账户报表
func GetAllElmHcMediaReports(db *gorm.DB) ([]ElmHcMediaReport, error) {
	var reports []ElmHcMediaReport
	err := db.Order("id DESC").Find(&reports).Error
	return reports, err
}

// GetByID 根据ID获取汇川饿了么账户报表
func GetElmHcMediaReportByID(db *gorm.DB, id int64) (*ElmHcMediaReport, error) {
	var report ElmHcMediaReport
	err := db.Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// GetByUniqueKey 根据唯一键查询（媒体账户ID+汇川账户ID）
func GetElmHcMediaReportByUniqueKey(db *gorm.DB, mediaAdvId string, huichuanAdvId int64) (*ElmHcMediaReport, error) {
	var report ElmHcMediaReport
	err := db.Where("media_adv_id = ? AND huichuan_adv_id = ?", mediaAdvId, huichuanAdvId).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// Create 创建汇川饿了么账户报表
func CreateElmHcMediaReport(db *gorm.DB, report *ElmHcMediaReport) error {
	return db.Create(report).Error
}

// Update 更新汇川饿了么账户报表
func UpdateElmHcMediaReport(db *gorm.DB, report *ElmHcMediaReport) error {
	return db.Model(&ElmHcMediaReport{}).Where("id = ?", report.ID).Updates(map[string]interface{}{
		"performance_id":  report.PerformanceID,
		"media_adv_id":    report.MediaAdvId,
		"media_adv_name":  report.MediaAdvName,
		"huichuan_adv_id": report.HuichuanAdvId,
		"redirect_num":    report.RedirectNum,
		"pay_num":         report.PayNum,
	}).Error
}

// Delete 删除汇川饿了么账户报表
func DeleteElmHcMediaReport(db *gorm.DB, id int64) error {
	return db.Delete(&ElmHcMediaReport{}, id).Error
}

// GetByPerformanceId 根据客户表ID获取汇川饿了么账户报表
func GetElmHcMediaReportsByPerformanceId(db *gorm.DB, performanceId int) ([]ElmHcMediaReport, error) {
	var reports []ElmHcMediaReport
	err := db.Where("performance_id = ?", performanceId).Order("id DESC").Find(&reports).Error
	return reports, err
}
