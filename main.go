package main

import (
	config "Server/Config"
	"Server/database"
)

func main() {
	//加载配置
	config.ConfigInit()
	//初始化数据库
	database.DabaseInit()
	//初始化服务器
	ServerInit()

}
