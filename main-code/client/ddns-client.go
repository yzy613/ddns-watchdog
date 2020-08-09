package main

import (
	"ddns/client"
	"ddns/common"
	"flag"
	"log"
)

var (
	enforcement = flag.Bool("f", false, "强制检查 DNS 解析记录")
	moreTips    = flag.Bool("mt", false, "显示更多的提示")
	version     = flag.Bool("version", false, "查看当前版本并检查更新")
	initOption  = flag.Bool("init", false, "初始化配置文件")
	confPath    = flag.String("conf_path", "", "手动设置配置文件路径（绝对路径）（有空格用双引号）")
)

func main() {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		client.ConfPath = *confPath
	}

	// 初始化配置
	if *initOption {
		conf := client.ClientConf{}
		conf.APIUrl = common.DefaultAPIServer
		conf.LatestIP = "0:0:0:0:0:0:0:0"
		conf.IsIPv6 = true
		err := common.MarshalAndSave(conf, client.ConfPath+"/client.json")
		if err != nil {
			log.Fatal(err)
		}
		dpc := client.DNSPodConf{}
		err = common.MarshalAndSave(dpc, client.ConfPath+"/dnspod.json")
		if err != nil {
			log.Fatal(err)
		}
		ayc := client.AliyunConf{}
		err = common.MarshalAndSave(ayc, client.ConfPath+"/aliyun.json")
		if err != nil {
			log.Fatal(err)
		}
		cfc := client.CloudflareConf{}
		err = common.MarshalAndSave(cfc, client.ConfPath+"/cloudflare.json")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 加载配置
	conf := client.ClientConf{}
	err := common.LoadAndUnmarshal(client.ConfPath+"/client.json", &conf)
	if err != nil {
		log.Fatal(err)
	}

	// 检查版本
	if *version {
		conf.CheckLatestVersion()
		return
	}

	// 检查启用 ddns
	if !conf.Services.DNSPod && !conf.Services.Aliyun && !conf.Services.Cloudflare {
		log.Fatal("请打开客户端配置文件 " + client.ConfPath + "/client.json 启用需要使用的服务并重新启动")
	}

	// 获取 IP
	acquiredIP, isIPv6, err := client.GetOwnIP(conf.APIUrl, conf.EnableNetworkCard, conf.NetworkCard)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case acquiredIP != conf.LatestIP || *enforcement:
		if acquiredIP != conf.LatestIP {
			conf.LatestIP = acquiredIP
			conf.IsIPv6 = isIPv6
			err = common.MarshalAndSave(conf, client.ConfPath+"/client.json")
			if err != nil {
				log.Fatal(err)
			}
		}
		waitDNSPod := make(chan bool)
		waitAliyun := make(chan bool)
		waitCloudflare := make(chan bool)
		if conf.Services.DNSPod {
			go startDNSPod(acquiredIP, waitDNSPod)
		}
		if conf.Services.Aliyun {
			go startAliyun(acquiredIP, waitAliyun)
		}
		if conf.Services.Cloudflare {
			go startCloudflare(acquiredIP, waitCloudflare)
		}
		if conf.Services.DNSPod {
			<-waitDNSPod
		}
		if conf.Services.Aliyun {
			<-waitAliyun
		}
		if conf.Services.Cloudflare {
			<-waitCloudflare
		}
	case *moreTips:
		log.Println("因为获取的 IP 和当前本地记录的 IP 相同，所以跳过检查解析记录\n" +
			"若需要强制检查 DNS 解析记录，请添加启动参数 -f")
	}
}

func startDNSPod(ipAddr string, done chan bool) {
	err := client.DNSPod(ipAddr)
	if err != nil {
		log.Println(err)
	}
	done <- true
}

func startAliyun(ipAddr string, done chan bool) {
	err := client.Aliyun(ipAddr)
	if err != nil {
		log.Println(err)
	}
	done <- true
}

func startCloudflare(ipAddr string, done chan bool) {
	err := client.Cloudflare(ipAddr)
	if err != nil {
		log.Fatal(err)
	}
	done <- true
}