package controller

import (
	service "Server/Service"
	"Server/database"
	"Server/tools"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	emailRe        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	emailMaxLen    = 50 //邮箱最大长度
	passwordRe     = regexp.MustCompile(`^[a-zA-Z0-9@#_]+$`)
	passwordMaxLen = 50 //密码最大长度
	nicknameMaxLen = 50 //昵称最大长度
)

// 验证注册
func Regist(ctx *gin.Context) {
	type UserInfo struct {
		Email      string `json:"email" binding:"required"`
		Password   string `json:"password" binding:"required"`
		NickName   string `json:"nickname" binding:"required"`
		VerifyCode string `json:"verifyCode" binding:"required"`
	}
	//解析数据
	var user UserInfo
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "传入的数据有误!"})
		return
	}
	//
	email := user.Email
	password := user.Password
	nickname := user.NickName
	verifyCode := user.VerifyCode
	//检测验证码是否正确
	exist := CheckRegistVerifyCodeLegal(email, verifyCode)
	if !exist {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "验证码错误!"})
		return
	}
	//检测用户是否存在
	_, exist = GetUser(email)
	if exist {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "用户已存在!"})
		return
	}
	//检测账号 密码 昵称的合法性
	if err := CheckEmailLegal(email); err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: err.Error()})
		return
	}
	if err := CheckPasswordLegal(password); err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: err.Error()})
		return
	}
	if !CheckNickNameLegal(nickname) {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "昵称不合法!"})
		return
	}
	//加密密码
	encrpy_passwd, err := EncryptPassword(password)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "后端加密密码出错!"})
		return
	}
	//添加入数据库
	var usr database.User
	usr.Email = email
	usr.Password = encrpy_passwd
	usr.NickName = user.NickName
	err = database.AddUser(usr)
	//
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "无法添加该用户!"})
		return
	}
	//
	tools.RespondACK(ctx, &tools.RespondMSG{Status: true, Msg: "添加用户成功!"})
}

// 验证登录
func Login(ctx *gin.Context) {
	type UserInfo struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	//解析数据
	var user UserInfo
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "传入的数据有误!"})
		return
	}
	//去数据库查找该用户
	record_user, exist := GetUser(user.Email)
	if !exist {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "用户不存在!"})
		return
	}
	//验证密码是否匹配
	if !CheckPasswordMatch(user.Password, record_user.Password) {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "用户名或者密码错误!"})
		return
	}
	//验证通过
	//设置set-cookie字段
	cookie, cookieValue := GenerateUserCookie(user.Email)
	ctx.SetSameSite(cookie.SameSite)
	ctx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	//将cookie写入数据库
	var cookieInfo database.Cookie
	cookieInfo.CookieValue = cookieValue
	cookieInfo.Email = user.Email
	err = database.AddUserCookie(cookieInfo)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "添加cookie失败!"})
		return
	}
	//
	//登陆成功
	tools.RespondACK(ctx, &tools.RespondMSG{Status: true, Msg: "登录成功!"})
}

// 获取验证码
func RequestRegistVerifyCode(ctx *gin.Context) {
	type Info struct {
		Email string `json:"email" binding:"required"`
	}
	//解析数据
	var user Info
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "传入的数据有误!"})
		return
	}
	//检查是距离上一次需求验证码小于一分钟
	if !checkRegistVerifyCodeIntervalLegal(user.Email) {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "间隔小于一分钟!"})
		return
	}
	//随机生成验证码
	code := GenerateRegistVerifyCode()
	//存入数据库
	verifyCode := &database.VerifyCode{Email: user.Email, Code: code, GenerateDate: nil, Usage: database.VerifyCode_Regist}
	err = database.AddVerifyCode(verifyCode)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "数据库异常!"})
		return
	}
	//放入验证码队列
	var msg service.VerifyCodeMsg
	msg.Email = user.Email
	msg.Code = code
	err = service.AddRegistVerifyCodeMsg(&msg)
	if err != nil {
		tools.RespondACK(ctx, &tools.RespondMSG{Status: false, Msg: "邮件服务器异常!"})
		return
	}
	//
	tools.RespondACK(ctx, &tools.RespondMSG{Status: true, Msg: "验证码发送成功!"})
}

// 获取指定邮箱的用户
func GetUser(email string) (database.User, bool) {
	user, err := database.GetUser(email)
	if err != nil {
		return user, false
	}
	return user, true
}

// 每一个存储入数据库的密码都是经过加密的
func EncryptPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// 检测账号的合法性

func CheckEmailLegal(email string) error {
	//长度不能超过emailMaxLen并且格式为A@B.C,只能为字母和数字
	if len(email) > emailMaxLen {
		return errors.New("邮箱长度不能超过" + strconv.Itoa(emailMaxLen))
	} else {
		match := passwordRe.MatchString(email)
		if !match {
			return errors.New("邮箱格式不正确")
		}
	}
	return nil
}

// 检测密码的合法性
func CheckPasswordLegal(password string) error {
	//长度不能超过emailMaxLen并且格式为A@B.C,只能为字母和数字
	if len(password) > passwordMaxLen {
		return errors.New("密码长度不能超过" + strconv.Itoa(passwordMaxLen))
	} else {
		match := emailRe.MatchString(password)
		if !match {
			return errors.New("密码格式不正确")
		}
	}
	return nil
}

// 检查昵称的合法性
func CheckNickNameLegal(nickname string) bool {
	return len(nickname) <= nicknameMaxLen
}

// 检测原密码和加密后的密码是否匹配
func CheckPasswordMatch(password_original string, password_encrypt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password_encrypt), []byte(password_original)) == nil
}

// 生成一个cookie
func GenerateUserCookie(email string) (*http.Cookie, string) {
	name := "token"
	value := email
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   3600 * 24, //一天
		Path:     "/",
		Secure:   true,                  // 仅 HTTPS
		HttpOnly: true,                  // 禁止 JavaScript 访问
		SameSite: http.SameSiteNoneMode, // 跨站限制
	}
	return cookie, name + "=" + value
}

// 检测注册验证码是否正确
func CheckRegistVerifyCodeLegal(email string, code string) bool {
	res, err := database.GetTargetVerifyCode(email, database.VerifyCode_Regist)
	if err != nil {
		return false
	}
	for _, r := range res {
		if r.Code == code {
			return true
		}
	}
	return false
}

// 检查两次验证码间隔是否大于一分钟
func checkRegistVerifyCodeIntervalLegal(email string) bool {
	res, err := database.GetEarliestVerifyCode(email, database.VerifyCode_Regist)
	//如果没有记录(默认没找到)
	if err != nil || res.GenerateDate == nil {
		return true
	}
	//如果大于1分钟
	if time.Time(*res.GenerateDate).Before(time.Now().Add(time.Duration(1 * time.Minute))) {
		return true
	}
	//
	return false
}

// 生成注册验证码
func GenerateRegistVerifyCode() string {
	t := time.Now().UnixNano() / int64(1e6)
	str := strconv.Itoa(int(t))
	length := len(str)
	return str[length-5 : length]
}
