package main

import (
	"errors"
	"flag"
	"github.com/yzy613/ddns-watchdog/client"
	"github.com/yzy613/ddns-watchdog/common"
	"log"
	"time"
)

var (
	installOption   = flag.Bool("install", false, "安装服务")
	uninstallOption = flag.Bool("uninstall", false, "卸载服务")
	enforcement     = flag.Bool("f", false, "强制检查 DNS 解析记录")
	version         = flag.Bool("version", false, "查看当前版本并检查更新")
	initOption      = flag.String("init", "", "有选择地初始化配置文件，可以组合使用 (例 01)\n"+
		"0 -> "+client.ConfFileName+"\n"+
		"1 -> "+client.DNSPodConfFileName+"\n"+
		"2 -> "+client.AliDNSConfFileName+"\n"+
		"3 -> "+client.CloudflareConfFileName)
	confPath = flag.String("conf_path", "", "指定配置文件路径 (最好是绝对路径)(路径有空格请放在双引号中间)")
	conf     = client.ClientConf{}
	dpc      = client.DNSPodConf{}
	adc      = client.AliDNSConf{}
	cfc      = client.CloudflareConf{}
)

func main() {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		tempStr := *confPath
		if tempStr[len(tempStr)-1:] != "/" {
			tempStr = tempStr + "/"
		}
		client.ConfPath = tempStr
	}

	// 有选择地初始化配置文件
	if *initOption != "" {
		for _, event := range *initOption {
			err := runInit(string(event))
			if err != nil {
				log.Fatal(err)
			}
		}
		return
	}

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err := client.Install()
		if err != nil {
			log.Fatal(err)
		}
		return
	case *uninstallOption:
		err := client.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 加载客户端配置
	err := common.LoadAndUnmarshal(client.ConfPath+client.ConfFileName, &conf)
	if err != nil {
		log.Fatal(err)
	}
	// 检查版本
	if *version {
		conf.CheckLatestVersion()
		return
	}
	// 检查启用 ddns
	if !conf.Services.DNSPod && !conf.Services.AliDNS && !conf.Services.Cloudflare {
		log.Fatal("请打开客户端配置文件 " + client.ConfPath + client.ConfFileName + " 启用需要使用的服务并重新启动")
	}
	servicesErr := false
	// 加载服务配置
	if conf.Services.DNSPod {
		err = common.LoadAndUnmarshal(client.ConfPath+client.DNSPodConfFileName, &dpc)
		if err != nil {
			log.Fatal(err)
		}
		if dpc.Id == "" || dpc.Token == "" || dpc.Domain == "" {
			log.Println("请打开配置文件 " + client.ConfPath + client.DNSPodConfFileName + " 检查你的 id, token, domain 并重新启动")
			servicesErr = true
		}
		if len(dpc.SubDomain) == 0 {
			log.Println("请打开配置文件 " + client.ConfPath + client.DNSPodConfFileName + " 检查你的 sub_domain 并重新启动")
			servicesErr = true
		}
	}
	if conf.Services.AliDNS {
		err = common.LoadAndUnmarshal(client.ConfPath+client.AliDNSConfFileName, &adc)
		if err != nil {
			log.Fatal(err)
		}
		if adc.AccessKeyId == "" || adc.AccessKeySecret == "" || adc.Domain == "" {
			log.Println("请打开配置文件 " + client.ConfPath + client.AliDNSConfFileName + " 检查你的 accesskey_id, accesskey_secret, domain 并重新启动")
			servicesErr = true
		}
		if len(adc.SubDomain) == 0 {
			log.Println("请打开配置文件 " + client.ConfPath + client.AliDNSConfFileName + " 检查你的 sub_domain 并重新启动")
			servicesErr = true
		}
	}
	if conf.Services.Cloudflare {
		err = common.LoadAndUnmarshal(client.ConfPath+client.CloudflareConfFileName, &cfc)
		if err != nil {
			log.Fatal(err)
		}
		if cfc.Email == "" || cfc.APIKey == "" || cfc.ZoneID == "" {
			log.Println("请打开配置文件 " + client.ConfPath + client.CloudflareConfFileName + " 检查你的 email, api_key, zone_id 并重新启动")
			servicesErr = true
		}
		for len(cfc.Domain) == 0 {
			log.Println("请打开配置文件 " + client.ConfPath + client.CloudflareConfFileName + " 检查你的 domain 并重新启动")
			servicesErr = true
		}
	}
	if servicesErr {
		log.Fatal("请检查以上错误")
	}

	// 周期循环
	waitCheckDone := make(chan bool, 1)
	if conf.CheckCycle == 0 {
		go asyncCheck(&conf, waitCheckDone)
		<-waitCheckDone
	} else {
		cycle := time.NewTicker(time.Duration(conf.CheckCycle) * time.Minute)
		for {
			go asyncCheck(&conf, waitCheckDone)
			<-waitCheckDone
			<-cycle.C
		}
	}
}

