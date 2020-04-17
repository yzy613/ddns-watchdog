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
	version = flag.Bool("version", false, "查看当前版本")
)

func main() {
	flag.Parse()
	switch {
	case *installMode:
		server.Install()
		return
	case *uninstallMode:
		server.Uninstall()
		return
	}

	// 加载配置
	conf := server.ServerConf{}
	getErr := common.IsDirExistAndCreate("./conf")
	if getErr != nil {
		fmt.Println(getErr)
		return
	}
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
		conf.IsRoot = false
		conf.RootServerAddr = "https://yzyweb.cn/ddns"
		if server.IsWindows() == true {
			getErr = common.MarshalAndSave(conf, "./conf/server.json")
		} else {
			getErr = common.MarshalAndSave(conf, server.ConfPath+"server.json")
		}
		if getErr != nil {
			fmt.Println(getErr)
		}
	}
	if *version {
		server.CheckLatestVersion(conf)
		return
	}


	ddnsServerHandler := func(w http.ResponseWriter, req *http.Request) {
		info := common.PublicInfo{
			IP:            server.GetClientIP(req),
			LatestVersion: server.GetLatestVersion(conf),
		}
		sendJson, getErr := json.Marshal(info)
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
