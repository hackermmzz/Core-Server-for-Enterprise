package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteInit() *gin.Engine {
	
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	//配置路由
	user := server.Group("/user")
	{
		user.GET("/b/a", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "a"})
		})
		user.GET("/b", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "b"})
		})
	}
	//
	return server
}
