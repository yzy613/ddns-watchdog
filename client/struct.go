package client

type Service struct {
	DNSPod     bool `json:"dnspod"`
	Aliyun     bool `json:"aliyun"`
	Cloudflare bool `json:"cloudflare"`
}

type ClientConf struct {
	APIUrl            string  `json:"api_url"`
	LatestIP          string  `json:"latest_ip"`
	IsIPv6            bool    `json:"is_ipv6"`
	EnableNetworkCard bool    `json:"enable_network_card"`
	NetworkCard       string  `json:"network_card"`
	Services          Service `json:"services"`
}

type DNSPodConf struct {
	Id           string `json:"id"`
	Token        string `json:"token"`
	Domain       string `json:"domain"`
	SubDomain    string `json:"sub_domain"`
	RecordId     string `json:"record_id"`
	RecordLineId string `json:"record_line_id"`
}

type AliyunConf struct {
	AccessKeyId     string `json:"accesskey_id"`
	AccessKeySecret string `json:"accesskey_secret"`
	Domain          string `json:"domain"`
	SubDomain       string `json:"sub_domain"`
	RecordId        string `json:"record_id"`
}

type CloudflareConf struct {
	Email    string `json:"email"`
	APIKey   string `json:"api_key"`
	ZoneID   string `json:"zone_id"`
	Domain   string `json:"domain"`
	DomainID string `json:"domain_id"`
}

type CloudflareUpdateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}
