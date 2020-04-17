package client

import (
	"ddns/common"
	"encoding/json"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetOwnIP(webAddr string) (ipAddr string, isIPv6 bool, err error) {
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
	var ipInfo common.PublicInfo
	err = json.Unmarshal(recvJson, &ipInfo)
	if err != nil {
		return
	}
	ipAddr = ipInfo.IP
	if ipAddr[0] == '[' {
		isIPv6 = true
	} else {
		isIPv6 = false
	}
	return
}

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

func RecordModifyInit(dpc DNSPodConf, ipAddr string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + dpc.RecordId +
		"&sub_domain=" + dpc.SubDomain +
		"&record_type=" + dpc.RecordType +
		"&record_line_id=" + dpc.RecordLineId +
		"&value=" + ipAddr
	return
}

func Postman(url, src string) (dst []byte, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "ddns-client/0.1.0 ()")
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

func DNSPod(dpc *DNSPodConf, ipAddr string) (err error) {
	recordId, recordIP, lineId, err := GetParseRecordId(*dpc)
	if err != nil {
		return
	}
	if recordId != dpc.RecordId || lineId != dpc.RecordLineId {
		dpc.RecordId = recordId
		dpc.RecordLineId = lineId
		err = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		if err != nil {
			return
		}
	}
	if recordIP == ipAddr {
		err = errors.New("域名解析记录上的 IP 和当前公网 IP 一致")
		return
	}
	err = UpdateParseRecord(*dpc, ipAddr)
	if err != nil {
		return
	}
	return
}

func CheckStatus(jsonObj *simplejson.Json) (err error) {
	statusCode := jsonObj.Get("status").Get("code").MustString()
	if statusCode != "1" {
		err = errors.New("DNSPod return " + statusCode + "\n" + jsonObj.Get("status").Get("message").MustString())
		return
	}
	return
}

func GetParseRecordId(dpc DNSPodConf) (recordId string, recordIP string, lineId string, err error) {
	postContent := PublicRequestInit(dpc)
	postContent = postContent + "&" + RecordListInit(dpc)
	recvJson, err := Postman("https://dnsapi.cn/Record.List", postContent)
	if err != nil {
		return
	}
	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = CheckStatus(jsonObj)
	if err != nil {
		return
	}
	records, err := jsonObj.Get("records").Array()
	if len(records) == 0 {
		err = errors.New("解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == dpc.SubDomain {
			recordId = element["id"].(string)
			recordIP = element["value"].(string)
			lineId = element["line_id"].(string)
			break
		}
	}
	return
}

func UpdateParseRecord(dpc DNSPodConf, ipAddr string) (err error) {
	postContent := PublicRequestInit(dpc)
	postContent = postContent + "&" + RecordModifyInit(dpc, ipAddr)
	recvJson, err := Postman("https://dnsapi.cn/Record.Modify", postContent)
	if err != nil {
		return
	}
	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = CheckStatus(jsonObj)
	if err != nil {
		return
	}
	return
}
