package main

import (
	"ddns-watchdog/internal/common"
	"ddns-watchdog/internal/server"
	"encoding/json"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"net/http"
)

var (
	installOption   = flag.BoolP("install", "I", false, "安装服务并退出")
	uninstallOption = flag.BoolP("uninstall", "U", false, "卸载服务并退出")
	version         = flag.BoolP("version", "v", false, "查看当前版本并检查更新后退出")
	confPath        = flag.StringP("conf", "c", "", "指定配置文件目录 (目录有空格请放在双引号中间)")
	initOption      = flag.BoolP("init", "i", false, "初始化配置文件并退出")
)

func main() {
	flag.Parse()
	// 加载自定义配置文件目录
	if *confPath != "" {
		server.ConfDirectoryName = common.FormatDirectoryPath(*confPath)
	}

	exit := false

	// 初始化配置
	if *initOption {
		err := RunInit()
		if err != nil {
			log.Fatal(err)
		}
		exit = true
	}

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err := server.Install()
		if err != nil {
			log.Fatal(err)
		}
		exit = true
	case *uninstallOption:
		err := server.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		exit = true
	}

	if exit {
		return
	}

	// 进入工作流程
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
		w.Header().Add("Cache-Control", "no-store")
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
	if conf.TLS.Enable {
		log.Println("Work on", conf.Port, "with TLS")
		err = http.ListenAndServeTLS(conf.Port, server.ConfDirectoryName+"/"+conf.TLS.CertFile, server.ConfDirectoryName+"/"+conf.TLS.KeyFile, nil)
	} else {
		log.Println("Work on", conf.Port)
		err = http.ListenAndServe(conf.Port, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func RunInit() (err error) {
	conf := server.ServerConf{
		Port:           ":10032",
		IsRoot:         false,
		RootServerAddr: "https://yzyweb.cn/ddns-watchdog",
	}
	err = common.MarshalAndSave(conf, server.ConfDirectoryName+"/"+server.ConfFileName)
	if err != nil {
		return
	}
	log.Println("初始化 " + server.ConfDirectoryName + "/" + server.ConfFileName)
	return
}
