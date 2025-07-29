package logger

import "fmt"

func ErrorLog(msg interface{}) {
	fmt.Printf("Error:%v:%v\n", getTimeString(), msg)
}
