package client

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yzy613/ddns-watchdog/common"
	"log"
)

func (adc *aliDNSConf) InitConf() (msg string, err error) {
	*adc = aliDNSConf{}
	adc.AccessKeyId = "在 https://ram.console.aliyun.com/users 获取"
	adc.AccessKeySecret = adc.AccessKeyId
	adc.Domain = "example.com"
	adc.SubDomain.A = "ipv4"
	adc.SubDomain.AAAA = "ipv6"
	err = common.MarshalAndSave(adc, ConfPath+AliDNSConfFileName)
	msg = "初始化 " + ConfPath + AliDNSConfFileName
	return
}

func (adc *aliDNSConf) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfPath+AliDNSConfFileName, &adc)
	if err != nil {
		return
	}
	if adc.AccessKeyId == "" || adc.AccessKeySecret == "" || adc.Domain == "" || (adc.SubDomain.A == "" && adc.SubDomain.AAAA == "") {
		log.Println("请打开配置文件 " + ConfPath + AliDNSConfFileName + " 检查你的 accesskey_id, accesskey_secret, domain, sub_domain 并重新启动")
	}
	return
}

func (adc aliDNSConf) Run(enabled enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if enabled.IPv4 {
		// 获取解析记录
		recordIP, err := adc.GetParseRecord(adc.SubDomain.A, "A")
		if err != nil {
			errs = append(errs, err)
		} else {
			if recordIP != ipv4 {
				// 更新解析记录
				err = adc.UpdateParseRecord(ipv4, "A", adc.SubDomain.A)
				if err != nil {
					errs = append(errs, err)
				} else {
					msg = append(msg, "AliDNS: "+adc.SubDomain.A+"."+adc.Domain+" 已更新解析记录 "+ipv4)
				}
			}
		}
	}
	if enabled.IPv6 {
		// 获取解析记录
		recordIP, err := adc.GetParseRecord(adc.SubDomain.AAAA, "AAAA")
		if err != nil {
			errs = append(errs, err)
		} else {
			if recordIP != ipv6 {
				// 更新解析记录
				err = adc.UpdateParseRecord(ipv6, "AAAA", adc.SubDomain.AAAA)
				if err != nil {
					errs = append(errs, err)
				} else {
					msg = append(msg, "AliDNS: "+adc.SubDomain.AAAA+"."+adc.Domain+" 已更新解析记录 "+ipv6)
				}
			}
		}
	}
	return
}

func (adc *aliDNSConf) GetParseRecord(subDomain, recordType string) (recordIP string, err error) {
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

func (adc aliDNSConf) UpdateParseRecord(ipAddr, recordType, subDomain string) (err error) {
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
