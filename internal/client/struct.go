package client

import (
	"ddns-watchdog/internal/common"
)

const (
	RunningName            = "ddns-watchdog-client"
	ConfFileName           = "client.json"
	DNSPodConfFileName     = "dnspod.json"
	AliDNSConfFileName     = "alidns.json"
	CloudflareConfFileName = "cloudflare.json"
	NetworkCardFileName    = "network_card.json"
)

var (
	RunningPath = common.GetRunningPath()
	InstallPath = "/etc/systemd/system/" + RunningName + ".service"
	ConfPath    = RunningPath + "conf/"
	Conf        = clientConf{}
	Dpc         = dnspodConf{}
	Adc         = aliDNSConf{}
	Cfc         = cloudflareConf{}
)

type clientConf struct {
	APIUrl            apiUrl      `json:"api_url"`
	Enable            enable      `json:"enable"`
	NetworkCard       networkCard `json:"network_card"`
	Services          service     `json:"services"`
	CheckCycleMinutes int         `json:"check_cycle_minutes"`
	LatestIPv4        string      `json:"-"`
	LatestIPv6        string      `json:"-"`
}

type apiUrl struct {
	IPv4    string `json:"ipv4"`
	IPv6    string `json:"ipv6"`
	Version string `json:"version"`
}

type enable struct {
	IPv4        bool `json:"ipv4"`
	IPv6        bool `json:"ipv6"`
	NetworkCard bool `json:"network_card"`
}

type networkCard struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

type service struct {
	DNSPod     bool `json:"dnspod"`
	AliDNS     bool `json:"alidns"`
	Cloudflare bool `json:"cloudflare"`
}

type subdomain struct {
	A    string `json:"a"`
	AAAA string `json:"aaaa"`
}

type dnspodConf struct {
	Id           string    `json:"id"`
	Token        string    `json:"token"`
	Domain       string    `json:"domain"`
	SubDomain    subdomain `json:"sub_domain"`
	RecordId     string    `json:"-"`
	RecordLineId string    `json:"-"`
}

type aliDNSConf struct {
	AccessKeyId     string    `json:"accesskey_id"`
	AccessKeySecret string    `json:"accesskey_secret"`
	Domain          string    `json:"domain"`
	SubDomain       subdomain `json:"sub_domain"`
	RecordId        string    `json:"-"`
}

type cloudflareConf struct {
	ZoneID   string    `json:"zone_id"`
	APIToken string    `json:"api_token"`
	Domain   subdomain `json:"domain"`
	DomainID string    `json:"-"`
}

type cloudflareUpdateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

// AsyncServerCallback 异步服务回调函数类型
type AsyncServerCallback func(enabledServices enable, ipv4, ipv6 string) (msg []string, errs []error)