package main

import (
	"encoding/json"
	"flag"
	"github.com/yzy613/ddns-watchdog/common"
	"github.com/yzy613/ddns-watchdog/server"
	"io"
	"log"
	"net/http"
)

var (
	installOption   = flag.Bool("install", false, "安装服务")
	uninstallOption = flag.Bool("uninstall", false, "卸载服务")
	version         = flag.Bool("version", false, "查看当前版本并检查更新")
	confPath        = flag.String("conf_path", "", "手动设置配置文件路径（绝对路径）（有空格用双引号）")
	initOption      = flag.Bool("init", false, "初始化配置文件")
)

func main() {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		tempStr := *confPath
		if tempStr[len(tempStr)-1:] != "/" {
			tempStr = tempStr + "/"
		}
		server.ConfPath = tempStr
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
	err := common.LoadAndUnmarshal(server.ConfPath+"/server.json", &conf)
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
	err = common.MarshalAndSave(conf, server.ConfPath+"/server.json")
	if err != nil {
		return
	}
	return
}
