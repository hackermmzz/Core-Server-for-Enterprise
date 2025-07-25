package main

import "fmt"

//

// 初始化服务器
func ServerInit() {
	server := RouteInit()
	err := server.RunTLS(":443", "ssl/cert.pem", "ssl/key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	Log("服务器初始化成功!")

}
