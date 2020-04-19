package client

import (
	"ddns/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetOwnIP(webAddr string) (ipAddr string, isIPv6 bool, err error) {
	if webAddr == "" {
		webAddr = common.RootServer
	}
	res, err := http.Get(webAddr)
	if err != nil {
		return
	}
	defer res.Body.Close()
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var ipInfo common.PublicInfo
	err = json.Unmarshal(recvJson, &ipInfo)
	if err != nil {
		return
	}
	ipAddr = ipInfo.IP
	// 判断 IP 类型
	if strings.Contains(ipAddr, ":") {
		isIPv6 = true
	} else {
		isIPv6 = false
	}
	return
}

func GetLatestVersion(conf ClientConf) string {
	res, getErr := http.Get(conf.WebAddr)
	if getErr != nil {
		return common.LocalVersion
	}
	defer res.Body.Close()
	recvJson, getErr := ioutil.ReadAll(res.Body)
	if getErr != nil {
		return common.LocalVersion
	}
	recv := common.PublicInfo{}
	getErr = json.Unmarshal(recvJson, &recv)
	if getErr != nil {
		return common.LocalVersion
	}
	return recv.Version
}

func CheckLatestVersion(conf ClientConf) {
	LatestVersion := GetLatestVersion(conf)
	fmt.Println("当前版本 ", common.LocalVersion)
	fmt.Println("最新版本 ", LatestVersion)
	if common.CompareVersionString(LatestVersion, common.LocalVersion) {
		fmt.Println("\n发现新版本，请前往 https://github.com/yzy613/ddns/releases 下载")
	}
}
