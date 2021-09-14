package common

const (
	LocalVersion      = "1.4.3"
	DefaultAPIUrl     = "https://yzyweb.cn/ddns-watchdog"
	DefaultIPv6APIUrl = "https://yzyweb.cn/ddns-watchdog6"
	ProjectUrl        = "https://github.com/yzy613/ddns-watchdog"
)

type PublicInfo struct {
	IP      string `json:"ip"`
	Version string `json:"latest_version"`
}
