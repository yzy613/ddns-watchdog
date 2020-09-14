package main

import (
	"flag"
	"log"
	"time"
	"watchdog-ddns/client"
	"watchdog-ddns/common"
)

var (
	enforcement = flag.Bool("f", false, "强制检查 DNS 解析记录")
	version     = flag.Bool("version", false, "查看当前版本并检查更新")
	initOption  = flag.Bool("init", false, "初始化配置文件")
	confPath    = flag.String("conf_path", "", "手动设置配置文件路径（绝对路径）（有空格用双引号）")
	conf        = client.ClientConf{}
	dpc         = client.DNSPodConf{}
	ayc         = client.AliyunConf{}
	cfc         = client.CloudflareConf{}
)

func main() {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		client.ConfPath = *confPath
	}

	// 初始化配置
	if *initOption {
		conf.APIUrl = common.DefaultAPIServer
		conf.CheckCycle = 0
		err := common.MarshalAndSave(conf, client.ConfPath+"/client.json")
		if err != nil {
			log.Fatal(err)
		}
		err = common.MarshalAndSave(dpc, client.ConfPath+"/dnspod.json")
		if err != nil {
			log.Fatal(err)
		}
		err = common.MarshalAndSave(ayc, client.ConfPath+"/aliyun.json")
		if err != nil {
			log.Fatal(err)
		}
		err = common.MarshalAndSave(cfc, client.ConfPath+"/cloudflare.json")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 加载客户端配置
	err := common.LoadAndUnmarshal(client.ConfPath+"/client.json", &conf)
	if err != nil {
		log.Fatal(err)
	}
	// 检查启用 ddns
	if !conf.Services.DNSPod && !conf.Services.Aliyun && !conf.Services.Cloudflare {
		log.Fatal("请打开客户端配置文件 " + client.ConfPath + "/client.json 启用需要使用的服务并重新启动")
	}
	servicesErr := false
	// 加载服务配置
	if conf.Services.DNSPod {
		err = common.LoadAndUnmarshal(client.ConfPath+"/dnspod.json", &dpc)
		if err != nil {
			log.Fatal(err)
		}
		if dpc.Id == "" || dpc.Token == "" || dpc.Domain == "" || dpc.SubDomain == "" {
			log.Println("请打开配置文件 " + client.ConfPath + "/dnspod.json 检查你的 id, token, domain, sub_domain 并重新启动")
			servicesErr = true
		}
	}
	if conf.Services.Aliyun {
		err = common.LoadAndUnmarshal(client.ConfPath+"/aliyun.json", &ayc)
		if err != nil {
			log.Fatal(err)
		}
		if ayc.AccessKeyId == "" || ayc.AccessKeySecret == "" || ayc.Domain == "" || ayc.SubDomain == "" {
			log.Println("请打开配置文件 " + client.ConfPath + "/aliyun.json 检查你的 accesskey_id, accesskey_secret, domain, sub_domain 并重新启动")
			servicesErr = true
		}
	}
	if conf.Services.Cloudflare {
		err = common.LoadAndUnmarshal(client.ConfPath+"/cloudflare.json", &cfc)
		if err != nil {
			log.Fatal(err)
		}
		if cfc.Email == "" || cfc.APIKey == "" || cfc.ZoneID == "" || cfc.Domain == "" {
			log.Println("请打开配置文件 " + client.ConfPath + "/cloudflare.json 检查你的 email, api_key, zone_id, domain 并重新启动")
			servicesErr = true
		}
	}
	if servicesErr {
		log.Fatal("请检查以上错误")
	}

	// 检查版本
	if *version {
		conf.CheckLatestVersion()
		return
	}

	// 周期循环
	skip := false
	waitCheckDone := make(chan bool, 1)
	for !skip {
		go asyncCheck(&conf, waitCheckDone)
		if conf.CheckCycle != 0 {
			time.Sleep(time.Duration(conf.CheckCycle) * time.Minute)
		} else {
			skip = true
		}
		<-waitCheckDone
	}
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
		if conf.Services.Aliyun {
			servicesCount++
		}
		if conf.Services.Cloudflare {
			servicesCount++
		}
		waitServicesDone := make(chan bool, servicesCount)
		if conf.Services.DNSPod {
			go asyncDNSPod(acquiredIP, waitServicesDone)
		}
		if conf.Services.Aliyun {
			go asyncAliyun(acquiredIP, waitServicesDone)
		}
		if conf.Services.Cloudflare {
			go asyncCloudflare(acquiredIP, waitServicesDone)
		}
		for i := 0; i < servicesCount; i++ {
			<-waitServicesDone
		}
	} else {
		log.Println("当前获取的 IP 和上一次获取的 IP 相同")
	}
	done <- true
}

func asyncDNSPod(ipAddr string, done chan bool) {
	err := client.DNSPod(dpc, ipAddr)
	if err != nil {
		log.Println(err)
	}
	done <- true
}

func asyncAliyun(ipAddr string, done chan bool) {
	err := client.Aliyun(ayc, ipAddr)
	if err != nil {
		log.Println(err)
	}
	done <- true
}

func asyncCloudflare(ipAddr string, done chan bool) {
	err := client.Cloudflare(cfc, ipAddr)
	if err != nil {
		log.Println(err)
	}
	done <- true
}
