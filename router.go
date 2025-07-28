package main

import (
	"Server/controller"
	"Server/database"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	filterRoutePrefix = [...]string{"/api/user/login", "/api/user/regist", "/api/user/registVerifyCode", "/download"}
)

// 路由初始化
func RouteInit(config map[string]interface{}) *gin.Engine {
	//设置服务器模式
	gin.SetMode(config["mode"].(string))
	//
	server := gin.Default()
	//设置过滤器
	setFilter(server)
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
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.Next()
	})
}

// 设置过滤器
func setFilter(server *gin.Engine) {
	server.Use(func(ctx *gin.Context) {
		//对指定路由不执行过滤操作
		route := ctx.Request.URL.Path
		if FilterDisable(route) {
			ctx.Next()
			return
		}
		//否则需要进行验证cookie
		cookie := ctx.Request.Header.Get("cookie")
		cookie, _ = url.QueryUnescape(cookie) //解码
		//如果cookie不合法直接过滤掉
		if !checkCookieLegal(cookie) {
			ctx.Abort()
			return
		}
		//
		ctx.Next()
	})
}

// 判断是否需要禁用过滤器
func FilterDisable(route string) bool {
	for _, s := range filterRoutePrefix {
		if strings.HasPrefix(route, s) {
			return true
		}
	}
	return false
}

// 判断cookie是否合法
func checkCookieLegal(cookie string) bool {
	return database.CheckCookieExist(cookie)
}

// 检查指定origin是否跨域访问服务器
func CheckOriginAllowed(origin string) bool {
	return true
}
