package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// FzMediaAdvertiser 飞猪媒体账户表
type FzMediaAdvertiser struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Media        string    `gorm:"column:media" json:"media"`                         // 媒体简称，oppo, xiaomi, adn, honor
	MediaAdvId   string    `gorm:"column:media_adv_id" json:"media_adv_id"`           // 媒体账户ID
	MediaAdvName string    `gorm:"column:media_adv_name" json:"media_adv_name"`       // 媒体账户名称
	ClientID     string    `gorm:"column:client_id" json:"client_id"`                 // OAuth2 Client ID（honor专用）
	ClientSecret string    `gorm:"column:client_secret" json:"client_secret"`         // OAuth2 Client Secret（honor专用）
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (FzMediaAdvertiser) TableName() string {
	return "fz_media_advertiser"
}

// FzMediaAdvertiserModel 飞猪媒体账户数据访问接口
type FzMediaAdvertiserModel interface {
	// FindByMedia 根据媒体类型查询所有账户
	FindByMedia(media string) ([]*FzMediaAdvertiser, error)
	// FindAll 查询所有账户
	FindAll() ([]*FzMediaAdvertiser, error)
	// FindOne 根据ID查询单个账户
	FindOne(id int64) (*FzMediaAdvertiser, error)
	// FindByMediaAdvId 根据媒体账户ID查询
	FindByMediaAdvId(mediaAdvId string) (*FzMediaAdvertiser, error)
	// Insert 插入新账户
	Insert(advertiser *FzMediaAdvertiser) error
	// Update 更新账户信息
	Update(advertiser *FzMediaAdvertiser) error
	// Delete 删除账户
	Delete(id int64) error
}

type defaultFzMediaAdvertiserModel struct {
	db *gorm.DB
}

// NewFzMediaAdvertiserModel 创建媒体账户模型实例
func NewFzMediaAdvertiserModel(db *gorm.DB) FzMediaAdvertiserModel {
	return &defaultFzMediaAdvertiserModel{
		db: db,
	}
}

// FindByMedia 根据媒体类型查询所有账户
func (m *defaultFzMediaAdvertiserModel) FindByMedia(media string) ([]*FzMediaAdvertiser, error) {
	var advertisers []*FzMediaAdvertiser
	err := m.db.Where("media = ?", media).Find(&advertisers).Error
	if err != nil {
		return nil, fmt.Errorf("查询媒体账户失败: %w", err)
	}
	return advertisers, nil
}

// FindOne 根据ID查询单个账户
func (m *defaultFzMediaAdvertiserModel) FindOne(id int64) (*FzMediaAdvertiser, error) {
	var adv FzMediaAdvertiser
	err := m.db.Where("id = ?", id).First(&adv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("账户不存在: id=%d", id)
		}
		return nil, fmt.Errorf("查询账户失败: %w", err)
	}
	return &adv, nil
}

// FindByMediaAdvId 根据媒体账户ID查询
func (m *defaultFzMediaAdvertiserModel) FindByMediaAdvId(mediaAdvId string) (*FzMediaAdvertiser, error) {
	var adv FzMediaAdvertiser
	err := m.db.Where("media_adv_id = ?", mediaAdvId).First(&adv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("账户不存在: media_adv_id=%s", mediaAdvId)
		}
		return nil, fmt.Errorf("查询账户失败: %w", err)
	}
	return &adv, nil
}

// Insert 插入新账户
func (m *defaultFzMediaAdvertiserModel) Insert(advertiser *FzMediaAdvertiser) error {
	return m.db.Create(advertiser).Error
}

// FindAll 查询所有账户
func (m *defaultFzMediaAdvertiserModel) FindAll() ([]*FzMediaAdvertiser, error) {
	var advertisers []*FzMediaAdvertiser
	err := m.db.Order("id DESC").Find(&advertisers).Error
	if err != nil {
		return nil, fmt.Errorf("查询所有账户失败: %w", err)
	}
	return advertisers, nil
}

// Update 更新账户信息
func (m *defaultFzMediaAdvertiserModel) Update(advertiser *FzMediaAdvertiser) error {
	result := m.db.Model(&FzMediaAdvertiser{}).Where("id = ?", advertiser.Id).Updates(advertiser)
	if result.Error != nil {
		return fmt.Errorf("更新账户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("账户不存在: id=%d", advertiser.Id)
	}
	return nil
}

// Delete 删除账户
func (m *defaultFzMediaAdvertiserModel) Delete(id int64) error {
	result := m.db.Where("id = ?", id).Delete(&FzMediaAdvertiser{})
	if result.Error != nil {
		return fmt.Errorf("删除账户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("账户不存在: id=%d", id)
	}
	return nil
}
