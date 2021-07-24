package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yzy613/ddns-watchdog/client"
	"github.com/yzy613/ddns-watchdog/common"

	"github.com/yzy613/ddns-watchdog/golang.org/x/sys/windows/svc"
	"github.com/yzy613/ddns-watchdog/golang.org/x/sys/windows/svc/debug"
	"github.com/yzy613/ddns-watchdog/golang.org/x/sys/windows/svc/eventlog"
)

var (
	installOption   = flag.Bool("I", false, "安装服务")
	uninstallOption = flag.Bool("U", false, "卸载服务")
	enforcement     = flag.Bool("f", false, "强制检查 DNS 解析记录")
	version         = flag.Bool("v", false, "查看当前版本并检查更新")
	runAsService    = flag.Bool("s", false, "以 Windows 服务模式运行（请不要自行添加此参数！）")
	debugMode       = flag.Bool("d", false, "在以 Windows 服务模式运行时开启调试模式（请不要自行添加此参数！）")
	initOption      = flag.String("i", "", "有选择地初始化配置文件，可以组合使用 (例 01)\n"+
		"0 -> "+client.ConfFileName+"\n"+
		"1 -> "+client.DNSPodConfFileName+"\n"+
		"2 -> "+client.AliDNSConfFileName+"\n"+
		"3 -> "+client.CloudflareConfFileName)
	confPath = flag.String("c", "", "指定配置文件路径 (最好是绝对路径)(路径有空格请放在双引号中间)")
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

	if *runAsService { // Windows 服务模式
		if common.IsWindows() { // 是 Windows 系统
			runService(client.WindowsServiceName, *debugMode)
		} else { // 不是 Windows 系统
			log.Fatal("请不要手动传入 -s 参数！")
		}
	} else {
		waitCheckDone := make(chan bool, 1)
		if client.Conf.CheckCycleMinutes <= 0 {
			go asyncCheck(waitCheckDone)
			<-waitCheckDone
		} else {
			cycle := time.NewTicker(time.Duration(client.Conf.CheckCycleMinutes) * time.Minute)
			for {
				go asyncCheck(waitCheckDone)
				<-waitCheckDone
				<-cycle.C
			}
		}
	}
	// 周期循环

}

func runFlag() (exit bool, err error) {
	flag.Parse()
	// 加载自定义配置文件路径
	if *confPath != "" {
		tempStr := *confPath
		if tempStr[len(tempStr)-1:] != "/" {
			tempStr = tempStr + "/"
		}
		client.ConfPath = tempStr
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

	// 安装 / 卸载服务
	switch {
	case *installOption:
		err = client.Install()
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

	// 加载客户端配置
	// 不得不放在这个地方，因为有下面的检查版本
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
		err = client.Dpc.LoadCOnf()
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

func asyncCheck(done chan bool) {
	// 获取 IP
	ipv4, ipv6, err := client.GetOwnIP(client.Conf.Enable, client.Conf.APIUrl, client.Conf.NetworkCard)
	if err != nil {
		log.Println(err)
		done <- true
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
		servicesCount := 0
		if client.Conf.Services.DNSPod {
			servicesCount++
		}
		if client.Conf.Services.AliDNS {
			servicesCount++
		}
		if client.Conf.Services.Cloudflare {
			servicesCount++
		}
		waitServicesDone := make(chan bool, servicesCount)
		if client.Conf.Services.DNSPod {
			go asyncDNSPod(ipv4, ipv6, waitServicesDone)
		}
		if client.Conf.Services.AliDNS {
			go asyncAliDNS(ipv4, ipv6, waitServicesDone)
		}
		if client.Conf.Services.Cloudflare {
			go asyncCloudflare(ipv4, ipv6, waitServicesDone)
		}
		for i := 0; i < servicesCount; i++ {
			<-waitServicesDone
		}
	}
	done <- true
}

func asyncDNSPod(ipv4, ipv6 string, done chan bool) {
	msg, err := client.Dpc.Run(client.Conf.Enable, ipv4, ipv6)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}

func asyncAliDNS(ipv4, ipv6 string, done chan bool) {
	msg, err := client.Adc.Run(client.Conf.Enable, ipv4, ipv6)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}

func asyncCloudflare(ipv4, ipv6 string, done chan bool) {
	msg, err := client.Cfc.Run(client.Conf.Enable, ipv4, ipv6)
	for _, row := range err {
		log.Println(row)
	}
	for _, row := range msg {
		log.Println(row)
	}
	done <- true
}

// Windows 服务
var elog debug.Log

type WindowsService struct{}

func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptPauseAndContinue | svc.AcceptShutdown | svc.AcceptStop
	changes <- svc.Status{State: svc.StartPending}
	runTick := time.NewTicker(time.Duration(client.Conf.CheckCycleMinutes) * time.Minute)
	pauseTick := time.NewTicker(17280 * time.Hour)
	tick := runTick
	waitCheckDone := make(chan bool, 1)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	elog.Info(1, fmt.Sprintf("服务 %s 启动成功！", client.WindowsServiceName))
	go asyncCheck(waitCheckDone)
	<-waitCheckDone
	elog.Info(3, "动态域名解析更新成功！")
	if client.Conf.CheckCycleMinutes <= 0 {
		changes <- svc.Status{State: svc.StopPending}
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
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = pauseTick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = runTick
				go asyncCheck(waitCheckDone)
				<-waitCheckDone
				elog.Info(3, "动态域名解析更新成功！")
			default:
				elog.Error(2, fmt.Sprintf("无法识别的控制命令 #%d", c))
			}
		case <-tick.C:
			elog.Info(3, "动态域名解析更新成功！")
			go asyncCheck(waitCheckDone)
			<-waitCheckDone
		}
	}
	changes <- svc.Status{State: svc.Stopped}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("服务 %s 正在启动中……", name))
	run := svc.Run
	if isDebug {
		elog.Warning(1, fmt.Sprintf("服务 %s 将以调试模式运行！", name))
		run = debug.Run
	}
	err = run(name, &WindowsService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("服务 %s 启动失败: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("服务 %s 已停止！", name))
}
