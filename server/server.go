package server

import (
	"ddns/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	WorkPath = "/opt/ddns/"
	ConfPath = WorkPath + "conf/"
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
		fmt.Println("当前版本 ", common.LocalVersion)
		fmt.Println("最新版本 ", LatestVersion)
		switch {
		case strings.Contains(LatestVersion, "N/A"):
			fmt.Println("\n需要手动检查更新，请前往 " + common.ProjectAddr + " 查看")
		case common.CompareVersionString(LatestVersion, common.LocalVersion):
			fmt.Println("\n发现新版本，请前往 " + common.ProjectAddr + " 下载")
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
		log.Println("Windows 暂不支持安装到系统")
	} else {
		// 复制文件到工作目录
		getErr := common.CopyFile("./ddns-server", WorkPath+"ddns-server")
		if getErr != nil {
			log.Fatal(getErr)
		}

		// 注册系统服务
		getErr = ioutil.WriteFile("/etc/systemd/system/ddns-server.service", serviceContent, 0664)
		if getErr != nil {
			log.Fatal(getErr)
		}
		log.Println("可以使用 systemctl 控制 ddns-server 服务了")
	}
}

func Uninstall() {
	if IsWindows() {
		log.Println("Windows 暂不支持安装到系统")
	} else {
		err := os.Remove("/etc/systemd/system/ddns-server.service")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("卸载服务成功")
		log.Println("\n若要完全删除，请移步到 /opt/ddns 进行完全删除")
	}
}
