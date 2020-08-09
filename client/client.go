package client

import (
	"ddns/common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var ConfPath = common.GetRunningPath() + "/conf"

func NetworkCardRespond() (map[string]string, error) {
	networkCardInfo := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		ipAddr, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addrAndMask := range ipAddr {
			// 分离 IP 和子网掩码
			addr := strings.Split(addrAndMask.String(), "/")[0]
			if strings.Contains(addr, ":") {
				addr = common.DecodeIPv6(addr)
				networkCardInfo[i.Name+" IPv6"] = addr
			} else {
				networkCardInfo[i.Name+" IPv4"] = addr
			}
		}
	}
	return networkCardInfo, nil
}

func GetOwnIP(apiUrl string, enableNetworkCard bool, networkCard string) (acquiredIP string, isIPv6 bool, err error) {
	if enableNetworkCard {
		// 网卡获取
		if networkCard == "" {
			ncr, getErr := NetworkCardRespond()
			err = getErr
			if err != nil {
				return
			}
			err = common.MarshalAndSave(ncr, ConfPath+"/network_card.json")
			if err != nil {
				return
			}
			err = errors.New("请打开 " + ConfPath + "/network_card.json 选择一个网卡填入 " +
				ConfPath + "/client.json 的 \"network_card\"")
			return
		} else {
			ncr, getErr := NetworkCardRespond()
			err = getErr
			if err != nil {
				return
			}
			acquiredIP = ncr[networkCard]
			if acquiredIP == "" {
				err = errors.New("选择了不存在的网卡")
				return
			}
		}
	} else {
		// 远程获取
		if apiUrl == "" {
			apiUrl = common.DefaultAPIServer
		}
		res, getErr := http.Get(apiUrl)
		err = getErr
		if err != nil {
			return
		}
		defer res.Body.Close()
		recvJson, getErr := ioutil.ReadAll(res.Body)
		err = getErr
		if err != nil {
			return
		}
		var ipInfo common.PublicInfo
		err = json.Unmarshal(recvJson, &ipInfo)
		if err != nil {
			return
		}
		acquiredIP = ipInfo.IP
	}
	// 判断 IP 类型
	if strings.Contains(acquiredIP, ":") {
		isIPv6 = true
	} else {
		isIPv6 = false
	}
	return
}

func (conf ClientConf) GetLatestVersion() string {
	res, err := http.Get(conf.APIUrl)
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

func (conf ClientConf) CheckLatestVersion() {
	LatestVersion := conf.GetLatestVersion()
	fmt.Println("当前版本 ", common.LocalVersion)
	fmt.Println("最新版本 ", LatestVersion)
	switch {
	case strings.Contains(LatestVersion, "N/A"):
		fmt.Println("\n需要手动检查更新，请前往 " + common.ProjectUrl + " 查看")
	case common.CompareVersionString(LatestVersion, common.LocalVersion):
		fmt.Println("\n发现新版本，请前往 " + common.ProjectUrl + " 下载")
	}
}
