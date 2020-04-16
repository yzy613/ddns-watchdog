package common

type IpInfoFormat struct {
	Ip string `json:"ip"`
}

type ServerConf struct {
	Port string `json:"port"`
}

type ClientConf struct {
	WebAddr    string `json:"web_addr"`
	LastIP     string `json:"last_ip"`
	EnableDdns bool   `json:"enable_ddns"`
	DNSPod     bool   `json:"dnspod"`
}

type DNSPodSecret struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}
