package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var Configuration = make(map[string]interface{})

func ConfigInit() {
	//读取配置
	config, err := os.ReadFile("Config/config.db")
	if err != nil {
		fmt.Println("配置加载失败!", err)
		return
	}
	//处理config,去掉#后面的注释
	for {
		res := string(config)
		idx_beg := strings.Index(res, "#")
		if idx_beg == -1 {
			break
		}
		//从#开始一直找到换行或者到达文本末尾
		idx_end := idx_beg
		for idx_end < len(res) && res[idx_end] != '\n' {
			idx_end++
		}
		//重新拷贝
		res = res[:idx_beg] + res[min(idx_end+1, len(res)):]
		//
		config = []byte(res)
	}
	//以json格式解析config
	err = json.Unmarshal(config, &Configuration)
	if err != nil {
		fmt.Println("配置加载失败!", err)
		return
	}
	//
	//
	fmt.Println("配置加载成功!")
}
