package main

import (
	"ddns-watchdog/internal/client"
	"ddns-watchdog/internal/common"
	"ddns-watchdog/internal/server"
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"time"
)

var (
	installOption   = flag.BoolP("install", "I", false, "安装服务并退出")
	addToken        = flag.BoolP("add-token", "a", false, "添加 token 到白名单")
	generateToken   = flag.BoolP("generate-token", "g", false, "生成 token 并输出")
	tokenLength     = flag.IntP("token-length", "l", 48, "指定生成 token 的长度")
	token           = flag.StringP("token", "t", "", "指定 token (长度在 [16,127] 之间，支持 UTF-8 字符)")
	message         = flag.StringP("message", "m", "undefined", "备注 token 信息")
	uninstallOption = flag.BoolP("uninstall", "U", false, "卸载服务并退出")
	version         = flag.BoolP("version", "v", false, "查看当前版本并检查更新后退出")
	confPath        = flag.StringP("conf", "c", "", "指定配置文件目录 (目录有空格请放在双引号中间)")
	initOption      = flag.StringP("init", "i", "", "有选择地初始化配置文件并退出，可以组合使用 (例 01)\n"+
		"0 -> "+server.ConfFileName+"\n"+
		"1 -> "+server.WhitelistFileName+"\n"+
		"2 -> "+server.ServiceConfFileName)
	insecure = flag.BoolP("insecure", "k", false, "使用 https 链接时不检查 TLS 证书合法性")
)

func main() {
	// 处理 flag
	exit, err := processFlag()
	if err != nil {
		log.Fatal(err)
		return
	}
	if exit {
		return
	}

	// 加载白名单
	if server.Srv.CenterService {
		err = server.Service.LoadConf()
		if err != nil {
			log.Fatal(err)
		}
		// 路由绑定函数
		http.HandleFunc(server.Srv.Route.Center, server.RespCenterReq)
	}

	// 路由绑定函数
	http.HandleFunc(server.Srv.Route.GetIP, server.RespGetIPReq)

	// 设置超时参数
	httpSrv := http.Server{
		Addr:              server.Srv.ServerAddr,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       2 * time.Second,
	}

	// 启动监听
	if server.Srv.TLS.Enable {
		log.Println("Work on", server.Srv.ServerAddr, "with TLS")
		err = httpSrv.ListenAndServeTLS(server.ConfDirectoryName+"/"+server.Srv.TLS.CertFile, server.ConfDirectoryName+"/"+server.Srv.TLS.KeyFile)
	} else {
		log.Println("Work on", server.Srv.ServerAddr)
		err = httpSrv.ListenAndServe()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func processFlag() (exit bool, err error) {
	flag.Parse()
	if *confPath != "" {
		server.ConfDirectoryName = common.FormatDirectoryPath(*confPath)
	}

	// 初始化配置
	if *initOption != "" {
		for _, event := range *initOption {
			err = initConf(string(event))
			if err != nil {
				return
			}
		}
		exit = true
		return
	}

	// 加载配置
	err = server.Srv.LoadConf()
	if err != nil {
		return
	}

	// 版本信息
	if *version {
		server.Srv.CheckLatestVersion()
		exit = true
		return
	}

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err = server.Install()
		if err != nil {
			return
		}
		exit = true
		return
	case *uninstallOption:
		err = server.Uninstall()
		if err != nil {
			return
		}
		exit = true
		return
	}

	currentToken := ""
	// 获取 token
	switch {
	case *token != "":
		currentToken = *token
	case *generateToken:
		length := *tokenLength
		if length < 16 {
			length = 16
		}
		if length > 127 {
			length = 127
		}
		currentToken = server.GenerateToken(length)
		fmt.Printf("Token: %v\nMessage: %v\n", currentToken, *message)
		exit = true
	}

	// 添加 token 到白名单
	if *addToken {
		if currentToken == "" || len(currentToken) < 16 || len(currentToken) > 127 {
			err = errors.New("token 不符合要求")
		} else {
			err = server.AddTokenToWhitelist(currentToken, *message)
		}
		if err != nil {
			return
		}
		exit = true
		fmt.Printf("Added %v(%v) to whitelist.\n", currentToken, *message)
	}

	if *insecure {
		client.HttpsInsecure = *insecure
	}
	return
}

func initConf(event string) (err error) {
	msg := ""
	switch event {
	case "0":
		msg, err = server.Srv.InitConf()
	case "1":
		msg, err = server.InitWhitelist()
	case "2":
		msg, err = server.Service.InitConf()
	default:
		err = errors.New("你初始化了一个寂寞")
	}
	if err != nil {
		return err
	}
	log.Println(msg)
	return
}
