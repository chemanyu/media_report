package model

import (
	"time"

	"gorm.io/gorm"
)

// TanxMonitor 淘宝联盟监控数据模型
type TanxMonitor struct {
	ID            int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Ds            string  `gorm:"column:ds;type:varchar(20);not null" json:"ds"`                             // 日期
	Pid           string  `gorm:"column:pid;type:varchar(50);not null" json:"pid"`                           // 广告位ID
	AdzoneName    string  `gorm:"column:adzone_name;type:varchar(200);default:''" json:"adzone_name"`        // 广告位名称
	Qingqiupv     int64   `gorm:"column:qingqiupv;default:0" json:"qingqiupv"`                               // tanx有效请求
	ActiveRatioDf string  `gorm:"column:active_ratio_df;type:varchar(50);default:''" json:"active_ratio_df"` // 东风手淘换端率-同步点击
	TanxEffectPv  int64   `gorm:"column:tanx_effect_pv;default:0" json:"tanx_effect_pv"`                     // TANX曝光数
	TanxClk       int64   `gorm:"column:tanx_clk;default:0" json:"tanx_clk"`                                 // TANX点击数
	DongfengEf    float64 `gorm:"column:dongfeng_ef;type:decimal(10,2);default:0" json:"dongfeng_ef"`        // TANX预估收益
	CreateTime    int64   `gorm:"column:create_time;autoCreateTime" json:"create_time"`                      // 创建时间
	//UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`                      // 更新时间
}

// TableName 指定表名
func (TanxMonitor) TableName() string {
	return "tanx_monitor"
}

// GetAll 获取所有淘宝联盟监控数据
func GetAllTanxMonitors(db *gorm.DB) ([]TanxMonitor, error) {
	var monitors []TanxMonitor
	err := db.Order("ds DESC, id DESC").Find(&monitors).Error
	return monitors, err
}

// GetByID 根据ID获取淘宝联盟监控数据
func GetTanxMonitorByID(db *gorm.DB, id int64) (*TanxMonitor, error) {
	var monitor TanxMonitor
	err := db.Where("id = ?", id).First(&monitor).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// GetByDsAndPid 根据日期和广告位ID获取监控数据
func GetTanxMonitorByDsAndPid(db *gorm.DB, ds, pid string) (*TanxMonitor, error) {
	var monitor TanxMonitor
	err := db.Where("ds = ? AND pid = ?", ds, pid).First(&monitor).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// Create 创建淘宝联盟监控数据
func CreateTanxMonitor(db *gorm.DB, monitor *TanxMonitor) error {
	return db.Create(monitor).Error
}

// Update 更新淘宝联盟监控数据
func UpdateTanxMonitor(db *gorm.DB, monitor *TanxMonitor) error {
	return db.Save(monitor).Error
}

// Upsert 创建或更新淘宝联盟监控数据（根据 ds 和 pid）
func UpsertTanxMonitor(db *gorm.DB, monitor *TanxMonitor) error {
	var existing TanxMonitor
	err := db.Where("ds = ? AND pid = ?", monitor.Ds, monitor.Pid).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		return db.Create(monitor).Error
	} else if err != nil {
		// 查询错误
		return err
	}

	// 存在，更新记录
	monitor.ID = existing.ID
	monitor.CreateTime = time.Now().Unix()
	return db.Save(monitor).Error
}

// GetByDateRange 根据日期范围获取监控数据
func GetTanxMonitorsByDateRange(db *gorm.DB, startDate, endDate string) ([]TanxMonitor, error) {
	var monitors []TanxMonitor
	err := db.Where("ds >= ? AND ds <= ?", startDate, endDate).Order("ds DESC, pid ASC").Find(&monitors).Error
	return monitors, err
}

// GetByPid 根据广告位ID获取所有监控数据
func GetTanxMonitorsByPid(db *gorm.DB, pid string) ([]TanxMonitor, error) {
	var monitors []TanxMonitor
	err := db.Where("pid = ?", pid).Order("ds DESC").Find(&monitors).Error
	return monitors, err
}
