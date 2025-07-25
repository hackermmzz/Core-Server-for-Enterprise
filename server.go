package main

import (
	"fmt"
)

//

// 初始化服务器
func ServerInit() {
	server := RouteInit()
	fmt.Println("服务器初始化成功!")
	err := server.RunTLS(":443", "ssl/cert.pem", "ssl/key.pem")
	if err != nil {
		panic("服务器初始化错误:" + err.Error())
	}
}
