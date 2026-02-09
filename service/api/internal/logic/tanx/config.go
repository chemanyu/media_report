package tanx

import (
	"sync"

	"media_report/service/api/internal/config"
)

var (
	tanxConfig *config.TanxConfig
	configLock sync.RWMutex
	smtpConfig *config.SMTPConfig
)

// InitTanxConfig 从配置文件初始化 Tanx 配置
func InitTanxConfig(cfg config.TanxConfig) {
	configLock.Lock()
	defer configLock.Unlock()

	tanxConfig = &config.TanxConfig{
		AdSlots:    cfg.AdSlots,
		Recipients: cfg.Recipients,
	}
	smtpConfig = &cfg.SMTP
}

// GetTanxConfig 获取配置
func GetTanxConfig() *config.TanxConfig {
	configLock.RLock()
	defer configLock.RUnlock()

	// 返回配置副本，避免外部修改
	config := &config.TanxConfig{
		AdSlots:    make([]string, len(tanxConfig.AdSlots)),
		Recipients: make([]string, len(tanxConfig.Recipients)),
	}
	copy(config.AdSlots, tanxConfig.AdSlots)
	copy(config.Recipients, tanxConfig.Recipients)
	return config
}

// GetSMTPConfig 获取SMTP配置
func GetSMTPConfig() *config.SMTPConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return smtpConfig
}
