package logger

import (
	config "Server/Config"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getTimeString() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())
	hour := strconv.Itoa(now.Hour())
	minute := strconv.Itoa(now.Minute())
	second := strconv.Itoa(now.Second())
	time_ := year + "-" + month + "-" + day + "-" + hour + "点" + minute + "分" + second + "秒"
	return time_
}

// 对所有的log进行初始化
func LogInit() {
	debug := config.Configuration["server"].(map[string]interface{})["mode"].(string) == "debug"
	//debug模式直接输出,release模式输出到文件
	if debug {

	} else {
		//
		dir := "LogFile/" //指定存放目录
		var err_s []error
		logfile, err := os.Create(dir + "log.txt")
		err_s = append(err_s, err)
		errfile, err := os.Create(dir + "errlog.txt")
		err_s = append(err_s, err)
		//检查出错
		for _, err := range err_s {
			if err != nil {
				panic("日志文件打开失败!")
			}
		}
		//设置输出流
		os.Stdout = logfile
		os.Stderr = errfile
		//设置gin的输出
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}
}

// 普通log
func Log(msg interface{}) {
	fmt.Printf("Log:%v:%v\n", getTimeString(), msg)
}
