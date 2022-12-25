package server

import (
	"ddns-watchdog/internal/common"
)

const (
	ServiceConfFileName = "services.json"
)

type service struct {
	DNSPod      dnspod      `json:"dnspod"`
	AliDNS      alidns      `json:"alidns"`
	Cloudflare  cloudflare  `json:"cloudflare"`
	HuaweiCloud huaweiCloud `json:"huawei_cloud"`
}

type dnspod struct {
	Enable bool   `json:"enable"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

type alidns struct {
	Enable          bool   `json:"enable"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type cloudflare struct {
	Enable   bool   `json:"enable"`
	ZoneID   string `json:"zone_id"`
	APIToken string `json:"api_token"`
}

type huaweiCloud struct {
	Enable          bool   `json:"enable"`
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
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
