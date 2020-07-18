package client

import (
	"ddns/common"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"strings"
)

func DNSPod(ipAddr string) (err error) {
	dpc := DNSPodConf{}
	err = common.LoadAndUnmarshal("./conf/dnspod.json", &dpc)
	if err != nil {
		err = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		if err != nil {
			return
		}
		err = errors.New("请打开配置文件 ./conf/dnspod.json 填入你的 id, token, domain, sub_domain 并重新启动")
		return
	}
	if dpc.Id == "" || dpc.Token == "" || dpc.Domain == "" || dpc.SubDomain == "" {
		err = errors.New("请打开配置文件 ./conf/dnspod.json 填入你的 id, token, domain, sub_domain 并重新启动")
		return
	}

	recordId, recordType, recordIP, lineId, err := dpc.GetParseRecordId()
	if err != nil {
		return
	}
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
	} else {
		recordType = "A"
	}
	if recordId != dpc.RecordId || recordType != dpc.RecordType || lineId != dpc.RecordLineId {
		dpc.RecordId = recordId
		dpc.RecordType = recordType
		dpc.RecordLineId = lineId
		err = common.MarshalAndSave(dpc, "./conf/dnspod.json")
		if err != nil {
			return
		}
	}
	if recordIP == ipAddr {
		return
	}
	err = dpc.UpdateParseRecord(ipAddr)
	if err != nil {
		return
	}
	return
}

func (dpc DNSPodConf) CheckDNSPodStatus(jsonObj *simplejson.Json) (err error) {
	statusCode := jsonObj.Get("status").Get("code").MustString()
	if statusCode != "1" {
		err = errors.New("DNSPod return " + statusCode + "\n" + jsonObj.Get("status").Get("message").MustString())
		return
	}
	return
}

func (dpc DNSPodConf) GetParseRecordId() (recordId, recordType, recordIP, lineId string, err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordListInit()
	recvJson, err := postman("https://dnsapi.cn/Record.List", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = dpc.CheckDNSPodStatus(jsonObj)
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

func (dpc DNSPodConf) UpdateParseRecord(ipAddr string) (err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordModifyInit(ipAddr)
	recvJson, err := postman("https://dnsapi.cn/Record.Modify", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = dpc.CheckDNSPodStatus(jsonObj)
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

func (dpc DNSPodConf) RecordListInit() (ri string) {
	ri = "domain=" + dpc.Domain +
		"&sub_domain=" + dpc.SubDomain
	return
}

func (dpc DNSPodConf) RecordModifyInit(ipAddr string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + dpc.RecordId +
		"&sub_domain=" + dpc.SubDomain +
		"&record_type=" + dpc.RecordType +
		"&record_line_id=" + dpc.RecordLineId +
		"&value=" + ipAddr
	return
}
