package main

import (
	"ddns-watchdog/internal/client"
	"ddns-watchdog/internal/common"
	"errors"
	"flag"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var (
	installOption   = flag.Bool("I", false, "安装服务并退出")
	uninstallOption = flag.Bool("U", false, "卸载服务并退出")
	enforcement     = flag.Bool("f", false, "强制检查 DNS 解析记录")
	version         = flag.Bool("v", false, "查看当前版本并检查更新后退出")
	initOption      = flag.String("i", "", "有选择地初始化配置文件并退出，可以组合使用 (例 01)\n"+
		"0 -> "+client.ConfFileName+"\n"+
		"1 -> "+client.DNSPodConfFileName+"\n"+
		"2 -> "+client.AliDNSConfFileName+"\n"+
		"3 -> "+client.CloudflareConfFileName)
	confPath             = flag.String("c", "conf", "指定配置文件目录 (目录有空格请放在双引号中间)")
	printNetworkCardInfo = flag.Bool("n", false, "输出网卡信息并退出")
	serviceDebug         = flag.Bool("d", false, "启用 Windows 服务调试模式（配合 Windows 服务使用，请不要在生产环境中使用！）")
)

func main() {
	// 初始化并处理 flag
	exit, err := runFlag()
	if err != nil {
		log.Fatal(err)
	}
	if exit {
		return
	}

	// 加载服务配置
	err = runLoadConf()
	if err != nil {
		log.Fatal(err)
	}

	if common.IsWindowsService {
		runService(*serviceDebug)
	} else {
		if client.Conf.CheckCycleMinutes <= 0 {
			check(elog)
		} else {
			cycle := time.NewTicker(time.Duration(client.Conf.CheckCycleMinutes) * time.Minute)
			for {
				check(elog)
				<-cycle.C
			}
		}
	}
	// 周期循环

}

func runFlag() (exit bool, err error) {
	flag.Parse()
	// 打印网卡信息
	if *printNetworkCardInfo {
		ncr, err2 := client.NetworkCardRespond()
		if err2 != nil {
			err = err2
			return
		}
		var arr []string
		for key := range ncr {
			arr = append(arr, key)
		}
		sort.Strings(arr)
		for _, key := range arr {
			fmt.Printf("%v\n\t%v\n", key, ncr[key])
		}
		exit = true
		return
	}

	// 加载自定义配置文件目录
	if *confPath != "" {
		client.ConfDirectoryName = common.FormatDirectoryPath(*confPath)
	}

	// 有选择地初始化配置文件
	if *initOption != "" {
		for _, event := range *initOption {
			err = runInitConf(string(event))
			if err != nil {
				return
			}
		}
		exit = true
		return
	}

	// 加载客户端配置
	// 不得不放在这个地方，因为有下面的检查版本和安装 / 卸载服务
	err = client.Conf.LoadConf()
	if err != nil {
		return
	}

	// 检查版本
	if *version {
		client.Conf.CheckLatestVersion()
		exit = true
		return
	}

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err = client.Install(*confPath)
		if err != nil {
			return
		}
		exit = true
		return
	case *uninstallOption:
		err = client.Uninstall()
		if err != nil {
			return
		}
		exit = true
		return
	}
	return
}

func runInitConf(event string) error {
	switch event {
	case "0":
		msg, err := client.Conf.InitConf()
		if err != nil {
			return err
		}
		log.Println(msg)
	case "1":
		msg, err := client.Dpc.InitConf()
		if err != nil {
			return err
		}
		log.Println(msg)
	case "2":
		msg, err := client.Adc.InitConf()
		if err != nil {
			return err
		}
		log.Println(msg)
	case "3":
		msg, err := client.Cfc.InitConf()
		if err != nil {
			return err
		}
		log.Println(msg)
	default:
		err := errors.New("你初始化了一个寂寞")
		return err
	}
	return nil
}

func runLoadConf() (err error) {
	if client.Conf.Services.DNSPod {
		err = client.Dpc.LoadConf()
		if err != nil {
			return
		}
	}
	if client.Conf.Services.AliDNS {
		err = client.Adc.LoadConf()
		if err != nil {
			return
		}
	}
	if client.Conf.Services.Cloudflare {
		err = client.Cfc.LoadConf()
		if err != nil {
			return
		}
	}
	return
}

