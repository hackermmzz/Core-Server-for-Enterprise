package main

import (
	"Server/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteInit(config map[string]interface{}) *gin.Engine {
	//设置服务器模式
	gin.SetMode(config["mode"].(string))
	//
	server := gin.Default()
	//设置跨域
	setCors(server)
	//配置路由
	configRoute(server)
	//
	return server
}

// 配置路由
func configRoute(server *gin.Engine) {
	//配置api
	api := server.Group("/api")
	{
		//用户操作
		user := api.Group("/user")
		{
			user.POST("/regist", controller.Regist)
			user.POST("/login", controller.Login)
			user.POST("/registVerifyCode", controller.RequestRegistVerifyCode)
		}
	}
	//下载操作
	download := server.Group("/download")
	{
		download.GET("/*path", controller.DownLoad)
	}
}

// 设置跨域请求
func setCors(server *gin.Engine) {
	server.Use(func(ctx *gin.Context) {
		//获取请求头里面的域
		origin := ctx.Request.Header.Get("Origin")
		legal_origin := "https://www.adhn.asia"
		if CheckOriginAllowed(origin) {
			legal_origin = origin
		}
		//
		ctx.Header("Access-Control-Allow-Origin", legal_origin)
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Header.Get("Method") == "OPTIONS" {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.Next()
	})

}

// 检查指定origin是否跨域访问服务器
func CheckOriginAllowed(origin string) bool {
	return true
}
