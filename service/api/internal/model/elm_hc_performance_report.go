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
	MediaAdvId        string    `gorm:"column:media_adv_id;type:varchar(100);not null" json:"media_adv_id"`               // 媒体账户ID（如：巨量引擎平台账户ID）
	MediaAdvName      string    `gorm:"column:media_adv_name;type:varchar(200);not null" json:"media_adv_name"`           // 媒体账户名称（如：巨量引擎平台账户名称）
	HuichuanAdvId     int64     `gorm:"column:huichuan_adv_id;not null" json:"huichuan_adv_id"`                           // 汇川账户ID
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

// GetByUniqueKey 根据唯一键查询（媒体平台名称+媒体账户ID+汇川账户ID）
func GetElmHcPerformanceReportByUniqueKey(db *gorm.DB, mediaPlatformName, mediaAdvId string, huichuanAdvId int64) (*ElmHcPerformanceReport, error) {
	var report ElmHcPerformanceReport
	err := db.Where("media_platform_name = ? AND media_adv_id = ? AND huichuan_adv_id = ?",
		mediaPlatformName, mediaAdvId, huichuanAdvId).First(&report).Error
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
		"media_adv_id":        report.MediaAdvId,
		"media_adv_name":      report.MediaAdvName,
		"huichuan_adv_id":     report.HuichuanAdvId,
		"redirect_num":        report.RedirectNum,
		"pay_num":             report.PayNum,
	}).Error
}

// Delete 删除汇川饿了么数据报表
func DeleteElmHcPerformanceReport(db *gorm.DB, id int64) error {
	return db.Delete(&ElmHcPerformanceReport{}, id).Error
}

// CreateOrUpdate 创建或更新（根据唯一键判断）
func CreateOrUpdateElmHcPerformanceReport(db *gorm.DB, report *ElmHcPerformanceReport) error {
	// 先尝试查询是否存在
	existing, err := GetElmHcPerformanceReportByUniqueKey(db, report.MediaPlatformName, report.MediaAdvId, report.HuichuanAdvId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在，创建新记录
			return CreateElmHcPerformanceReport(db, report)
		}
		return err
	}

	// 已存在，更新记录
	report.ID = existing.ID
	return UpdateElmHcPerformanceReport(db, report)
}

// BatchCreate 批量创建汇川饿了么数据报表
func BatchCreateElmHcPerformanceReports(db *gorm.DB, reports []ElmHcPerformanceReport) error {
	if len(reports) == 0 {
		return nil
	}
	return db.Create(&reports).Error
}

// BatchCreateOrUpdate 批量创建或更新
func BatchCreateOrUpdateElmHcPerformanceReports(db *gorm.DB, reports []ElmHcPerformanceReport) error {
	if len(reports) == 0 {
		return nil
	}

	// 使用事务处理
	return db.Transaction(func(tx *gorm.DB) error {
		for _, report := range reports {
			if err := CreateOrUpdateElmHcPerformanceReport(tx, &report); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByMediaPlatform 根据媒体平台名称查询
func GetElmHcPerformanceReportsByMediaPlatform(db *gorm.DB, mediaPlatformName string) ([]ElmHcPerformanceReport, error) {
	var reports []ElmHcPerformanceReport
	err := db.Where("media_platform_name = ?", mediaPlatformName).Order("id DESC").Find(&reports).Error
	return reports, err
}

// GetByHuichuanAdvId 根据汇川账户ID查询
func GetElmHcPerformanceReportsByHuichuanAdvId(db *gorm.DB, huichuanAdvId int64) ([]ElmHcPerformanceReport, error) {
	var reports []ElmHcPerformanceReport
	err := db.Where("huichuan_adv_id = ?", huichuanAdvId).Order("id DESC").Find(&reports).Error
	return reports, err
}

// GetByCustomerShort 根据客户简称查询
func GetElmHcPerformanceReportsByCustomerShort(db *gorm.DB, customerShort string) ([]ElmHcPerformanceReport, error) {
	var reports []ElmHcPerformanceReport
	err := db.Where("customer_short = ?", customerShort).Order("id DESC").Find(&reports).Error
	return reports, err
}

// GetByAgentShort 根据代理简称查询
func GetElmHcPerformanceReportsByAgentShort(db *gorm.DB, agentShort string) ([]ElmHcPerformanceReport, error) {
	var reports []ElmHcPerformanceReport
	err := db.Where("agent_short = ?", agentShort).Order("id DESC").Find(&reports).Error
	return reports, err
}
