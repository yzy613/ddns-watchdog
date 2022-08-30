package client

var (
	installPath = "/etc/systemd/system/" + RunningName + ".service"
)

type subdomain struct {
	A    string `json:"a"`
	AAAA string `json:"aaaa"`
}
