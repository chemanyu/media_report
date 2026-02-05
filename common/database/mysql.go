package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	ConnMaxLifetime int    // 连接最大生命周期（秒）
	ConnMaxIdleTime int    // 连接最大空闲时间（秒）
	LogFile         string // SQL日志文件路径
	LogLevel        string // SQL日志级别: silent, error, warn, info
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

	// 创建自定义的 SQL logger
	sqlLogger := createSQLLogger(config.LogFile, config.LogLevel)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
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

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func createSQLLogger(logFile string, logLevel string) logger.Interface {
	// 如果没有指定日志文件，使用默认的控制台日志
	if logFile == "" {
		return logger.Default.LogMode(parseLogLevel(logLevel))
	}

	// 确保日志目录存在
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("创建SQL日志目录失败: %v, 将使用控制台输出", err)
		return logger.Default.LogMode(parseLogLevel(logLevel))
	}

	// 打开或创建日志文件（追加模式）
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开SQL日志文件失败: %v, 将使用控制台输出", err)
		return logger.Default.LogMode(parseLogLevel(logLevel))
	}

	// 创建自定义logger，输出到文件
	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,  // 慢查询阈值
			LogLevel:                  parseLogLevel(logLevel), // 日志级别
			IgnoreRecordNotFoundError: true,                    // 忽略 ErrRecordNotFound 错误
			Colorful:                  false,                   // 文件日志不需要颜色
		},
	)

	return newLogger
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info // 默认info级别
	}
}

// MustNewMySQLConnection 创建 MySQL 数据库连接，失败时 panic
func MustNewMySQLConnection(config MySQLConfig) *gorm.DB {
	db, err := NewMySQLConnection(config)
	if err != nil {
		panic(err)
	}
	return db
}

// SetLogOutput 动态设置日志输出（用于在运行时切换日志输出位置）
func SetLogOutput(db *gorm.DB, writer io.Writer) {
	db.Logger = logger.New(
		log.New(writer, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
