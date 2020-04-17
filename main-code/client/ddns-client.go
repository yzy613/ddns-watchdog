package main

import (
	"ddns/client"
	"ddns/common"
	"encoding/json"
	"flag"
	"fmt"
)

var forcibly = flag.Bool("f", false, "强制刷新 DNS 解析记录")

func main() {
	flag.Parse()

	// 加载配置
	conf := common.ClientConf{}
	getErr := common.LoadAndUnmarshal("./conf/client.json", &conf)
	if getErr != nil {
		fmt.Println(getErr)
	}

	// 对比上一次的 IP
	if conf.WebAddr == "" {
		conf.WebAddr = "https://yzyweb.cn/ddns"
	}
	ipAddr, getErr := client.GetOwnIP(conf.WebAddr)
	if getErr != nil {
		fmt.Println(getErr)
	}
	if ipAddr != conf.LastIP || *forcibly {
		conf.LastIP = ipAddr
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
	} else {
		fmt.Println("你的公网 IP 没有变化")
	}
}

func DNSPod(ipAddr string) {
	dpc:= client.DNSPodConf{}
	getErr := common.LoadAndUnmarshal("conf/dnspod.json", &dpc)
	if getErr != nil {
		fmt.Println(getErr)
	}
	if dpc.Id == "" || dpc.Token == "" {
		fmt.Println("请打开配置文件 dnspod.json 填入你的 Id 或 Token")
		getErr = common.MarshalAndSave(dpc, "conf/dnspod.json")
		if getErr != nil {
			fmt.Println(getErr)
		}
		return
	}
	DPContent, recvMsg, getErr := client.DNSPod(dpc, ipAddr)
	if getErr != nil {
		fmt.Println(getErr)
	}
	if recvMsg != "" {
		fmt.Println(recvMsg)
	}
	recvMap := make(map[string]interface{})
	getErr = json.Unmarshal(DPContent, &recvMap)
	if getErr != nil {
		fmt.Println(getErr)
	}
	getErr = common.MarshalAndSave(recvMap, "conf/recv.json")
	/*getErr = ioutil.WriteFile("conf/recv.json", DPContent, 0666)
	if getErr != nil {
		fmt.Println(getErr)
	}*/
}