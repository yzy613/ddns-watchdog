package main

import (
	"ddns/client"
	"ddns/common"
	"flag"
	"fmt"
)

var (
	enforcement = flag.Bool("f", false, "强制检查 DNS 解析记录")
	moreTips    = flag.Bool("mt", false, "显示更多的提示")
	version     = flag.Bool("version", false, "查看当前版本并检查更新")
)

func main() {
	flag.Parse()

	// 加载配置
	conf := client.ClientConf{}
	getErr := common.IsDirExistAndCreate("./conf")
	if getErr != nil {
		fmt.Println(getErr)
		return
	}
	getErr = common.LoadAndUnmarshal("./conf/client.json", &conf)
	if getErr != nil {
		fmt.Println(getErr)
		// 这里不能 return
	}

	saveMark := false
	if conf.WebAddr == "" {
		conf.WebAddr = common.RootServer
		saveMark = true
	}
	if conf.LatestIP == "" {
		conf.LatestIP = "0:0:0:0:0:0:0:0"
		conf.IsIPv6 = true
		saveMark = true
	}
	if saveMark {
		getErr = common.MarshalAndSave(conf, "./conf/client.json")
		if getErr != nil {
			fmt.Println(getErr)
		}
		// 以后可能不需要
		if !conf.Services.DNSPod && !conf.Services.Aliyun {
			fmt.Println("请打开客户端配置文件 client.json 启用需要使用的服务并重新启动")
			// 需要用户手动设置
			return
		}
	}
	if *version {
		conf.CheckLatestVersion()
		return
	}
	// 以后可能不需要
	if !conf.Services.DNSPod && !conf.Services.Aliyun {
		fmt.Println("请打开客户端配置文件 client.json 启用需要使用的服务并重新启动")
	}

	// 对比上一次的 IP
	acquiredIP, isIPv6, getErr := client.GetOwnIP(conf.WebAddr)
	if getErr != nil {
		fmt.Println(getErr)
		return
	}
	switch {
	case acquiredIP != conf.LatestIP || *enforcement:
		if acquiredIP != conf.LatestIP {
			conf.LatestIP = acquiredIP
			conf.IsIPv6 = isIPv6
			getErr = common.MarshalAndSave(conf, "./conf/client.json")
			if getErr != nil {
				fmt.Println(getErr)
				return
			}
		}
		waitDNSPod := make(chan bool)
		waitAliyun := make(chan bool)
		if conf.Services.DNSPod {
			go startDNSPod(acquiredIP, waitDNSPod)
		}
		if conf.Services.Aliyun {
			go startAliyun(acquiredIP, waitAliyun)
		}
		switch {
		case conf.Services.DNSPod && conf.Services.Aliyun:
			_, _ = <-waitDNSPod, <-waitAliyun
		case conf.Services.DNSPod:
			<-waitDNSPod
		case conf.Services.Aliyun:
			<-waitAliyun
		}
	case *moreTips:
		fmt.Println("因为获取的 IP 和当前本地记录的 IP 相同，所以跳过检查解析记录\n" +
			"若需要强制检查 DNS 解析记录，请添加启动参数 -f")
	}
}

func startDNSPod(ipAddr string, done chan bool) {
	err := client.DNSPod(ipAddr)
	if err != nil {
		fmt.Println(err)
	}
	done <- true
}

func startAliyun(ipAddr string, done chan bool) {
	err := client.Aliyun(ipAddr)
	if err != nil {
		fmt.Println(err)
	}
	done <- true
}
