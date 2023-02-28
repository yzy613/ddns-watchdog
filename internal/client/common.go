package client

import (
	"crypto/tls"
	"net/http"
)

var (
	installPath   = "/etc/systemd/system/" + ProjName + ".service"
	HttpsInsecure = false
)

func getGeneralHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: HttpsInsecure,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
}
