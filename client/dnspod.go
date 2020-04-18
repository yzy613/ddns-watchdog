package client

import (
	"ddns/common"
	"errors"
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

func DNSPod(ipAddr string) (err error) {
	dpc := DNSPodConf{}
	err = common.LoadAndUnmarshal("./conf/dnspod.json", &dpc)
	if err != nil {
		err = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		err = errors.New("请打开配置文件 dnspod.json 填入你的 id, token, domain, sub_domain\n并重新启动")
		return
	}
	if dpc.Id == "" || dpc.Token == "" || dpc.Domain == "" || dpc.SubDomain == "" {
		err = errors.New("请打开配置文件 dnspod.json 填入你的 id, token, domain, sub_domain\n并重新启动")
		return
	}

	recordId, recordType, recordIP, lineId, err := GetParseRecordId(dpc)
	if err != nil {
		return
	}
	if ipAddr[0] == '[' {
		recordType = "AAAA"
	} else {
		recordType = "A"
	}
	if recordId != dpc.RecordId || recordType != dpc.RecordType || lineId != dpc.RecordLineId {
		dpc.RecordId = recordId
		dpc.RecordLineId = lineId
		dpc.RecordType = recordType
		err = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		if err != nil {
			return
		}
	}
	if recordIP == ipAddr {
		err = errors.New("域名解析记录上的 IP 和当前公网 IP 一致")
		return
	}
	err = UpdateParseRecord(dpc, ipAddr)
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

func GetParseRecordId(dpc DNSPodConf) (recordId, recordType, recordIP, lineId string, err error) {
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
			recordType = element["type"].(string)
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
