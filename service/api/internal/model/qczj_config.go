package model

import (
	"time"

	"gorm.io/gorm"
)

// QczjConfig QCZJ ADN 数据配置
type QczjConfig struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TotalUv    int64     `gorm:"column:total_uv;not null;default:110000" json:"total_uv"`
	Ratio      float64   `gorm:"column:ratio;type:decimal(5,4);not null;default:0.4000" json:"ratio"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (QczjConfig) TableName() string {
	return "qczj_config"
}

// GetQczjConfig 获取配置（ID=1）
func GetQczjConfig(db *gorm.DB) (*QczjConfig, error) {
	var config QczjConfig
	err := db.Where("id = ?", 1).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateQczjConfig 更新配置（ID=1）
func UpdateQczjConfig(db *gorm.DB, totalUv int64, ratio float64) error {
	return db.Model(&QczjConfig{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{
			"total_uv": totalUv,
			"ratio":    ratio,
		}).Error
}
