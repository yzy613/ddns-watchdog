package client

import (
	"ddns/common"
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func PublicRequestInit(dpc DNSPodConf) (pp string) {
	pp = "login_token=" + dpc.Id + "," + dpc.Token +
		"&format=" + "json" +
		"&lang=" + "cn" +
		"&error_on_empty=" + "no"
	return
}

func RecordListInit(dpc DNSPodConf) (ri string) {
	ri = "domain=" + dpc.Domain +
		"&sub_domain=" + dpc.SubDomain
	return
}

func RecordDdnsInit(dpc DNSPodConf, ipAddr string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + "0" +
		"&sub_domain=" + dpc.SubDomain +
		"&record_line=" + "默认" +
		"&record_line_id=" + "" +
		"&value=" + ipAddr
	return
}

func Postman(url, src string) (dst []byte, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "ddns/0.0.1-beta ()")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dst, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

func DNSPod(dpc DNSPodConf, ipAddr string) error {
	postContent := PublicRequestInit(dpc)
	postContent = postContent + "&" + RecordListInit(dpc)
	recvJson, err := Postman("https://dnsapi.cn/Record.List", postContent)
	if err != nil {
		return err
	}
	jsonObj, err := simplejson.NewJson(recvJson)
	// undone
	jsonObj.Get("status")
	return nil
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
