package client

import (
	"ddns-watchdog/internal/common"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const AliDNSConfFileName = "alidns.json"

type aliDNSConf struct {
	AccessKeyId     string    `json:"accesskey_id"`
	AccessKeySecret string    `json:"accesskey_secret"`
	Domain          string    `json:"domain"`
	SubDomain       subdomain `json:"sub_domain"`
	RecordId        string    `json:"-"`
}

func (adc *aliDNSConf) InitConf() (msg string, err error) {
	*adc = aliDNSConf{}
	adc.AccessKeyId = "在 https://ram.console.aliyun.com/users 获取"
	adc.AccessKeySecret = adc.AccessKeyId
	adc.Domain = "example.com"
	adc.SubDomain.A = "A记录子域名"
	adc.SubDomain.AAAA = "AAAA记录子域名"
	err = common.MarshalAndSave(adc, ConfDirectoryName+"/"+AliDNSConfFileName)
	msg = "初始化 " + ConfDirectoryName + "/" + AliDNSConfFileName
	return
}

func (adc *aliDNSConf) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+AliDNSConfFileName, &adc)
	if err != nil {
		return
	}
	if adc.AccessKeyId == "" || adc.AccessKeySecret == "" || adc.Domain == "" || (adc.SubDomain.A == "" && adc.SubDomain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + AliDNSConfFileName + " 检查你的 accesskey_id, accesskey_secret, domain, sub_domain 并重新启动")
	}
	return
}

func (adc aliDNSConf) Run(enabled enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if enabled.IPv4 && adc.SubDomain.A != "" {
		// 获取解析记录
		recordIP, err := adc.getParseRecord(adc.SubDomain.A, "A")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv4 {
			// 更新解析记录
			err = adc.updateParseRecord(ipv4, "A", adc.SubDomain.A)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "AliDNS: "+adc.SubDomain.A+"."+adc.Domain+" 已更新解析记录 "+ipv4)
			}
		}
	}
	if enabled.IPv6 && adc.SubDomain.AAAA != "" {
		// 获取解析记录
		recordIP, err := adc.getParseRecord(adc.SubDomain.AAAA, "AAAA")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv6 {
			// 更新解析记录
			err = adc.updateParseRecord(ipv6, "AAAA", adc.SubDomain.AAAA)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "AliDNS: "+adc.SubDomain.AAAA+"."+adc.Domain+" 已更新解析记录 "+ipv6)
			}
		}
	}
	return
}

func (adc *aliDNSConf) getParseRecord(subDomain, recordType string) (recordIP string, err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", adc.AccessKeyId, adc.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = adc.Domain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return
	}

	for i := range response.DomainRecords.Record {
		if response.DomainRecords.Record[i].RR == subDomain {
			adc.RecordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			if response.DomainRecords.Record[i].Type == recordType {
				break
			}
		}
	}
	if adc.RecordId == "" || recordIP == "" {
		err = errors.New("AliDNS: " + subDomain + "." + adc.Domain + " 解析记录不存在")
	}
	return
}

func (adc aliDNSConf) updateParseRecord(ipAddr, recordType, subDomain string) (err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", adc.AccessKeyId, adc.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = adc.RecordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ipAddr

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	return
}
