package client

import (
	"ddns-watchdog/internal/common"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	installPath       = "/etc/systemd/system/" + RunningName + ".service"
	ConfDirectoryName = "conf"
	Conf              = clientConf{}
	Dpc               = dnspodConf{}
	Adc               = aliDNSConf{}
	Cfc               = cloudflareConf{}
)

type subdomain struct {
	A    string `json:"a"`
	AAAA string `json:"aaaa"`
}

// AsyncServiceCallback 异步服务回调函数类型
type AsyncServiceCallback func(enabledServices enable, ipv4, ipv6 string) (msg []string, errs []error)

func exePath(confPath string) (path string, err error) { // 获取可执行文件路径
	confPath, err = filepath.Abs(confPath)
	if err != nil {
		return "", err
	}
	prog := os.Args[0]
	path, err = filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(path)
	if err == nil {
		if !fi.Mode().IsDir() {
			return `"` + path + `" -c "` + confPath + `"`, nil
		}
		err = errors.New(path + " 是一个目录！")
	}
	if filepath.Ext(path) == "" {
		path += ".exe"
		fi, err := os.Stat(path)
		if err == nil {
			if !fi.Mode().IsDir() {
				return `"` + path + `" -c "` + confPath + `"`, nil
			}
			return path, errors.New(path + " 是一个目录！")
		}
	}
	path = `"` + path + `" -c "` + confPath + `"`
	return
}

