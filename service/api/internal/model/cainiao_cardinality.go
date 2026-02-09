package model

import (
	"time"

	"gorm.io/gorm"
)

// CainiaoCardinality 菜鸟基数配置
type CainiaoCardinality struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Cardinality float64   `gorm:"column:cardinality;type:decimal(3,1);not null" json:"cardinality"` // 基数值，如1.4、4.1等
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (CainiaoCardinality) TableName() string {
	return "cainiao_cardinality"
}

// GetCardinality 获取基数配置（ID=1）
func GetCardinality(db *gorm.DB) (*CainiaoCardinality, error) {
	var config CainiaoCardinality
	err := db.Where("id = ?", 1).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateCardinality 更新基数配置（ID=1）
func UpdateCardinality(db *gorm.DB, cardinality float64) error {
	return db.Model(&CainiaoCardinality{}).
		Where("id = ?", 1).
		Update("cardinality", cardinality).Error
}
