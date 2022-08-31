package client

import (
	"crypto/tls"
	"net/http"
)

var (
	installPath   = "/etc/systemd/system/" + RunningName + ".service"
	HttpsInsecure = false
)

type subdomain struct {
	A    string `json:"a"`
	AAAA string `json:"aaaa"`
}

func getGeneralHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: HttpsInsecure,
			},
		},
	}
}