func check(elog debug.Log) {
	// 获取 IP
	ipv4, ipv6, err := client.GetOwnIP(client.Conf.Enable, client.Conf.APIUrl, client.Conf.NetworkCard)
	if err != nil {
		if common.IsWindowsService {
			elog.Error(101, err.Error())
		} else {
			log.Println(err)
		}
		return
	}

	// 进入更新流程
	if ipv4 != client.Conf.LatestIPv4 || ipv6 != client.Conf.LatestIPv6 || *enforcement {
		if ipv4 != client.Conf.LatestIPv4 {
			client.Conf.LatestIPv4 = ipv4
		}
		if ipv6 != client.Conf.LatestIPv6 {
			client.Conf.LatestIPv6 = ipv6
		}
		wg := sync.WaitGroup{}
		if client.Conf.Services.DNSPod {
			wg.Add(1)
			go asyncServiceInterface(ipv4, ipv6, client.Dpc.Run, &wg, elog)
		}
		if client.Conf.Services.AliDNS {
			wg.Add(1)
			go asyncServiceInterface(ipv4, ipv6, client.Adc.Run, &wg, elog)
		}
		if client.Conf.Services.Cloudflare {
			wg.Add(1)
			go asyncServiceInterface(ipv4, ipv6, client.Cfc.Run, &wg, elog)
		}
		wg.Wait()
	}
}

func asyncServiceInterface(ipv4, ipv6 string, callback client.AsyncServiceCallback, wg *sync.WaitGroup, elog debug.Log) {
	defer wg.Done()
	msg, err := callback(client.Conf.Enable, ipv4, ipv6)
	for _, row := range err {
		if common.IsWindowsService {
			elog.Error(102, row.Error())
		} else {
			log.Println(row)
		}
	}
	for _, row := range msg {
		if common.IsWindowsService {
			elog.Info(100, row)
		} else {
			log.Println(row)
		}
	}
}

// Windows 服务
type WindowsService struct{}

var elog debug.Log

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	changes <- svc.Status{State: svc.StartPending}
	var tick = time.NewTicker(time.Duration(client.Conf.CheckCycleMinutes) * time.Minute)
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptPauseAndContinue | svc.AcceptShutdown | svc.AcceptStop}
	elog.Info(2, fmt.Sprintf("服务 %s 启动成功！", client.RunningName))
	check(elog)
	elog.Info(3, "动态域名解析更新完成！")
	if client.Conf.CheckCycleMinutes <= 0 {
		changes <- svc.Status{State: svc.StopPending}
		tick.Stop()
		changes <- svc.Status{State: svc.Stopped}
		return
	}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				elog.Info(6, fmt.Sprintf("服务 %s 正在停止……", client.RunningName))
				tick.Stop()
				break loop
			case svc.Pause:
				tick.Stop()
				changes <- svc.Status{State: svc.Paused, Accepts: svc.AcceptPauseAndContinue | svc.AcceptShutdown | svc.AcceptStop}
				elog.Info(4, fmt.Sprintf("服务 %s 已暂停！", client.RunningName))
			case svc.Continue:
				tick = time.NewTicker(time.Duration(client.Conf.CheckCycleMinutes) * time.Minute)
				changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptPauseAndContinue | svc.AcceptShutdown | svc.AcceptStop}
				elog.Info(5, fmt.Sprintf("服务 %s 已恢复！", client.RunningName))
				check(elog)
				elog.Info(3, "动态域名解析更新成功！")
			default:
				elog.Error(9, fmt.Sprintf("无法识别的控制命令 #%d", c))
			}
		case <-tick.C:
			elog.Info(3, "动态域名解析更新成功！")
			check(elog)
		}
	}
	changes <- svc.Status{State: svc.Stopped}
	return
}

func runService(isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(client.RunningName)
	} else {
		elog, err = eventlog.Open(client.RunningName)
		if err != nil {
			return
		}
	}

	defer elog.Close()
	elog.Info(1, fmt.Sprintf("服务 %s 正在启动中……", client.RunningName))
	run := svc.Run
	if isDebug {
		elog.Warning(50, fmt.Sprintf("服务 %s 将以调试模式运行！", client.RunningName))
		run = debug.Run
	}
	err = run(client.RunningName, &WindowsService{})
	if err != nil {
		elog.Error(8, fmt.Sprintf("服务 %s 启动失败，错误信息: %v", client.RunningName, err))
		return
	}
	elog.Info(7, fmt.Sprintf("服务 %s 已停止！", client.RunningName))
}
