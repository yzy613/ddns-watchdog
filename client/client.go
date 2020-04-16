package client

import (
	"ddns/common"
	"encoding/json"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"io/ioutil"
	"net/http"
)

func DNSPod(dps common.DNSPodSecret) {
	credential := tcommon.NewCredential(dps.SecretId, dps.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "GET"
	cpf.HttpProfile.ReqTimeout = 5
	cpf.SignMethod = "HmacSHA1"

	if credential == nil {
		return
	}
}

func GetOwnIP(webAddr string) (ipAddr string, err error) {
	if webAddr == "" {
		webAddr = "https://yzyweb.cn/ddns"
	}
	res, err := http.Get(webAddr)
	defer res.Body.Close()
	if err != nil {
		return
	}
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var ipInfo common.IpInfoFormat
	err = json.Unmarshal(recvJson, &ipInfo)
	if err != nil {
		return
	}
	ipAddr = ipInfo.Ip
	return
}
