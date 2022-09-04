package client

import (
	"ddns-watchdog/internal/common"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const AliDNSConfFileName = "alidns.json"

type AliDNS struct {
	AccessKeyId     string           `json:"accesskey_id"`
	AccessKeySecret string           `json:"accesskey_secret"`
	Domain          string           `json:"domain"`
	SubDomain       common.Subdomain `json:"sub_domain"`
	RecordId        string           `json:"-"`
}

func (ad *AliDNS) InitConf() (msg string, err error) {
	*ad = AliDNS{}
	ad.AccessKeyId = "在 https://ram.console.aliyun.com/users 获取"
	ad.AccessKeySecret = ad.AccessKeyId
	ad.Domain = "example.com"
	ad.SubDomain.A = "A记录子域名"
	ad.SubDomain.AAAA = "AAAA记录子域名"
	err = common.MarshalAndSave(ad, ConfDirectoryName+"/"+AliDNSConfFileName)
	msg = "初始化 " + ConfDirectoryName + "/" + AliDNSConfFileName
	return
}

func (ad *AliDNS) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+AliDNSConfFileName, &ad)
	if err != nil {
		return
	}
	if ad.AccessKeyId == "" || ad.AccessKeySecret == "" || ad.Domain == "" || (ad.SubDomain.A == "" && ad.SubDomain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + AliDNSConfFileName + " 检查你的 accesskey_id, accesskey_secret, domain, sub_domain 并重新启动")
	}
	return
}

func (ad AliDNS) Run(enabled common.Enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if enabled.IPv4 && ad.SubDomain.A != "" {
		// 获取解析记录
		recordIP, err := ad.getParseRecord(ad.SubDomain.A, "A")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv4 {
			// 更新解析记录
			err = ad.updateParseRecord(ipv4, "A", ad.SubDomain.A)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "AliDNS: "+ad.SubDomain.A+"."+ad.Domain+" 已更新解析记录 "+ipv4)
			}
		}
	}
	if enabled.IPv6 && ad.SubDomain.AAAA != "" {
		// 获取解析记录
		recordIP, err := ad.getParseRecord(ad.SubDomain.AAAA, "AAAA")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv6 {
			// 更新解析记录
			err = ad.updateParseRecord(ipv6, "AAAA", ad.SubDomain.AAAA)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "AliDNS: "+ad.SubDomain.AAAA+"."+ad.Domain+" 已更新解析记录 "+ipv6)
			}
		}
	}
	return
}

func (ad *AliDNS) getParseRecord(subDomain, recordType string) (recordIP string, err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ad.AccessKeyId, ad.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = ad.Domain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return
	}

	for i := range response.DomainRecords.Record {
		if response.DomainRecords.Record[i].RR == subDomain {
			ad.RecordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			if response.DomainRecords.Record[i].Type == recordType {
				break
			}
		}
	}
	if ad.RecordId == "" || recordIP == "" {
		err = errors.New("AliDNS: " + subDomain + "." + ad.Domain + " 解析记录不存在")
	}
	return
}

func (ad AliDNS) updateParseRecord(ipAddr, recordType, subDomain string) (err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ad.AccessKeyId, ad.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = ad.RecordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ipAddr

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	return
}