// 注册系统服务
func Install(confPath string) (err error) {
	if Conf.CheckCycleMinutes == 0 {
		err = errors.New("设置一下 " + ConfDirectoryName + "/" + ConfFileName + " 的 check_cycle_minutes 吧")
		return
	}
	if common.IsWindows {
		exepath, err := exePath(confPath)
		if err != nil {
			log.Fatalf("exePath 生成失败！错误信息：%v", err)
		}
		m, err := mgr.Connect()
		if err != nil {
			log.Fatalf("Windows 服务管理器连接失败！错误信息：%v", err)
		}
		defer m.Disconnect()
		s, err := m.OpenService(RunningName)
		if err == nil {
			s.Close()
			log.Fatalf("服务 %s 已存在！", RunningName)
		}
		config := mgr.Config{
			DisplayName:      "DDNS-Watchdog 动态域名解析客户端",
			Description:      "DDNS-Watchdog 动态域名解析客户端服务，用于在动态域名解析（DDNS）中自动化更新解析记录。",
			ServiceType:      windows.SERVICE_WIN32_OWN_PROCESS,
			StartType:        windows.SERVICE_AUTO_START,
			ErrorControl:     windows.SERVICE_ERROR_NORMAL,
			ServiceStartName: "NT AUTHORITY\\NetworkService",
			DelayedAutoStart: true,
		}
		s, err = m.CreateService(RunningName, exepath, config)
		if err != nil {
			return err
		}
		defer s.Close()
		recoveryActions := []mgr.RecoveryAction{
			{
				Type:  windows.SC_ACTION_RESTART,
				Delay: (5 * time.Minute),
			},
			{
				Type:  windows.SC_ACTION_RESTART,
				Delay: (5 * time.Minute),
			},
			{
				Type:  windows.SC_ACTION_NONE,
				Delay: (5 * time.Minute),
			},
		}
		err = s.SetRecoveryActions(recoveryActions, 2*86400)
		if err != nil {
			s.Delete()
			log.Fatalf("设置错误处理程序时发生错误：%v", err)
		}
		err = eventlog.InstallAsEventCreate(RunningName, eventlog.Error|eventlog.Warning|eventlog.Info)
		if err != nil {
			s.Delete()
			log.Fatalf("调用 SetupEventLogSource() 时发生错误：%v", err)
		}
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\`+RunningName, registry.ALL_ACCESS)
		if err != nil {
			log.Fatalf("服务被不完全安装，问题出现在二次写入注册表过程中，请尝试重新安装服务！信息：%v", err)
		}
		err = key.SetStringValue(`ImagePath`, exepath)
		if err != nil {
			log.Fatalf("服务被不完全安装，问题出现在二次写入注册表过程中，请尝试重新安装服务！信息：%v", err)
		}
		log.Printf("服务 %s 安装成功！", RunningName)
		return nil
	} else {
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
		err = os.WriteFile(installPath, serviceContent, 0664)
		if err != nil {
			return err
		}
		log.Println("可以使用 systemctl 控制 " + RunningName + " 服务了")
	}
	return
}

func Uninstall() (err error) {
	if common.IsWindows {
		m, err := mgr.Connect()
		if err != nil {
			log.Fatalf("在连接至 Windows 服务管理器时发生错误！详细信息：%v", err)
		}
		defer m.Disconnect()
		s, err := m.OpenService(RunningName)
		if err != nil {
			err = errors.New("此程序尚未被安装为服务！")
			return err
		}
		defer s.Close()
		err = s.Delete()
		if err != nil {
			log.Fatalf("删除服务时发生错误：%v", err)
		}
		err = eventlog.Remove(RunningName)
		if err != nil {
			log.Fatalf("调用 RemoveEventLogSource() 时发生错误：%v", err)
		}
		log.Printf("服务 %s 删除成功！", RunningName)
		return nil
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Remove(installPath)
		if err != nil {
			return err
		}
		log.Println("卸载服务成功")
		log.Println("若要完全删除，请移步到 " + wd + " 和 " + ConfDirectoryName + " 完全删除")
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
		ipAddr, err2 := i.Addrs()
		if err2 != nil {
			return nil, err2
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

func GetOwnIP(enabled enable, apiUrl apiUrl, nc networkCard) (ipv4, ipv6 string, err error) {
	ncr := make(map[string]string)
	// 若需网卡信息，则获取网卡信息并提供给用户
	if enabled.NetworkCard && nc.IPv4 == "" && nc.IPv6 == "" {
		ncr, err = NetworkCardRespond()
		if err != nil {
			return
		}
		err = common.MarshalAndSave(ncr, ConfDirectoryName+"/"+NetworkCardFileName)
		if err != nil {
			return
		}
		err = errors.New("请打开 " + ConfDirectoryName + "/" + NetworkCardFileName + " 选择网卡填入 " +
			ConfDirectoryName + "/" + ConfFileName + " 的 network_card")
		return
	}

	// 若需网卡信息，则获取网卡信息
	if enabled.NetworkCard && (nc.IPv4 != "" || nc.IPv6 != "") {
		ncr, err = NetworkCardRespond()
		if err != nil {
			return
		}
	}

	// 启用 IPv4
	if enabled.IPv4 {
		// 启用网卡 IPv4
		if enabled.NetworkCard && nc.IPv4 != "" {
			ipv4 = ncr[nc.IPv4]
			if ipv4 == "" {
				err = errors.New("IPv4 选择了不存在的网卡")
				return
			}
		} else {
			// 使用 API 获取 IPv4
			if apiUrl.IPv4 == "" {
				apiUrl.IPv4 = common.DefaultAPIUrl
			}
			resp, err2 := http.Get(apiUrl.IPv4)
			if err2 != nil {
				err = err2
				return
			}
			defer func(Body io.ReadCloser) {
				t := Body.Close()
				if t != nil {
					err = t
				}
			}(resp.Body)
			recvJson, err2 := io.ReadAll(resp.Body)
			if err2 != nil {
				err = err2
				return
			}
			var ipInfo common.PublicInfo
			err = json.Unmarshal(recvJson, &ipInfo)
			if err != nil {
				return
			}
			ipv4 = ipInfo.IP
		}
		if strings.Contains(ipv4, ":") {
			err = errors.New("获取到的 IPv4 格式错误，意外获取到了 " + ipv4)
			return
		}
	}

	// 启用 IPv6
	if enabled.IPv6 {
		// 启用网卡 IPv6
		if enabled.NetworkCard && nc.IPv6 != "" {
			ipv6 = ncr[nc.IPv6]
			if ipv6 == "" {
				err = errors.New("IPv6 选择了不存在的网卡")
				return
			}
		} else {
			// 使用 API 获取 IPv4
			if apiUrl.IPv6 == "" {
				apiUrl.IPv6 = common.DefaultIPv6APIUrl
			}
			resp, err2 := http.Get(apiUrl.IPv6)
			if err2 != nil {
				err = err2
				return
			}
			defer func(Body io.ReadCloser) {
				t := Body.Close()
				if t != nil {
					err = t
				}
			}(resp.Body)
			recvJson, err2 := io.ReadAll(resp.Body)
			if err2 != nil {
				err = err2
				return
			}
			var ipInfo common.PublicInfo
			err = json.Unmarshal(recvJson, &ipInfo)
			if err != nil {
				return
			}
			ipv6 = ipInfo.IP
		}
		if strings.Contains(ipv6, ":") {
			ipv6 = common.DecodeIPv6(ipv6)
		} else {
			err = errors.New("获取到的 IPv6 格式错误，意外获取到了 " + ipv6)
			return
		}
	}
	return
}
