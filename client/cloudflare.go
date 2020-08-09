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

func Cloudflare(ipAddr string) (err error) {
	cfc := CloudflareConf{}
	// 获取配置
	err = common.LoadAndUnmarshal(ConfPath+"/cloudflare.json", &cfc)
	if err != nil {
		return
	}

	if cfc.Email == "" || cfc.APIKey == "" || cfc.ZoneID == "" || cfc.Domain == "" {
		err = errors.New("请打开配置文件 " + ConfPath + "/cloudflare.json 核对你的 email, api_key, zone_id, domain 并重新启动")
		return
	}

	domainID, recordIP, err := cfc.GetParseRecord()
	if err != nil {
		return
	}
	recordType := ""
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
	} else {
		recordType = "A"
	}
	if domainID != cfc.DomainID {
		cfc.DomainID = domainID
		err = common.MarshalAndSave(cfc, ConfPath+"/cloudflare.json")
		if err != nil {
			return
		}
	}
	if recordIP == ipAddr {
		err = errors.New("Cloudflare 记录的 IP 和当前获取的 IP 一致")
		return
	}

	err = cfc.UpdateParseRecord(ipAddr, recordType)
	if err != nil {
		return
	}
	return
}

func (cfc CloudflareConf) GetParseRecord() (domainID, recordIP string, err error) {
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
	records, err := jsonObj.Get("result").Array()
	if len(records) == 0 {
		err = errors.New("解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"] == cfc.Domain {
			domainID = element["id"].(string)
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
		Proxied: false,
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
	return
}
