package client

import (
	"encoding/json"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func Cloudflare(cfc CloudflareConf, ipAddr string) (err error) {
	// 获取解析记录
	recordIP, err := cfc.GetParseRecord()
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
		err = errors.New("Cloudflare 记录的 IP 和当前获取的 IP 一致")
		return
	}
	// 更新解析记录
	err = cfc.UpdateParseRecord(ipAddr, recordType)
	if err != nil {
		return
	}
	return
}

func (cfc *CloudflareConf) GetParseRecord() (recordIP string, err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records?name=" + cfc.Domain
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
		err = errors.New("Cloudflare: " + cfc.Domain + " 解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == cfc.Domain {
			cfc.DomainID = element["id"].(string)
			recordIP = element["content"].(string)
			break
		}
	}
	return
}

func (cfc CloudflareConf) UpdateParseRecord(ipAddr string, recordType string) (err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records/" + cfc.DomainID
	reqData := CloudflareUpdateRequest{
		Type:    recordType,
		Name:    cfc.Domain,
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
