package service

import (
	config "Server/Config"
	"fmt"
)

// 初始化所有的服务
func ServiceInit() {
	js := config.Configuration["service"].(map[string]interface{})
	//cookie过期检测服务
	go RemoveCookieExpireService(js["cookieExpire"].(map[string]interface{}))
	//验证码发送队列
	go VerifyCodeSendService(js["verifyCodeSend"].(map[string]interface{}))
	//
	fmt.Println("所有服务启动成功!")
}
