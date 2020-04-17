package client

type DNSPodConf struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	SubDomain string `json:"sub_domain"`
	RecordId  uint   `json:"record_id"`
}
