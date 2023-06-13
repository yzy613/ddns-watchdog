package client

import (
	"ddns-watchdog/internal/common"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/bitly/go-simplejson"
)

const CloudflareConfFileName = "cloudflare.json"

type Cloudflare struct {
	ZoneID   string           `json:"zone_id"`
	APIToken string           `json:"api_token"`
	Domain   common.Subdomain `json:"domain"`
}

type cloudflareUpdateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

func (cfc *Cloudflare) InitConf() (msg string, err error) {
	*cfc = Cloudflare{
		ZoneID:   "在你域名页面的右下角有个区域 ID",
		APIToken: "在 https://dash.cloudflare.com/profile/api-tokens 获取",
		Domain: common.Subdomain{
			A:    "A记录子域名1.example.com,A记录子域名2.example.com,……",
			AAAA: "AAAA记录子域名1.example.com,AAAA记录子域名2.example.com,……",
		},
	}

	err = common.MarshalAndSave(cfc, ConfDirectoryName+"/"+CloudflareConfFileName)
	msg = "初始化 " + ConfDirectoryName + "/" + CloudflareConfFileName
	return
}

func (cfc *Cloudflare) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+CloudflareConfFileName, &cfc)
	if err != nil {
		return
	}
	if cfc.ZoneID == "" || cfc.APIToken == "" || (cfc.Domain.A == "" && cfc.Domain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + CloudflareConfFileName + " 检查你的 zone_id, api_token, domain 并重新启动")
	}
	return
}

func (cfc *Cloudflare) processRecords(records []string, recordType string, ip string, msgChan chan<- string, errChan chan<- error, wg *sync.WaitGroup) {
	for _, record := range records {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			// 获取解析记录
			domainId, recordIP, err := cfc.getParseRecord(item, recordType)
			if err != nil {
				errChan <- err
			} else if recordIP != ip {
				// 更新解析记录
				err = cfc.updateParseRecord(ip, domainId, recordType, item)
				if err != nil {
					errChan <- err
				} else {
					msgChan <- "Cloudflare: " + item + " 已更新解析记录 " + ip
				}
			}
		}(record)
	}
}

func (cfc *Cloudflare) Run(enabled common.Enable, ipv4, ipv6 string) (msgs []string, errs []error) {
	AArr := common.DomainStr2Arr(cfc.Domain.A)
	AAAAArr := common.DomainStr2Arr(cfc.Domain.AAAA)

	var wg sync.WaitGroup
	msgChan := make(chan string)
	errChan := make(chan error)

	if enabled.IPv4 {
		cfc.processRecords(AArr, "A", ipv4, msgChan, errChan, &wg)
	}

	if enabled.IPv6 {
		cfc.processRecords(AAAAArr, "AAAA", ipv6, msgChan, errChan, &wg)
	}

	go func() {
		wg.Wait()
		close(msgChan)
		close(errChan)
	}()

	for msg := range msgChan {
		msgs = append(msgs, msg)
	}

	for err := range errChan {
		errs = append(errs, err)
	}

	return
}

func (cfc *Cloudflare) getParseRecord(domain, recordType string) (domainId, recordIP string, err error) {
	httpClient := getGeneralHttpClient()
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records?name=" + domain
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+cfc.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		t := Body.Close()
		if t != nil {
			err = t
		}
	}(resp.Body)
	recvJson, err := io.ReadAll(resp.Body)
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
	for _, v := range records {
		element := v.(map[string]any)
		if element["name"].(string) == domain && element["type"].(string) == recordType {
			domainId = element["id"].(string)
			recordIP = element["content"].(string)
			break
		}
	}
	if domainId == "" || recordIP == "" {
		err = errors.New("Cloudflare: " + domain + " 的 " + recordType + " 解析记录不存在")
	}
	return
}

func (cfc *Cloudflare) updateParseRecord(ipAddr, domainId, recordType, domain string) (err error) {
	httpClient := getGeneralHttpClient()
	url := "https://api.cloudflare.com/client/v4/zones/" + cfc.ZoneID + "/dns_records/" + domainId
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
		t := Body.Close()
		if t != nil {
			err = t
		}
	}(req.Body)
	req.Header.Set("Authorization", "Bearer "+cfc.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		t := Body.Close()
		if t != nil {
			err = t
		}
	}(resp.Body)
	recvJson, err := io.ReadAll(resp.Body)
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
			element := value.(map[string]any)
			errCode := element["code"].(json.Number)
			errMsg := element["message"].(string)
			errorsMsg = errorsMsg + "Cloudflare: " + errCode.String() + ": " + errMsg + "\n"
		}
		err = errors.New(errorsMsg)
		return
	}
	return
}
