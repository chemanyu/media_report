package model

import (
	"time"

	"gorm.io/gorm"
)

type CainiaoAdvertiser struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MediaAdvId string    `gorm:"column:media_adv_id;type:varchar(100);not null;uniqueIndex:uk_unique_record" json:"media_adv_id"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (CainiaoAdvertiser) TableName() string {
	return "cainiao_advertiser"
}

// GetAllAdvertiserIds 获取所有广告主ID列表
func GetAllAdvertiserIds(db *gorm.DB) ([]string, error) {
	var advertisers []CainiaoAdvertiser
	err := db.Select("media_adv_id").Find(&advertisers).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(advertisers))
	for _, adv := range advertisers {
		ids = append(ids, adv.MediaAdvId)
	}
	return ids, nil
}

// ListAdvertisers 获取所有广告主列表
func ListAdvertisers(db *gorm.DB) ([]CainiaoAdvertiser, error) {
	var list []CainiaoAdvertiser
	err := db.Order("id DESC").Find(&list).Error
	return list, err
}

// AddAdvertiser 添加广告主
func AddAdvertiser(db *gorm.DB, mediaAdvId string) error {
	adv := CainiaoAdvertiser{
		MediaAdvId: mediaAdvId,
	}
	return db.Create(&adv).Error
}

// DeleteAdvertiser 删除广告主
func DeleteAdvertiser(db *gorm.DB, id int64) error {
	return db.Delete(&CainiaoAdvertiser{}, id).Error
}
