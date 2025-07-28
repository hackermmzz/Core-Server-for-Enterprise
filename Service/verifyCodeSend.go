package service

import (
	"Server/logger"
	"crypto/tls"
	"fmt"

	"github.com/Workiva/go-datastructures/queue"
	"gopkg.in/gomail.v2"
)

type VerifyCodeMsg struct {
	Email string
	Code  string
}

var verifyCodeQueue *queue.Queue
var email = ""
var name = ""
var sender *gomail.Dialer

func VerifyCodeSendService(config map[string]interface{}) {
	verifyCodeQueue = queue.New(1024)
	email = config["email"].(string)
	password := config["password"].(string)
	name = config["name"].(string)
	host := config["host"].(string)
	port := int(config["port"].(float64))
	sender = gomail.NewDialer(host, port, email, password)
	sender.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	//开始发送
	for {
		if verifyCodeQueue.Len() != 0 {
			msg_, err := verifyCodeQueue.Get(1)
			if err != nil {
				logger.ErrorLog(err)
				continue
			}
			msg := msg_[0].(VerifyCodeMsg)
			//
			err = sendVerifyCode(msg.Email, msg.Code)
			if err != nil {
				logger.ErrorLog(err)
				continue
			}
		}
	}
}

// 发送验证码
func sendVerifyCode(recv_email string, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(email, name))
	msg.SetHeader("To", recv_email)
	msg.SetHeader("Subject", "账号注册验证码")
	msg.SetBody("text/plain", fmt.Sprintf("【%s】您的账号注册验证码为:%s，5 分钟内有效，感谢您使用公司相关服务。", name, code))
	err := sender.DialAndSend(msg)
	return err
}

// 向队列增添一份邮件
func AddRegistVerifyCodeMsg(msg *VerifyCodeMsg) error {
	return verifyCodeQueue.Put(*msg)
}
