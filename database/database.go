package database

import (
	config "Server/Config"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func DabaseInit() {
	var err error
	//获取数据库账户和密码
	dbJS := config.Configuration["database"].(map[string]interface{})
	user := dbJS["user"]
	password := dbJS["password"]
	database := dbJS["database"]
	address := dbJS["address"]
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?charset=utf8mb4&&parseTime=true&loc=UTC", user, password, address, database)
	//重定向gorm的输出
	gormLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags), // 输出到日志文件
		logger.Config{
			SlowThreshold:             time.Second,  // 慢查询阈值
			LogLevel:                  logger.Error, // 日志级别（Info/ Warn/ Error/ Silent）
			IgnoreRecordNotFoundError: true,         // 是否忽略记录未找到错误
			Colorful:                  false,        // 非控制台输出，关闭彩色日志
		},
	)
	//
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	fmt.Println("数据库连接成功!")
}
