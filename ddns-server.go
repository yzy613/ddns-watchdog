package main

import (
	"ddns/common"
	"ddns/server"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	installMode   = flag.Bool("install", false, "安装到系统")
	uninstallMode = flag.Bool("uninstall", false, "卸载服务")
)

func beforeStart() {
	flag.Parse()
	if *installMode == true {
		server.Install()
		return
	}
	if *uninstallMode == true {
		server.Uninstall()
		return
	}
}

func main() {
	beforeStart()

	// 加载配置
	var getErr error
	conf := common.ServerConf{}
	if server.IsWindows() == true {
		getErr = common.LoadAndUnmarshal("./conf/server.json", &conf)
	} else {
		getErr = common.LoadAndUnmarshal(server.ConfPath+"server.json", &conf)
	}
	if getErr != nil {
		fmt.Println(getErr)
	}
	if conf.Port == "" {
		conf.Port = ":10032"
	}
	if server.IsWindows() == true {
		getErr = common.MarshalAndSave(conf, "./conf/server.json")
	} else {
		getErr = common.MarshalAndSave(conf, server.ConfPath+"server.json")
	}
	if getErr != nil {
		fmt.Println(getErr)
	}

	// 处理请求
	ddnsServerHandler := func(w http.ResponseWriter, req *http.Request) {
		// 编码为 json 并发送
		ipInfo := common.IpInfoFormat{Ip: server.GetIP(req)}
		sendJson, getErr := json.Marshal(ipInfo)
		if getErr != nil {
			fmt.Println(getErr)
		}
		_, getErr = io.WriteString(w, string(sendJson))
		if getErr != nil {
			fmt.Println(getErr)
		}
	}

	// 路径绑定处理变量
	http.HandleFunc("/", ddnsServerHandler)

	// 启动监听
	fmt.Println("Work on ", conf.Port)
	getErr = http.ListenAndServe(conf.Port, nil)
	if getErr != nil {
		fmt.Println(getErr)
	}
}
