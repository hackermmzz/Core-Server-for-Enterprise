package main

import (
	config "Server/Config"
	"fmt"
)

//

// 初始化服务器
func ServerInit() {
	serverConfig := config.Configuration["server"].(map[string]interface{})
	//初始化路由
	server := RouteInit(serverConfig)
	//监听
	fmt.Println("服务器初始化成功!")
	err := server.RunTLS(":"+serverConfig["port"].(string), "ssl/cert.pem", "ssl/key.pem")
	if err != nil {
		panic("服务器初始化错误:" + err.Error())
	}
}
