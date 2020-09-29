package client

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/yzy613/watchdog-ddns/common"
	"io/ioutil"
	"net/http"
	"strings"
)

func Cloudflare(cfc CloudflareConf, ipAddr string) (msg []string, err []error) {
	recordType := ""
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
		ipAddr = common.DecodeIPv6(ipAddr)
	} else {
		recordType = "A"
	}

	for _, domain := range cfc.Domain {
		// 获取解析记录
		recordIP, currentErr := cfc.GetParseRecord(domain)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		if recordIP == ipAddr {
			continue
		}
		// 更新解析记录
		currentErr = cfc.UpdateParseRecord(ipAddr, recordType, domain)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		msg = append(msg, "Cloudflare: " + domain + " 已更新解析记录 " + ipAddr)
	}
	return
}

func (cfc *CloudflareConf) GetParseRecord(domain string) (recordIP string, err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records?name=" + domain
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("X-Auth-Email", cfc.Email)
	req.Header.Set("X-Auth-Key", cfc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	if getErr := jsonObj.Get("error").MustString(); getErr != "" {
		err = errors.New(getErr)
		return
	}
	records, err := jsonObj.Get("result").Array()
	if len(records) == 0 {
		err = errors.New("Cloudflare: " + domain + " 解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == domain {
			cfc.DomainID = element["id"].(string)
			recordIP = element["content"].(string)
			break
		}
	}
	return
}

func (cfc CloudflareConf) UpdateParseRecord(ipAddr, recordType, domain string) (err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records/" + cfc.DomainID
	reqData := CloudflareUpdateRequest{
		Type:    recordType,
		Name:    domain,
		Content: ipAddr,
		Ttl:     1,
	}
	reqJson, err := json.Marshal(reqData)
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(reqJson)))
	if err != nil {
		return
	}
	defer req.Body.Close()
	req.Header.Set("X-Auth-Email", cfc.Email)
	req.Header.Set("X-Auth-Key", cfc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	if getErr := jsonObj.Get("error").MustString(); getErr != "" {
		err = errors.New(getErr)
		return
	}
	if !jsonObj.Get("success").MustBool() {
		errorsMsg := ""
		errorsArr, getErr := jsonObj.Get("errors").Array()
		if getErr != nil {
			err = getErr
			return
		}
		for _, value := range errorsArr {
			element := value.(map[string]interface{})
			errCode := element["code"].(json.Number)
			errMsg := element["message"].(string)
			errorsMsg = errorsMsg + "Cloudflare: " + errCode.String() + ": " + errMsg + "\n"
		}
		err = errors.New(errorsMsg)
		return
	}
	return
}
