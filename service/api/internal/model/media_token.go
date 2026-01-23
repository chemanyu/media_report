package model

import (
	"time"

	"gorm.io/gorm"
)

// MediaToken 媒体 token 表
type MediaToken struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Media        string `gorm:"column:media;type:varchar(64);default:''" json:"media"`                       // 媒体名称
	Token        string `gorm:"column:token;type:varchar(255);default:''" json:"token"`                      // 媒体token
	RefreshToken string `gorm:"column:refresh_token;type:varchar(62);default:''" json:"refresh_token"`       // 媒体刷新token
	AgentID      string `gorm:"column:agent_id;type:varchar(64);default:''" json:"agent_id"`                 // 代理商id
	AdvertiserID string `gorm:"column:advertiser_id;type:varchar(64);default:''" json:"advertiser_id"`       // 账户ID
	DelFlag      int8   `gorm:"column:del_flag;type:tinyint(1);default:0" json:"del_flag"`                   // 删除(0:正常;1:删除)
	CreateTime   int    `gorm:"column:create_time;type:int(11);default:0;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime   int    `gorm:"column:update_time;type:int(11);default:0;autoUpdateTime" json:"update_time"` // 修改时间
}

// TableName 指定表名
func (MediaToken) TableName() string {
	return "media_token"
}

// BeforeCreate GORM 钩子：创建前设置时间戳
func (t *MediaToken) BeforeCreate(tx *gorm.DB) error {
	now := int(time.Now().Unix())
	if t.CreateTime == 0 {
		t.CreateTime = now
	}
	if t.UpdateTime == 0 {
		t.UpdateTime = now
	}
	return nil
}

// BeforeUpdate GORM 钩子：更新前设置时间戳
func (t *MediaToken) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = int(time.Now().Unix())
	return nil
}

// GetByMedia 通过媒体名称获取 token 信息
func GetByMedia(db *gorm.DB, media string) (*MediaToken, error) {
	var token MediaToken
	err := db.Where("media = ? AND del_flag = ?", media, 0).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokensByMedia 通过媒体名称获取 token 和 refresh_token
func GetTokensByMedia(db *gorm.DB, media string) (token string, refreshToken string, err error) {
	mediaToken, err := GetByMedia(db, media)
	if err != nil {
		return "", "", err
	}
	return mediaToken.Token, mediaToken.RefreshToken, nil
}
