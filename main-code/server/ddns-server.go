package main

import (
	"ddns/common"
	"ddns/server"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	installMode   = flag.Bool("install", false, "安装服务")
	uninstallMode = flag.Bool("uninstall", false, "卸载服务")
	version       = flag.Bool("version", false, "查看当前版本并检查更新")
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
	var getErr error
	if server.IsWindows() {
		getErr = common.IsDirExistAndCreate("./conf/")
	} else {
		getErr = common.IsDirExistAndCreate(server.ConfPath)
	}
	if getErr != nil {
		log.Fatal(getErr)
	}
	if server.IsWindows() {
		getErr = common.LoadAndUnmarshal("./conf/server.json", &conf)
	} else {
		getErr = common.LoadAndUnmarshal(server.ConfPath+"server.json", &conf)
	}
	if getErr != nil {
		log.Println(getErr)
		// 这里不能 return
	}

	saveMark := false
	if conf.Port == "" {
		conf.Port = ":10032"
		saveMark = true
	}
	if conf.RootServerAddr == "" && !conf.IsRoot {
		conf.IsRoot = false
		conf.RootServerAddr = "https://yzyweb.cn/ddns"
		saveMark = true
	}
	if saveMark {
		if server.IsWindows() {
			getErr = common.MarshalAndSave(conf, "./conf/server.json")
		} else {
			getErr = common.MarshalAndSave(conf, server.ConfPath+"server.json")
		}
		if getErr != nil {
			log.Fatal(getErr)
		}
	}
	if *version {
		conf.CheckLatestVersion()
		return
	}

	ddnsServerHandler := func(w http.ResponseWriter, req *http.Request) {
		info := common.PublicInfo{
			IP:      server.GetClientIP(req),
			Version: conf.GetLatestVersion(),
		}
		sendJson, getErr := json.Marshal(info)
		if getErr != nil {
			log.Fatal(getErr)
		}
		_, getErr = io.WriteString(w, string(sendJson))
		if getErr != nil {
			log.Fatal(getErr)
		}
	}

	// 路径绑定处理变量
	http.HandleFunc("/", ddnsServerHandler)

	// 启动监听
	log.Println("Work on ", conf.Port)
	getErr = http.ListenAndServe(conf.Port, nil)
	if getErr != nil {
		log.Fatal(getErr)
	}
}
