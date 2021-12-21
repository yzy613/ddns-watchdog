package client

import (
	"ddns-watchdog/internal/common"
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (cfc *cloudflareConf) InitConf() (msg string, err error) {
	*cfc = cloudflareConf{}
	cfc.APIToken = "在 https://dash.cloudflare.com/profile/api-tokens 获取"
	cfc.ZoneID = "在你域名页面的右下角有个区域 ID"
	cfc.Domain.A = "A记录子域名.example.com"
	cfc.Domain.AAAA = "AAAA记录子域名.example.com"
	err = common.MarshalAndSave(cfc, ConfPath+CloudflareConfFileName)
	msg = "初始化 " + ConfPath + CloudflareConfFileName
	return
}

func (cfc *cloudflareConf) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfPath+CloudflareConfFileName, &cfc)
	if err != nil {
		return
	}
	if cfc.ZoneID == "" || cfc.APIToken == "" || (cfc.Domain.A == "" && cfc.Domain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfPath + CloudflareConfFileName + " 检查你的 zone_id, api_token, domain 并重新启动")
	}
	return
}

func (cfc cloudflareConf) Run(enabled enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if enabled.IPv4 && cfc.Domain.A != "" {
		// 获取解析记录
		recordIP, err := cfc.GetParseRecord(cfc.Domain.A, "A")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv4 {
			// 更新解析记录
			err = cfc.UpdateParseRecord(ipv4, "A", cfc.Domain.A)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "Cloudflare: "+cfc.Domain.A+" 已更新解析记录 "+ipv4)
			}
		}
	}
	if enabled.IPv6 && cfc.Domain.AAAA != "" {
		// 获取解析记录
		recordIP, err := cfc.GetParseRecord(cfc.Domain.AAAA, "AAAA")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv6 {
			// 更新解析记录
			err = cfc.UpdateParseRecord(ipv6, "AAAA", cfc.Domain.AAAA)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "Cloudflare: "+cfc.Domain.AAAA+" 已更新解析记录 "+ipv6)
			}
		}
	}
	return
}

func (cfc *cloudflareConf) GetParseRecord(domain, recordType string) (recordIP string, err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records?name=" + domain
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+cfc.APIToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(res.Body)
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	if err2 := jsonObj.Get("error").MustString(); err2 != "" {
		err = errors.New("Cloudflare: " + err2)
		return
	}
	if !jsonObj.Get("success").MustBool() {
		err = errors.New("Cloudflare: 身份认证似乎有问题")
		return
	}
	records, err := jsonObj.Get("result").Array()
	if len(records) == 0 {
		err = errors.New("Cloudflare: " + domain + " 解析记录不存在")
		return
	}
	for _, value := range records {
		element := value.(map[string]interface{})
		if element["name"].(string) == domain {
			cfc.DomainID = element["id"].(string)
			recordIP = element["content"].(string)
			if element["type"].(string) == recordType {
				break
			}
		}
	}
	return
}

func (cfc cloudflareConf) UpdateParseRecord(ipAddr, recordType, domain string) (err error) {
	httpClient := &http.Client{}
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records/" + cfc.DomainID
	reqData := cloudflareUpdateRequest{
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
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(req.Body)
	req.Header.Set("Authorization", "Bearer "+cfc.APIToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(res.Body)
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
