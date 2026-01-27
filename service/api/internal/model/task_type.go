package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskType 任务类型表
type TaskType struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"column:name;type:varchar(50);not null" json:"name"`                           // 任务名称，如：app、首购
	Code            string    `gorm:"column:code;type:varchar(20);not null" json:"code"`                           // 任务编码
	SettlementPrice float64   `gorm:"column:settlement_price;type:decimal(10,2);not null" json:"settlement_price"` // 结算单价
	Media           string    `gorm:"column:media;type:varchar(50);not null" json:"media"`                         // 媒体平台，如：巨量
	Status          int8      `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`                 // 状态：1-启用, 0-停用
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (TaskType) TableName() string {
	return "task_type"
}

// GetAll 获取所有任务类型
func GetAllTaskTypes(db *gorm.DB) ([]TaskType, error) {
	var taskTypes []TaskType
	err := db.Order("id DESC").Find(&taskTypes).Error
	return taskTypes, err
}

// GetByID 根据ID获取任务类型
func GetTaskTypeByID(db *gorm.DB, id uint) (*TaskType, error) {
	var taskType TaskType
	err := db.Where("id = ?", id).First(&taskType).Error
	if err != nil {
		return nil, err
	}
	return &taskType, nil
}

// Create 创建任务类型
func CreateTaskType(db *gorm.DB, taskType *TaskType) error {
	return db.Create(taskType).Error
}

// Update 更新任务类型
func UpdateTaskType(db *gorm.DB, taskType *TaskType) error {
	return db.Model(&TaskType{}).Where("id = ?", taskType.ID).Updates(map[string]interface{}{
		"name":             taskType.Name,
		"code":             taskType.Code,
		"settlement_price": taskType.SettlementPrice,
		"media":            taskType.Media,
		"status":           taskType.Status,
	}).Error
}

// Delete 删除任务类型
func DeleteTaskType(db *gorm.DB, id uint) error {
	return db.Delete(&TaskType{}, id).Error
}

// LoadTaskTypeConfigMap 加载任务类型配置为Map
// 返回 map[任务代码] = 结算单价
// 只加载启用状态的任务类型
func LoadTaskTypeConfigMap(db *gorm.DB) (map[string]float64, error) {
	var taskTypes []TaskType
	result := make(map[string]float64)

	if err := db.Where("status = ?", 1).Find(&taskTypes).Error; err != nil {
		return result, err
	}

	for _, task := range taskTypes {
		result[task.Name] = task.SettlementPrice
	}

	return result, nil
}

// GetByMediaAndCode 根据媒体和编码获取任务类型
func GetTaskTypeByMediaAndCode(db *gorm.DB, media, code string) (*TaskType, error) {
	var taskType TaskType
	err := db.Where("media = ? AND code = ?", media, code).First(&taskType).Error
	if err != nil {
		return nil, err
	}
	return &taskType, nil
}
