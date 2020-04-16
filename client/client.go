package client

import (
	"ddns/common"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func DNSPod(dps common.DNSPodSecret) {
}

func GetOwnIP(webAddr string) (ipAddr string, err error) {
	if webAddr == "" {
		webAddr = "https://yzyweb.cn/ddns"
	}
	res, err := http.Get(webAddr)
	defer res.Body.Close()
	if err != nil {
		return
	}
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var ipInfo common.IpInfoFormat
	err = json.Unmarshal(recvJson, &ipInfo)
	if err != nil {
		return
	}
	ipAddr = ipInfo.Ip
	return
}
