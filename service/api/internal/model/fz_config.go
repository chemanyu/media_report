package model

import (
	"time"

	"gorm.io/gorm"
)

// FzConfig 飞猪现金消耗系数配置
type FzConfig struct {
	Id          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Coefficient float64   `gorm:"column:coefficient;type:decimal(6,4);not null;default:1.7000" json:"coefficient"` // 系数
	BaseNum     float64   `gorm:"column:base_num;type:decimal(6,4);not null;default:0.8500" json:"base_num"`        // 基数
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (FzConfig) TableName() string {
	return "fz_config"
}

// GetFzConfig 获取配置（ID=1）
func GetFzConfig(db *gorm.DB) (*FzConfig, error) {
	var config FzConfig
	err := db.Where("id = ?", 1).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateFzConfig 更新配置（ID=1）
func UpdateFzConfig(db *gorm.DB, coefficient, baseNum float64) error {
	return db.Model(&FzConfig{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{
			"coefficient": coefficient,
			"base_num":    baseNum,
		}).Error
}
