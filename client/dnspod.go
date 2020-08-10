package client

import (
	"ddns/common"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func DNSPod(dpc DNSPodConf, ipAddr string) (err error) {
	// 获取解析记录
	recordIP, err := dpc.GetParseRecord()
	if err != nil {
		return
	}
	recordType := ""
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
	} else {
		recordType = "A"
	}
	if recordIP == ipAddr {
		err = errors.New("DNSPod 记录的 IP 和当前获取的 IP 一致")
		return
	}
	// 更新解析记录
	err = dpc.UpdateParseRecord(ipAddr, recordType)
	if err != nil {
		return
	}
	return
}

func (dpc DNSPodConf) CheckRespondStatus(jsonObj *simplejson.Json) (err error) {
	statusCode := jsonObj.Get("status").Get("code").MustString()
	if statusCode != "1" {
		err = errors.New("DNSPod: " + statusCode + ": " + jsonObj.Get("status").Get("message").MustString())
		return
	}
	return
}

func (dpc *DNSPodConf) GetParseRecord() (recordIP string, err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordRequestInit()
	recvJson, err := postman("https://dnsapi.cn/Record.List", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = dpc.CheckRespondStatus(jsonObj)
	if err != nil {
		return
	}
	records, err := jsonObj.Get("records").Array()
	if len(records) == 0 {
		err = errors.New("DNSPod: " + dpc.SubDomain + "." + dpc.Domain + " 解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == dpc.SubDomain {
			dpc.RecordId = element["id"].(string)
			recordIP = element["value"].(string)
			dpc.RecordLineId = element["line_id"].(string)
			break
		}
	}
	return
}

func (dpc DNSPodConf) UpdateParseRecord(ipAddr string, recordType string) (err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordModifyRequestInit(ipAddr, recordType)
	recvJson, err := postman("https://dnsapi.cn/Record.Modify", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = dpc.CheckRespondStatus(jsonObj)
	if err != nil {
		return
	}
	return
}

func (dpc DNSPodConf) PublicRequestInit() (pp string) {
	pp = "login_token=" + dpc.Id + "," + dpc.Token +
		"&format=" + "json" +
		"&lang=" + "cn" +
		"&error_on_empty=" + "no"
	return
}

func (dpc DNSPodConf) RecordRequestInit() (rr string) {
	rr = "domain=" + dpc.Domain +
		"&sub_domain=" + dpc.SubDomain
	return
}

func (dpc DNSPodConf) RecordModifyRequestInit(ipAddr string, recordType string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + dpc.RecordId +
		"&sub_domain=" + dpc.SubDomain +
		"&record_type=" + recordType +
		"&record_line_id=" + dpc.RecordLineId +
		"&value=" + ipAddr
	return
}

func postman(url, src string) (dst []byte, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "ddns-client/"+common.LocalVersion+" ()")
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
