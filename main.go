package main

import (
	config "Server/Config"
	service "Server/Service"
	"Server/database"
	"Server/logger"
	"math/rand"
	"time"
)

func main() {
	//初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	//加载配置
	config.ConfigInit()
	//初始化log日志
	logger.LogInit()
	//初始化数据库
	database.DabaseInit()
	//初始化所有的服务
	service.ServiceInit()
	//初始化服务器
	ServerInit()

}
