package client

type ClientConf struct {
	APIUrl            string  `json:"api_url"`
	LatestIP          string  `json:"-"`
	EnableNetworkCard bool    `json:"enable_network_card"`
	NetworkCard       string  `json:"network_card"`
	Services          Service `json:"services"`
	CheckCycle        int     `json:"check_cycle"`
}

type Service struct {
	DNSPod     bool `json:"dnspod"`
	Alidns     bool `json:"alidns"`
	Cloudflare bool `json:"cloudflare"`
}

type DNSPodConf struct {
	Id           string   `json:"id"`
	Token        string   `json:"token"`
	Domain       string   `json:"domain"`
	SubDomain    []string `json:"sub_domain"`
	RecordId     string   `json:"-"`
	RecordLineId string   `json:"-"`
}

type AliyunConf struct {
	AccessKeyId     string   `json:"accesskey_id"`
	AccessKeySecret string   `json:"accesskey_secret"`
	Domain          string   `json:"domain"`
	SubDomain       []string `json:"sub_domain"`
	RecordId        string   `json:"-"`
}

type CloudflareConf struct {
	Email    string   `json:"email"`
	APIKey   string   `json:"api_key"`
	ZoneID   string   `json:"zone_id"`
	Domain   []string `json:"domain"`
	DomainID string   `json:"-"`
}

type CloudflareUpdateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}
