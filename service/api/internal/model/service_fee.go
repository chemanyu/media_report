package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ServiceFee 服务费配置表
type ServiceFee struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceProvider string    `gorm:"column:service_provider;type:varchar(50);not null;uniqueIndex:uniq_provider" json:"service_provider"` // 服务商名称：通途、蚁行、创效、凯旋、星河、云谷、美数
	FeeRate         float64   `gorm:"column:fee_rate;type:decimal(5,3);not null" json:"fee_rate"`                                          // 服务费率：如0.04表示4%
	Remark          string    `gorm:"column:remark;type:varchar(255);default:''" json:"remark"`                                            // 备注
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (ServiceFee) TableName() string {
	return "service_fee"
}

// GetAll 获取所有服务费配置
func GetAllServiceFees(db *gorm.DB) ([]ServiceFee, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库连接为空")
	}

	var fees []ServiceFee
	err := db.Order("id DESC").Find(&fees).Error
	if err != nil {
		return nil, fmt.Errorf("查询服务费配置失败: %w", err)
	}

	return fees, nil
}

// GetByID 根据ID获取服务费配置
func GetServiceFeeByID(db *gorm.DB, id uint) (*ServiceFee, error) {
	var fee ServiceFee
	err := db.Where("id = ?", id).First(&fee).Error
	if err != nil {
		return nil, err
	}
	return &fee, nil
}

// Create 创建服务费配置
func CreateServiceFee(db *gorm.DB, fee *ServiceFee) error {
	return db.Create(fee).Error
}

// Update 更新服务费配置
func UpdateServiceFee(db *gorm.DB, fee *ServiceFee) error {
	return db.Model(&ServiceFee{}).Where("id = ?", fee.ID).Updates(map[string]interface{}{
		"service_provider": fee.ServiceProvider,
		"fee_rate":         fee.FeeRate,
		"remark":           fee.Remark,
	}).Error
}

// Delete 删除服务费配置
func DeleteServiceFee(db *gorm.DB, id uint) error {
	return db.Delete(&ServiceFee{}, id).Error
}

// LoadServiceFeeConfigMap 加载服务费配置为Map
// 返回 map[服务商] = 服务费率
func LoadServiceFeeConfigMap(db *gorm.DB) (map[string]float64, error) {
	var serviceFees []ServiceFee
	result := make(map[string]float64)

	if err := db.Find(&serviceFees).Error; err != nil {
		return result, err
	}

	for _, fee := range serviceFees {
		result[fee.ServiceProvider] = fee.FeeRate
	}

	return result, nil
}
