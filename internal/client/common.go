package client

import (
	"crypto/tls"
	"net/http"
)

var (
	installPath = "/etc/systemd/system/" + ProjName + ".service"
)

func getGeneralHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
}
