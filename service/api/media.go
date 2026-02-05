package main

import (
	"flag"
	"fmt"

	"media_report/common/database"
	"media_report/service/api/internal/config"
	"media_report/service/api/internal/handler"
	"media_report/service/api/internal/script"
	"media_report/service/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/media-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 配置文件数据获取
	var config config.Config
	conf.MustLoad(*configFile, &config)

	// server ctx 获取
	server := rest.MustNewServer(config.RestConf)
	defer server.Stop()

	// 初始化数据库连接
	dbConfig := database.MySQLConfig{
		Host:            config.MySQL.Host,
		Port:            config.MySQL.Port,
		User:            config.MySQL.User,
		Password:        config.MySQL.Password,
		Database:        config.MySQL.Database,
		Charset:         config.MySQL.Charset,
		MaxIdleConns:    config.MySQL.MaxIdleConns,
		MaxOpenConns:    config.MySQL.MaxOpenConns,
		ConnMaxLifetime: config.MySQL.ConnMaxLifetime,
		ConnMaxIdleTime: config.MySQL.ConnMaxIdleTime,
		LogFile:         config.MySQL.LogFile,
		LogLevel:        config.MySQL.LogLevel,
	}
	db := database.MustNewMySQLConnection(dbConfig)

	// 获取上下文，注册接口
	ctx := svc.NewServiceContext(config, db)
	handler.RegisterHandlers(server, ctx)

	// 定时任务
	script.Cron(config, db)
	fmt.Printf("Starting server at %s:%d...\n", config.Host, config.Port)
	server.Start()
}
