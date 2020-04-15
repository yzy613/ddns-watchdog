package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"runtime"
)

type ipInfoFormat struct {
	Ip string `json:"ip"`
}

type confFormat struct {
	Port string `json:"port"`
}

func main() {
	// 加载配置
	var (
		confSrc []byte
		getErr error
	)
	if runtime.GOOS == "windows" {
		confSrc, getErr = ioutil.ReadFile("./conf/server.json")
	} else {
		confSrc, getErr = ioutil.ReadFile("/opt/ddns/conf/server.json")
	}
	if getErr != nil {
		fmt.Println(getErr)
	}
	conf := confFormat{}
	getErr = json.Unmarshal(confSrc, &conf)
	if getErr != nil {
		fmt.Println(getErr)
	}

	// 处理请求
	ddnsServerHandler := func(w http.ResponseWriter, req *http.Request) {
		// 获取IP(:port)
		var ipFinal string
		ipFinal = req.Header.Get("X-Real-IP")
		if ipFinal == "" {
			ipFinal = req.Header.Get("X-Forwarded-For")
		}
		if ipFinal == "" {
			ipSrc := req.RemoteAddr
			// 对ip:port切片
			if req.RemoteAddr[0] == '[' {
				// IPv6
				ipFinal = strings.Split(ipSrc, "]:")[0]
				ipFinal = fmt.Sprint(ipFinal, "]")
			} else {
				// IPv4
				ipFinal = strings.Split(req.RemoteAddr, ":")[0]
			}
		}

		// 编码为 json 并发送
		ipInfo := ipInfoFormat{Ip: ipFinal}
		figJson, getErr := json.Marshal(ipInfo)
		if getErr != nil {
			fmt.Println(getErr)
		}
		io.WriteString(w, string(figJson))
	}

	// 路径绑定处理变量
	http.HandleFunc("/", ddnsServerHandler)

	// 启动监听
	if conf.Port == "" {
		conf.Port = ":10032"
	}
	fmt.Println("Work on ", conf.Port)
	http.ListenAndServe(conf.Port, nil)
}
