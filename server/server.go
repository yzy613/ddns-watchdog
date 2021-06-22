package server

import (
	"encoding/json"
	"fmt"
	"github.com/yzy613/ddns-watchdog/common"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	RunningName  = "ddns-watchdog-server"
	ConfFileName = "server.json"
)

var (
	RunningPath = common.GetRunningPath()
	InstallPath = "/etc/systemd/system/" + RunningName + ".service"
	ConfPath    = RunningPath + "conf/"
)

func (conf ServerConf) GetLatestVersion() string {
	if !conf.IsRoot {
		res, err := http.Get(conf.RootServerAddr)
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
	return common.LocalVersion
}

func (conf ServerConf) CheckLatestVersion() {
	if !conf.IsRoot {
		LatestVersion := conf.GetLatestVersion()
		common.VersionTips(LatestVersion)
	} else {
		fmt.Println("本机是根服务器")
		fmt.Println("当前版本 ", common.LocalVersion)
	}
}

func GetClientIP(req *http.Request) (ipAddr string) {
	ipAddr = req.Header.Get("X-Real-IP")
	if ipAddr == "" {
		ipAddr = req.Header.Get("X-Forwarded-For")
	}
	if ipAddr == "" {
		// 把 port 从 ip:port 分离
		if strings.Contains(req.RemoteAddr, "[") {
			// IPv6
			ipAddr = strings.Split(req.RemoteAddr, "]:")[0]
			ipAddr = ipAddr + "]"
		} else {
			// IPv4
			ipAddr = strings.Split(req.RemoteAddr, ":")[0]
		}
	}

	// IPv6 转格式 和 :: 解压
	switch {
	case strings.Contains(ipAddr, "["):
		ipAddr = strings.Split(ipAddr[1:], "]")[0]
		ipAddr = common.DecodeIPv6(ipAddr)
	case strings.Contains(ipAddr, ":"):
		ipAddr = common.DecodeIPv6(ipAddr)
	}
	return
}

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
				"ExecStart=" + RunningPath + RunningName + " -c " + ConfPath +
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
