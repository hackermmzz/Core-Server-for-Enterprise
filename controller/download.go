package controller

import (
	config "Server/Config"
	"Server/tools"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var wordMap = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

func DownLoad(ctx *gin.Context) {
	sourceDir := config.Configuration["Source"].(map[string]interface{})["directory"].(string)
	path := sourceDir + "/" + ctx.Param("path")
	fmt.Println(path)
	//检查path，不允许出现../
	if !checkPathLegal(path) {
		tools.RespondMsg(ctx, http.StatusForbidden, &tools.RespondMSG{Status: false, Msg: "路径非法!"})
		return
	}
	//检查文件是否存在
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		tools.RespondNAK(ctx, &tools.RespondMSG{Status: false, Msg: "文件不存在!"})
		return
	}
	//获取文件名
	pathSlice := strings.Split(path, "/")
	filename := pathSlice[len(pathSlice)-1]
	//
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename) // 用path作为文件名
	ctx.File(path)                                                      // 发送文件内容
}

// 检查路径是否是恶意得
func checkPathLegal(path string) bool {
	return !strings.Contains(path, "../")
}

// 随机生成一个文件名
func generateRandomFileName(length int, sufix string) string {
	prelen := length - len(sufix)
	var ret []byte
	//
	for i := 0; i < prelen; i++ {
		idx := rand.Intn(len(wordMap))
		ch := wordMap[idx]
		ret = append(ret, ch)
	}
	//
	return string(ret) + sufix

}
