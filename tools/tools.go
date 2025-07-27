package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 消息回复的格式
type RespondMSG struct {
	Msg    string      `json:"msg"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

// 回复消息
func RespondMsg(ctx *gin.Context, code int, msg *RespondMSG) {
	ctx.JSON(code, msg)
}

// 回复code为200的消息
func RespondACK(ctx *gin.Context, msg *RespondMSG) {
	ctx.JSON(http.StatusOK, msg)
}

// 回复code为404的消息
func RespondNAK(ctx *gin.Context, msg *RespondMSG) {
	ctx.JSON(http.StatusNotFound, msg)
}
