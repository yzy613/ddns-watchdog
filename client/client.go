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
	res, err := http.Get(conf.WebAddr)
	if err != nil {
		return "N/A (请检查网络连接)"
	}
	defer res.Body.Close()
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "N/A (数据包错误)"
	}
	recv := common.PublicInfo{}
	err = json.Unmarshal(recvJson, &recv)
	if err != nil {
		return "N/A (数据包错误)"
	}
	if recv.Version == "" {
		return "N/A (没有获取到版本信息)"
	}
	return recv.Version
}

func CheckLatestVersion(conf ClientConf) {
	LatestVersion := GetLatestVersion(conf)
	fmt.Println("当前版本 ", common.LocalVersion)
	fmt.Println("最新版本 ", LatestVersion)
	switch {
	case strings.Contains(LatestVersion, "N/A"):
		fmt.Println("\n需要手动检查更新，请前往 " + common.ProjectAddr + " 查看")
	case common.CompareVersionString(LatestVersion, common.LocalVersion):
		fmt.Println("\n发现新版本，请前往 " + common.ProjectAddr + " 下载")
	}
}

func postman(url, src string) (dst []byte, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "ddns-client/"+common.LocalVersion+" ()")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dst, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}
