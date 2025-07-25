package database

import (
	config "Server/Config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?charset=utf8mb4", user, password, address, database)
	//
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	fmt.Println("数据库连接成功!")
}
