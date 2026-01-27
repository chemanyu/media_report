package model

import (
	"time"

	"gorm.io/gorm"
)

// Rebate 返点配置表
type Rebate struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Subject     string    `gorm:"column:subject;type:varchar(50);not null" json:"subject"`                 // 主体：新杰、晴川、美数、魔米、归流
	Port        string    `gorm:"column:port;type:varchar(50);not null" json:"port"`                       // 端口：优居、至也、各界
	RebateRate  float64   `gorm:"column:rebate_rate;type:decimal(5,4);not null" json:"rebate_rate"`        // 返点率：如0.025表示2.5%
	SubjectType int8      `gorm:"column:subject_type;type:tinyint;not null;default:1" json:"subject_type"` // 主体类型：1-京东主体, 2-三方主体
	Remark      string    `gorm:"column:remark;type:varchar(255);default:''" json:"remark"`                // 备注
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (Rebate) TableName() string {
	return "rebate"
}

// GetAll 获取所有返点配置
func GetAllRebates(db *gorm.DB) ([]Rebate, error) {
	var rebates []Rebate
	err := db.Order("id DESC").Find(&rebates).Error
	return rebates, err
}

// GetByID 根据ID获取返点配置
func GetRebateByID(db *gorm.DB, id uint) (*Rebate, error) {
	var rebate Rebate
	err := db.Where("id = ?", id).First(&rebate).Error
	if err != nil {
		return nil, err
	}
	return &rebate, nil
}

// Create 创建返点配置
func CreateRebate(db *gorm.DB, rebate *Rebate) error {
	return db.Create(rebate).Error
}

// Update 更新返点配置
func UpdateRebate(db *gorm.DB, rebate *Rebate) error {
	return db.Model(&Rebate{}).Where("id = ?", rebate.ID).Updates(map[string]interface{}{
		"subject":      rebate.Subject,
		"port":         rebate.Port,
		"rebate_rate":  rebate.RebateRate,
		"subject_type": rebate.SubjectType,
		"remark":       rebate.Remark,
	}).Error
}

// Delete 删除返点配置
func DeleteRebate(db *gorm.DB, id uint) error {
	return db.Delete(&Rebate{}, id).Error
}
