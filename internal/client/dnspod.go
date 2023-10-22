package client

import (
	"ddns-watchdog/internal/common"
	"errors"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"strings"
)

const DNSPodConfFileName = "dnspod.json"

type DNSPod struct {
	ID        string           `json:"id"`
	Token     string           `json:"token"`
	Domain    string           `json:"domain"`
	SubDomain common.Subdomain `json:"sub_domain"`
}

func (dpc *DNSPod) InitConf() (msg string, err error) {
	*dpc = DNSPod{
		ID:     "在 https://console.dnspod.cn/account/token/token 获取",
		Domain: "example.com",
		SubDomain: common.Subdomain{
			A:    "A记录子域名",
			AAAA: "AAAA记录子域名",
		},
	}
	dpc.Token = dpc.ID

	err = common.MarshalAndSave(dpc, ConfDirectoryName+"/"+DNSPodConfFileName)
	msg = "初始化 " + ConfDirectoryName + "/" + DNSPodConfFileName
	return
}

func (dpc *DNSPod) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+DNSPodConfFileName, &dpc)
	if err != nil {
		return
	}
	if dpc.ID == "" || dpc.Token == "" || dpc.Domain == "" || (dpc.SubDomain.A == "" && dpc.SubDomain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + DNSPodConfFileName + " 检查你的 id, token, domain, sub_domain 并重新启动")
	}
	return
}

func (dpc *DNSPod) Run(enabled common.Enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if enabled.IPv4 && dpc.SubDomain.A != "" {
		// 获取解析记录
		recordId, recordLineId, recordIP, err := dpc.getParseRecord(dpc.SubDomain.A, "A", msg)
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv4 {
			// 更新解析记录
			err = dpc.updateParseRecord(ipv4, recordId, recordLineId, "A", dpc.SubDomain.A)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "DNSPod: "+dpc.SubDomain.A+"."+dpc.Domain+" 已更新解析记录 "+ipv4)
			}
		}
	}
	if enabled.IPv6 && dpc.SubDomain.AAAA != "" {
		// 获取解析记录
		recordId, recordLineId, recordIP, err := dpc.getParseRecord(dpc.SubDomain.AAAA, "AAAA", msg)
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv6 {
			// 更新解析记录
			err = dpc.updateParseRecord(ipv6, recordId, recordLineId, "AAAA", dpc.SubDomain.AAAA)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "DNSPod: "+dpc.SubDomain.AAAA+"."+dpc.Domain+" 已更新解析记录 "+ipv6)
			}
		}
	}
	return
}

func checkRespondStatus(jsonObj *simplejson.Json) (err error) {
	statusCode := jsonObj.Get("status").Get("code").MustString()
	if statusCode != "1" {
		err = errors.New("DNSPod: " + statusCode + ": " + jsonObj.Get("status").Get("message").MustString())
		return
	}
	return
}

func (dpc *DNSPod) getParseRecord(subDomain, recordType string) (recordId, recordLineId, recordIP string, msg []string, err error) {
	postContent := dpc.publicRequestInit()
	postContent = postContent + "&" + dpc.recordRequestInit(subDomain)
	recvJson, err := postman("https://dnsapi.cn/Record.List", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = checkRespondStatus(jsonObj)
	if err != nil {
		return
	}
	records, err := jsonObj.Get("records").Array()
	if len(records) == 0 {
		msg = append(msg, "DNSPod: " + subDomain + "." + dpc.Domain + " 解析记录不存在")
	}
	for _, value := range records {
		element := value.(map[string]any)
		if element["name"].(string) == subDomain && element["type"].(string) == recordType {
			recordId = element["id"].(string)
			recordIP = element["value"].(string)
			recordLineId = element["line_id"].(string)
			break
		}
	}
	if recordId == "" || recordIP == "" || recordLineId == "" {
		msg = append(msg, "DNSPod: " + subDomain + "." + dpc.Domain + " 的 " + recordType + " 解析记录不存在")
	}
	return
}

func (dpc *DNSPod) updateParseRecord(ipAddr, recordId, recordLineId, recordType, subDomain string) (err error) {
	postContent := dpc.publicRequestInit()
	postContent = postContent + "&" + dpc.recordModifyRequestInit(ipAddr, recordId, recordLineId, recordType, subDomain)
	recvJson, err := postman("https://dnsapi.cn/Record.Modify", postContent)
	if err != nil {
		return
	}

	jsonObj, err := simplejson.NewJson(recvJson)
	if err != nil {
		return
	}
	err = checkRespondStatus(jsonObj)
	if err != nil {
		return
	}
	return
}

func (dpc *DNSPod) publicRequestInit() (pp string) {
	pp = "login_token=" + dpc.ID + "," + dpc.Token +
		"&format=" + "json" +
		"&lang=" + "cn" +
		"&error_on_empty=" + "no"
	return
}

func (dpc *DNSPod) recordRequestInit(subDomain string) (rr string) {
	rr = "domain=" + dpc.Domain +
		"&sub_domain=" + subDomain
	return
}

func (dpc *DNSPod) recordModifyRequestInit(ipAddr, recordId, recordLineId, recordType, subDomain string) (rm string) {
	rm = "domain=" + dpc.Domain +
		"&record_id=" + recordId +
		"&sub_domain=" + subDomain +
		"&record_type=" + recordType +
		"&record_line_id=" + recordLineId +
		"&value=" + ipAddr
	return
}

func postman(url, src string) (dst []byte, err error) {
	httpClient := getGeneralHttpClient()
	req, err := http.NewRequest("POST", url, strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		t := Body.Close()
		if t != nil {
			err = t
		}
	}(req.Body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ProjName+"/"+common.LocalVersion+" ()")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		t := Body.Close()
		if t != nil {
			err = t
		}
	}(resp.Body)
	dst, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}
