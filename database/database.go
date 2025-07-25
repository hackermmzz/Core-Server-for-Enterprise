package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func DabaseInit() {
	var err error
	dsn := "root:01190650asd@tcp(127.0.0.1:3306)/SF?charset=utf8mb4"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

}
