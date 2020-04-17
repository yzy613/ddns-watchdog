package main

import (
	"ddns/client"
	"ddns/common"
	"flag"
	"fmt"
)

var forcibly = flag.Bool("f", false, "强制刷新 DNS 解析记录")

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
		if conf.WebAddr == "" {
			conf.WebAddr = "https://yzyweb.cn/ddns"
		}
		if conf.LatestIP == "" {
			conf.LatestIP = "0.0.0.0"
		}
		common.MarshalAndSave(conf, "./conf/client.json")
		fmt.Println(getErr)
		fmt.Println("如果显示错误为 unexpected end of JSON input\n请打开客户端配置文件 client.json 填写信息")
		return
	}

	// 对比上一次的 IP
	if conf.WebAddr == "" {
		conf.WebAddr = "https://yzyweb.cn/ddns"
	}
	ipAddr, isIPv6, getErr := client.GetOwnIP(conf.WebAddr)
	if getErr != nil {
		fmt.Println(getErr)
	}
	if ipAddr != conf.LatestIP || *forcibly {
		conf.LatestIP = ipAddr
		conf.IsIPv6 = isIPv6
		fmt.Println("你的公网 IP: ", ipAddr)
		getErr = common.MarshalAndSave(conf, "./conf/client.json")
		if getErr != nil {
			fmt.Println(getErr)
		}
		if conf.EnableDdns {
			if conf.DNSPod {
				DNSPod(ipAddr)
			}
		}
	} /* else {
		fmt.Println("你的公网 IP 没有变化")
	}*/
}

func DNSPod(ipAddr string) {
	dpc := client.DNSPodConf{}
	getErr := common.LoadAndUnmarshal("./conf/dnspod.json", &dpc)
	if getErr != nil {
		fmt.Println(getErr)
		fmt.Println("如果显示错误为 unexpected end of JSON input\n请打开配置文件 dnspod.json 填入你的 Id 和 Token")
		return
	}
	if dpc.Id == "" || dpc.Token == "" {
		fmt.Println("请打开配置文件 dnspod.json 填入你的 Id 和 Token")
		getErr = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		if getErr != nil {
			fmt.Println(getErr)
		}
		return
	}
	if ipAddr[0] == '[' {
		dpc.RecordType = "AAAA"
	} else {
		dpc.RecordType = "A"
	}
	getErr = client.DNSPod(&dpc, ipAddr)
	if getErr != nil {
		fmt.Println(getErr)
	}
}
