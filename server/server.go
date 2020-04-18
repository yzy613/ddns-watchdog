package server

import (
	"ddns/common"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	WorkPath = "/opt/ddns/"
	ConfPath = WorkPath + "conf/"
)

func GetLatestVersion(conf ServerConf) string {
	if !conf.IsRoot {
		res, getErr := http.Get(conf.RootServerAddr)
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
	return common.LocalVersion
}

func CheckLatestVersion(conf ServerConf) {
	if !conf.IsRoot {
		LatestVersion := GetLatestVersion(conf)
		fmt.Println("当前版本 ", common.LocalVersion)
		fmt.Println("最新版本 ", LatestVersion)
		if common.CompareVersionString(LatestVersion, common.LocalVersion) {
			fmt.Println("\n发现新版本，请前往 https://github.com/yzy613/ddns/releases 下载")
		}
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

	// IPv6 转格式 和 ::解压
	switch {
	case strings.Contains(ipAddr, "["):
		ipAddr = strings.Split(ipAddr[1:], "]")[0]
		ipAddr = common.DecodeIPv6(ipAddr)
	case strings.Contains(ipAddr, ":"):
		ipAddr = common.DecodeIPv6(ipAddr)
	}
	return
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func Install() {
	serviceContent := []byte("[Unit]\nDescription=ddns-server Service\nAfter=network.target\n\n[Service]\nType=simple\nExecStart=/opt/ddns/ddns-server\nRestart=on-failure\n\n[Install]\nWantedBy=multi-user.target\n")
	if IsWindows() {
		fmt.Println("Windows 暂不支持安装到系统")
	} else {
		// 复制文件到工作目录
		srcFile, getErr := os.Open("./ddns-server")
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		defer srcFile.Close()
		getErr = os.MkdirAll(WorkPath, 0755)
		if getErr != nil {
			fmt.Println(getErr)
		}
		dstFile, getErr := os.OpenFile(WorkPath+"ddns-server", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0744)
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		defer dstFile.Close()
		buf := make([]byte, 1024)
		for {
			n, getErr := srcFile.Read(buf)
			if getErr != nil {
				if getErr == io.EOF {
					break
				} else {
					fmt.Println(getErr)
					return
				}
			}
			n, getErr = dstFile.Write(buf[:n])
			if getErr != nil {
				fmt.Println(getErr)
				return
			}
		}

		srcFile, getErr = os.Open("./conf/server.json")
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		getErr = os.MkdirAll(ConfPath, 0755)
		if getErr != nil {
			fmt.Println(getErr)
		}
		dstFile, getErr = os.OpenFile(ConfPath+"server.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		for {
			n, getErr := srcFile.Read(buf)
			if getErr != nil {
				if getErr == io.EOF {
					break
				} else {
					fmt.Println(getErr)
					return
				}
			}
			n, getErr = dstFile.Write(buf[:n])
			if getErr != nil {
				fmt.Println(getErr)
				return
			}
		}

		// 注册系统服务
		getErr = ioutil.WriteFile("/etc/systemd/system/ddns-server.service", serviceContent, 0664)
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		fmt.Println("可以使用 systemctl 控制 ddns-server 服务了")
	}
}

func Uninstall() {
	if IsWindows() {
		fmt.Println("Windows 暂不支持安装到系统")
	} else {
		getErr := os.Remove("/etc/systemd/system/ddns-server.service")
		if getErr != nil {
			fmt.Println(getErr)
			return
		}
		fmt.Println("卸载服务成功")
	}
}
