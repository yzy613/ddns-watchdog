package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"watchdog-ddns/common"
)

var (
	RunningName            = "watchdog-ddns-client"
	RunningPath            = common.GetRunningPath()
	InstallPath            = "/etc/systemd/system/" + RunningName + ".service"
	ConfPath               = RunningPath + "conf/"
	ConfFileName           = "client.json"
	DNSPodConfFileName     = "dnspod.json"
	AliDNSConfFileName     = "alidns.json"
	CloudflareConfFileName = "cloudflare.json"
	NetworkCardFileName    = "network_card.json"
)

func Install() (err error) {
	if common.IsWindows() {
		log.Println("Windows 暂不支持安装到系统")
	} else {
		// 注册系统服务
		serviceContent := []byte(
			"[Unit]\n" +
				"Description=" + RunningName + " Service\n" +
				"After=network.target\n\n" +
				"[Service]\n" +
				"Type=simple\n" +
				"ExecStart=" + RunningPath + RunningName + " -conf_path " + ConfPath +
				"\nRestart=on-failure\n" +
				"RestartSec=2\n\n" +
				"[Install]\n" +
				"WantedBy=multi-user.target\n")
		err = ioutil.WriteFile(InstallPath, serviceContent, 0664)
		if err != nil {
			return
		}
		log.Println("可以使用 systemctl 控制 " + RunningName + " 服务了")
	}
	return
}

func Uninstall() (err error) {
	if common.IsWindows() {
		log.Println("Windows 暂不支持安装到系统")
	} else {
		err = os.Remove(InstallPath)
		if err != nil {
			return
		}
		log.Println("卸载服务成功")
		log.Println("若要完全删除，请移步到 " + RunningPath + " 和 " + ConfPath + " 完全删除")
	}
	return
}

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
		for j, addrAndMask := range ipAddr {
			// 分离 IP 和子网掩码
			addr := strings.Split(addrAndMask.String(), "/")[0]
			if strings.Contains(addr, ":") {
				addr = common.DecodeIPv6(addr)
			}
			networkCardInfo[i.Name+" "+strconv.Itoa(j)] = addr
		}
	}
	return networkCardInfo, nil
}

func GetOwnIP(apiUrl string, enableNetworkCard bool, networkCard string) (acquiredIP string, err error) {
	if enableNetworkCard {
		// 网卡获取
		if networkCard == "" {
			ncr, getErr := NetworkCardRespond()
			err = getErr
			if err != nil {
				return
			}
			err = common.MarshalAndSave(ncr, ConfPath+NetworkCardFileName)
			if err != nil {
				return
			}
			err = errors.New("请打开 " + ConfPath + NetworkCardFileName + " 选择一个网卡填入 " +
				ConfPath + ConfFileName + " 的 network_card")
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
	common.VersionTips(LatestVersion)
}
