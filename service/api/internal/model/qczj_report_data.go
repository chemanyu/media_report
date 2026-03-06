package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QczjReportData QCZJ ADN 日报数据
type QczjReportData struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReportDate int       `gorm:"column:report_date;uniqueIndex:uk_report_date;not null" json:"report_date"`
	View       int64     `gorm:"column:view;not null;default:0" json:"view"`
	Click      int64     `gorm:"column:click;not null;default:0" json:"click"`
	Action     int64     `gorm:"column:action;not null;default:0" json:"action"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (QczjReportData) TableName() string {
	return "qczj_report_data"
}

// ListQczjReportData 查询报表数据列表（按日期倒序）
func ListQczjReportData(db *gorm.DB) ([]*QczjReportData, int64, error) {
	var list []*QczjReportData
	var total int64

	db.Model(&QczjReportData{}).Count(&total)
	err := db.Order("report_date DESC").Find(&list).Error
	return list, total, err
}
func InsertOrUpdateQczjReportData(db *gorm.DB, data *QczjReportData) error {
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "report_date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"view":        data.View,
			"click":       data.Click,
			"action":      data.Action,
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(data).Error

	if err != nil {
		return fmt.Errorf("保存 qczj 数据失败: %w", err)
	}
	return nil
}

// ListQczjReportDataByDate 查询某天所有小时的数据，返回 map[hour]*QczjReportData
// date 格式为 "20260306"，report_date 存储格式为 YYYYMMDDHH（10位整数）
func ListQczjReportDataByDate(db *gorm.DB, date string) (map[string]*QczjReportData, error) {
	minDate, err := strconv.Atoi(date + "00")
	if err != nil {
		return nil, fmt.Errorf("日期格式错误: %w", err)
	}
	maxDate, err := strconv.Atoi(date + "23")
	if err != nil {
		return nil, fmt.Errorf("日期格式错误: %w", err)
	}

	var list []*QczjReportData
	if err := db.Where("report_date BETWEEN ? AND ?", minDate, maxDate).Find(&list).Error; err != nil {
		return nil, err
	}

	result := make(map[string]*QczjReportData, len(list))
	for _, item := range list {
		// report_date 最后两位是小时
		hourStr := fmt.Sprintf("%02d", item.ReportDate%100)
		result[hourStr] = item
	}
	return result, nil
}
