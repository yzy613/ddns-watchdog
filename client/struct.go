package client

type Service struct {
	DNSPod bool `json:"dnspod"`
	Aliyun bool `json:"aliyun"`
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
	RecordType   string `json:"record_type"`
	RecordLineId string `json:"record_line_id"`
}

type AliyunConf struct {
	AccessKeyId     string `json:"accesskey_id"`
	AccessKeySecret string `json:"accesskey_secret"`
	Domain          string `json:"domain"`
	SubDomain       string `json:"sub_domain"`
	RecordId        string `json:"record_id"`
	RecordType      string `json:"record_type"`
}
