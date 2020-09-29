package client

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/yzy613/watchdog-ddns/common"
	"io/ioutil"
	"net/http"
	"strings"
)

func DNSPod(dpc DNSPodConf, ipAddr string) (msg []string, err []error) {
	recordType := ""
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
		ipAddr = common.DecodeIPv6(ipAddr)
	} else {
		recordType = "A"
	}

	for _, subDomain := range dpc.SubDomain {
		// 获取解析记录
		recordIP, currentErr := dpc.GetParseRecord(subDomain)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		if recordIP == ipAddr {
			continue
		}
		// 更新解析记录
		currentErr = dpc.UpdateParseRecord(ipAddr, recordType, subDomain)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		msg = append(msg, "DNSPod: " + subDomain + "." + dpc.Domain + " 已更新解析记录 " + ipAddr)
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

func (dpc *DNSPodConf) GetParseRecord(subDomain string) (recordIP string, err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordRequestInit(subDomain)
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
		err = errors.New("DNSPod: " + subDomain + "." + dpc.Domain + " 解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == subDomain {
			dpc.RecordId = element["id"].(string)
			recordIP = element["value"].(string)
			dpc.RecordLineId = element["line_id"].(string)
			break
		}
	}
	return
}

func (dpc DNSPodConf) UpdateParseRecord(ipAddr, recordType, subDomain string) (err error) {
	postContent := dpc.PublicRequestInit()
	postContent = postContent + "&" + dpc.RecordModifyRequestInit(ipAddr, recordType, subDomain)
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

func (dpc DNSPodConf) RecordRequestInit(subDomain string) (rr string) {
	rr = "domain=" + dpc.Domain +
		"&sub_domain=" + subDomain
	return
}

func (dpc DNSPodConf) RecordModifyRequestInit(ipAddr, recordType, subDomain string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + dpc.RecordId +
		"&sub_domain=" + subDomain +
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
