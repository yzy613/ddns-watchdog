package server

import (
	"ddns-watchdog/internal/common"
)

const (
	ServiceConfFileName = "services.json"
)

type service struct {
	DNSPod     dnspod     `json:"dnspod"`
	AliDNS     alidns     `json:"alidns"`
	Cloudflare cloudflare `json:"cloudflare"`
}

type dnspod struct {
	Enable bool   `json:"enable"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

type alidns struct {
	Enable          bool   `json:"enable"`
	AccessKeyId     string `json:"accesskey_id"`
	AccessKeySecret string `json:"accesskey_secret"`
}

type cloudflare struct {
	Enable   bool   `json:"enable"`
	ZoneID   string `json:"zone_id"`
	APIToken string `json:"api_token"`
}

func (conf *service) InitConf() (msg string, err error) {
	*conf = service{
		DNSPod: dnspod{
			ID:    "",
			Token: "",
		},
		AliDNS: alidns{
			AccessKeyId:     "",
			AccessKeySecret: "",
		},
		Cloudflare: cloudflare{
			ZoneID:   "",
			APIToken: "",
		},
	}
	err = common.MarshalAndSave(conf, ConfDirectoryName+"/"+ServiceConfFileName)
	if err != nil {
		return
	}
	msg = "初始化 " + ConfDirectoryName + "/" + ServiceConfFileName
	return
}

func (conf *service) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+ServiceConfFileName, &conf)
	if err != nil {
		return
	}
	err = LoadWhitelist()
	return
}
