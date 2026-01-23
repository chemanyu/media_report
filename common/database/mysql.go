package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Charset  string
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
	sqlDB.SetMaxIdleConns(10)      // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)     // 最大打开连接数
	sqlDB.SetConnMaxLifetime(3600) // 连接最大生命周期（秒）

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
