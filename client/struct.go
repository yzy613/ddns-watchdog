package client

type ClientConf struct {
	WebAddr    string `json:"web_addr"`
	LatestIP   string `json:"latest_ip"`
	IsIPv6     bool   `json:"is_ipv6"`
	EnableDdns bool   `json:"enable_ddns"`
	DNSPod     bool   `json:"dnspod"`
}

type DNSPodConf struct {
	Id           string `json:"id"`
	Token        string `json:"token"`
	Domain       string `json:"domain"`
	SubDomain    string `json:"sub_domain"`
	RecordType   string `json:"record_type"`
	RecordId     string `json:"record_id"`
	RecordLineId string `json:"record_line_id"`
}
