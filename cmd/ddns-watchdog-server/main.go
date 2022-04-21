package main

import (
	"ddns-watchdog/internal/common"
	"ddns-watchdog/internal/server"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	installOption   = flag.Bool("I", false, "安装服务并退出")
	uninstallOption = flag.Bool("U", false, "卸载服务并退出")
	version         = flag.Bool("v", false, "查看当前版本并检查更新后退出")
	confPath        = flag.String("c", "", "指定配置文件路径 (最好是绝对路径)(路径有空格请放在双引号中间)")
	initOption      = flag.Bool("i", false, "初始化配置文件并退出")
)

func main() {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		server.ConfDirectoryName = common.FormatDirectoryPath(*confPath)
	}

	// 初始化配置
	if *initOption {
		err := RunInit()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err := server.Install()
		if err != nil {
			log.Fatal(err)
		}
		// 初始化配置
		err = RunInit()
		if err != nil {
			log.Fatal(err)
		}
		return
	case *uninstallOption:
		err := server.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 加载配置
	conf := server.ServerConf{}
	err := common.LoadAndUnmarshal(server.ConfDirectoryName+"/"+server.ConfFileName, &conf)
	if err != nil {
		log.Fatal(err)
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
		sendJson, err := json.Marshal(info)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.WriteString(w, string(sendJson))
		if err != nil {
			log.Fatal(err)
		}
	}

	// 路径绑定处理变量
	http.HandleFunc("/", ddnsServerHandler)

	// 启动监听
	log.Println("Work on", conf.Port)
	err = http.ListenAndServe(conf.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func RunInit() (err error) {
	conf := server.ServerConf{}
	conf.Port = ":10032"
	conf.IsRoot = false
	conf.RootServerAddr = "https://yzyweb.cn/ddns-watchdog"
	err = common.MarshalAndSave(conf, server.ConfDirectoryName+"/"+server.ConfFileName)
	if err != nil {
		return
	}
	return
}