func runInit(event string) (err error) {
	switch event {
	case "0":
		conf.APIUrl = common.DefaultAPIServer
		conf.CheckCycle = 0
		err = common.MarshalAndSave(conf, client.ConfPath+client.ConfFileName)
		if err != nil {
			return
		}
		log.Println("初始化 " + client.ConfPath + client.ConfFileName)
	case "1":
		dpc.SubDomain = append(dpc.SubDomain, "example")
		err = common.MarshalAndSave(dpc, client.ConfPath+client.DNSPodConfFileName)
		if err != nil {
			return
		}
		log.Println("初始化 " + client.ConfPath + client.DNSPodConfFileName)
	case "2":
		adc.SubDomain = append(adc.SubDomain, "example")
		err = common.MarshalAndSave(adc, client.ConfPath+client.AliDNSConfFileName)
		if err != nil {
			return
		}
		log.Println("初始化 " + client.ConfPath + client.AliDNSConfFileName)
	case "3":
		cfc.Domain = append(cfc.Domain, "example")
		err = common.MarshalAndSave(cfc, client.ConfPath+client.CloudflareConfFileName)
		if err != nil {
			return
		}
		log.Println("初始化 " + client.ConfPath + client.CloudflareConfFileName)
	default:
		err = errors.New("你初始化了一个寂寞")
	}
	return
}

func asyncCheck(conf *client.ClientConf, done chan bool) {
	// 获取 IP
	acquiredIP, err := client.GetOwnIP(conf.APIUrl, conf.EnableNetworkCard, conf.NetworkCard)
	if err != nil {
		log.Fatal(err)
	}

	if acquiredIP != conf.LatestIP || *enforcement {
		if acquiredIP != conf.LatestIP {
			conf.LatestIP = acquiredIP
		}
		servicesCount := 0
		if conf.Services.DNSPod {
			servicesCount++
		}
		if conf.Services.AliDNS {
			servicesCount++
		}
		if conf.Services.Cloudflare {
			servicesCount++
		}
		waitServicesDone := make(chan bool, servicesCount)
		if conf.Services.DNSPod {
			go asyncDNSPod(acquiredIP, waitServicesDone)
		}
		if conf.Services.AliDNS {
			go asyncAliDNS(acquiredIP, waitServicesDone)
		}
		if conf.Services.Cloudflare {
			go asyncCloudflare(acquiredIP, waitServicesDone)
		}
		for i := 0; i < servicesCount; i++ {
			<-waitServicesDone
		}
	}
	done <- true
}

func asyncDNSPod(ipAddr string, done chan bool) {
	msg, err := client.DNSPod(dpc, ipAddr)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}

func asyncAliDNS(ipAddr string, done chan bool) {
	msg, err := client.AliDNS(adc, ipAddr)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}

func asyncCloudflare(ipAddr string, done chan bool) {
	msg, err := client.Cloudflare(cfc, ipAddr)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}
