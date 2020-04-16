package client

type DNSPodConf struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	SubDomain string `json:"sub_domain"`
	RecordId  uint   `json:"record_id"`
}

type PublicParameter struct {
	LoginToken   string `json:"login_token"`
	Format       string `json:"format"`
	Lang         string `json:"lang"`
	ErrorOnEmpty string `json:"error_on_empty"`
	// 代理用，个人用户不需要
	UserID string `json:"user_id,omitempty"`
}

// https://www.dnspod.cn/docs/records.html#record-list
type RecordList struct {
	Domain string `json:"domain"`
	Offset uint   `json:"offset,omitempty"`
	// max 3000
	Length       uint   `json:"length,omitempty"`
	SubDomain    string `json:"sub_domain,omitempty"`
	RecordType   string `json:"record_type,omitempty"`
	RecordLine   string `json:"record_line,omitempty"`
	RecordLineId string `json:"record_line_id,omitempty"`
	Keyword      string `json:"keyword,omitempty"`
}

// https://www.dnspod.cn/docs/records.html#record-modify
type RecordModify struct {
	Domain       string `json:"domain"`
	RecordId     uint   `json:"record_id"`
	SubDomain    string `json:"sub_domain,omitempty"`
	RecordType   string `json:"record_type"`
	RecordLine   string `json:"record_line"`
	RecordLineId string `json:"record_line_id,omitempty"`
	Value        string `json:"value"`
	// 1~20
	Mx uint `json:"mx"`
	// 1~604800
	Ttl uint `json:"ttl,omitempty"`
	// enable or disable
	Status string `json:"status,omitempty"`
	// 0~100
	Weight uint `json:"weight,omitempty"`
}
