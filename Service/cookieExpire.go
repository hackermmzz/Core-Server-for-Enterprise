package service

import (
	"Server/database"
	"Server/logger"
	"time"
)

func RemoveCookieExpireService(config map[string]interface{}) {
	interval := int(config["interval"].(float64))
	for {
		//删除过期的cookie
		err := database.RemoveExpireCookie()
		if err != nil {
			logger.ErrorLog(err)
		}
		//休眠
		time.Sleep(time.Duration(interval * int(time.Minute)))
	}

}
