package main

import (
	config "Server/Config"
	"Server/controller"

	"github.com/gin-gonic/gin"
)

func RouteInit() *gin.Engine {
	//设置服务器模式
	gin.SetMode(config.Configuration["server"].(map[string]interface{})["mode"].(string))
	//
	server := gin.Default()
	//
	server.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"msg": "ok"})
	})
	//配置路由
	user := server.Group("/user")
	{
		user.POST("/regist", controller.AddUser)
	}
	//
	return server
}
