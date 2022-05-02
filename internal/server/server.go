package server

import (
	"ddns-watchdog/internal/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	InstallPath       = "/etc/systemd/system/" + RunningName + ".service"
	ConfDirectoryName = "conf"
)

type TLSConf struct {
	Enable   bool   `json:"enable"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type ServerConf struct {
	Port           string  `json:"port"`
	IsRoot         bool    `json:"is_root"`
	RootServerAddr string  `json:"root_server_addr"`
	TLS            TLSConf `json:"tls"`
}

func (conf ServerConf) GetLatestVersion() (str string) {
	if !conf.IsRoot {
		resp, err := http.Get(conf.RootServerAddr)
		if err != nil {
			return "N/A (请检查网络连接)"
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				str = err.Error()
			}
		}(resp.Body)
		recvJson, err := io.ReadAll(resp.Body)
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
	ipAddr = req.Header.Get("X-Forwarded-For")
	if ipAddr != "" && strings.Contains(ipAddr, ",") {
		// 如果只取第零个切片，这行其实可有可无
		//ipAddr = strings.ReplaceAll(ipAddr, " ", "")
		ipAddr = strings.Split(ipAddr, ",")[0]
	}
	if ipAddr == "" {
		ipAddr = req.Header.Get("X-Real-IP")
	}
	if ipAddr == "" {
		// 只保留 ip:port 的 ip
		if strings.Contains(req.RemoteAddr, "[") {
			// IPv6
			ipAddr = strings.Split(req.RemoteAddr[1:], "]:")[0]
		} else {
			// IPv4
			ipAddr = strings.Split(req.RemoteAddr, ":")[0]
		}
	}

	// IPv6 转格式 和 :: 解压
	if strings.Contains(ipAddr, ":") {
		ipAddr = common.DecodeIPv6(ipAddr)
	}
	return
}

func Install() (err error) {
	if common.IsWindows() {
		err = errors.New("windows 暂不支持安装到系统")
	} else {
		// 注册系统服务
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		serviceContent := []byte(
			"[Unit]\n" +
				"Description=" + RunningName + " Service\n" +
				"After=network.target\n\n" +
				"[Service]\n" +
				"Type=simple\n" +
				"WorkingDirectory=" + wd +
				"\nExecStart=" + wd + "/" + RunningName + " -c " + ConfDirectoryName +
				"\nRestart=on-failure\n" +
				"RestartSec=2\n\n" +
				"[Install]\n" +
				"WantedBy=multi-user.target\n")
		err = os.WriteFile(InstallPath, serviceContent, 0664)
		if err != nil {
			return err

		}
		log.Println("可以使用 systemctl 控制 " + RunningName + " 服务了")
	}
	return
}

func Uninstall() (err error) {
	if common.IsWindows() {
		err = errors.New("windows 暂不支持安装到系统")
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Remove(InstallPath)
		if err != nil {
			return err
		}
		log.Println("卸载服务成功")
		log.Println("若要完全删除，请移步到 " + wd + " 和 " + ConfDirectoryName + " 完全删除")
	}
	return
}
