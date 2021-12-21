package client

import (
	"ddns-watchdog/internal/common"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func (conf *clientConf) InitConf() (msg string, err error) {
	*conf = clientConf{}
	conf.APIUrl.IPv4 = common.DefaultAPIUrl
	conf.APIUrl.IPv6 = common.DefaultIPv6APIUrl
	conf.APIUrl.Version = common.DefaultAPIUrl
	conf.CheckCycleMinutes = 0
	err = common.MarshalAndSave(conf, ConfPath+ConfFileName)
	msg = "初始化 " + ConfPath + ConfFileName
	return
}

func (conf *clientConf) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfPath+ConfFileName, &conf)
	// 检查启用 IP 类型
	if !conf.Enable.IPv4 && !conf.Enable.IPv6 {
		err = errors.New("请打开客户端配置文件 " + ConfPath + ConfFileName + " 启用需要使用的 IP 类型并重新启动")
		return
	}
	// 检查启用服务
	if !conf.Services.DNSPod && !conf.Services.AliDNS && !conf.Services.Cloudflare {
		err = errors.New("请打开客户端配置文件 " + ConfPath + ConfFileName + " 启用需要使用的服务并重新启动")
		return
	}
	return
}

func (conf clientConf) GetLatestVersion() (str string) {
	res, err := http.Get(conf.APIUrl.Version)
	if err != nil {
		return "N/A (请检查网络连接)"
	}
	defer func(Body io.ReadCloser) {
		t := Body.Close()
		if t != nil {
			str = t.Error()
		}
	}(res.Body)
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

func (conf *clientConf) CheckLatestVersion() {
	if conf.APIUrl.Version == "" {
		conf.APIUrl.Version = common.DefaultAPIUrl
	}
	common.VersionTips(conf.GetLatestVersion())
}
