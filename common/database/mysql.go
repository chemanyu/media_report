package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	Charset         string
	MaxIdleConns    int // 最大空闲连接数
	MaxOpenConns    int // 最大打开连接数
	ConnMaxLifetime int // 连接最大生命周期（秒）
	ConnMaxIdleTime int // 连接最大空闲时间（秒）
}

// NewMySQLConnection 创建 MySQL 数据库连接
func NewMySQLConnection(config MySQLConfig) (*gorm.DB, error) {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层的 sql.DB 对象，配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	maxIdleConns := config.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10 // 默认值
	}

	maxOpenConns := config.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 50 // 默认值
	}

	connMaxLifetime := time.Duration(config.ConnMaxLifetime) * time.Second
	if config.ConnMaxLifetime <= 0 {
		connMaxLifetime = time.Hour // 默认1小时
	}

	connMaxIdleTime := time.Duration(config.ConnMaxIdleTime) * time.Second
	if config.ConnMaxIdleTime <= 0 {
		connMaxIdleTime = 10 * time.Minute // 默认10分钟
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)       // 最大空闲连接数
	sqlDB.SetMaxOpenConns(maxOpenConns)       // 最大打开连接数
	sqlDB.SetConnMaxLifetime(connMaxLifetime) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime) // 连接最大空闲时间

	return db, nil
}

// MustNewMySQLConnection 创建 MySQL 数据库连接，失败时 panic
func MustNewMySQLConnection(config MySQLConfig) *gorm.DB {
	db, err := NewMySQLConnection(config)
	if err != nil {
		panic(err)
	}
	return db
}
