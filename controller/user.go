package controller

import (
	"Server/database"
	"Server/tools"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 每一个存储入数据库的密码都是经过加密的
func EncryptPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// 检测账号的合法性
func CheckEmailLegal(email string) bool {
	return true
}

// 检测密码的合法性
func CheckPasswordLegal(password string) bool {
	return true
}

// 增加用户
func AddUser(ctx *gin.Context) {
	type UserInfo struct {
		email    string
		password string
	}
	//解析数据
	var user UserInfo
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "传入的数据有误!"})
		return
	}
	//
	email := user.email
	password := user.password
	//检测用户是否存在
	_, exist := GetUser(email)
	if exist {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "用户已存在!"})
		return
	}
	//检测账号和密码的合法性
	if !CheckEmailLegal(email) {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "邮箱格式不合法!"})
		return
	}
	if !CheckPasswordLegal(password) {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "密码格式不合法!"})
		return
	}
	//加密密码
	encrpy_passwd, err := EncryptPassword(password)
	if err != nil {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "后端加密密码出错!"})
		return
	}
	//添加入数据库
	err = database.AddUser(email, encrpy_passwd)
	//
	if err != nil {
		tools.RespondACK(ctx, tools.RespondMSG{Status: false, Msg: "无法添加该用户!"})
		return
	}
	//
	tools.RespondACK(ctx, tools.RespondMSG{Status: true, Msg: "添加用户成功!"})
}

// 获取指定邮箱的用户
func GetUser(email string) (database.User, bool) {
	user, err := database.GetUser(email)
	if err != nil {
		return user, false
	}
	return user, true
}
