package client

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"log"
	"strings"
)

func AliDNS(ayc AliDNSConf, ipAddr string) (err error) {
	for _, subDomain := range ayc.SubDomain {
		// 获取解析记录
		recordIP, err := ayc.GetParseRecord(subDomain)
		if err != nil {
			log.Println(err)
			continue
		}
		recordType := ""
		if strings.Contains(ipAddr, ":") {
			recordType = "AAAA"
		} else {
			recordType = "A"
		}
		if recordIP == ipAddr {
			continue
		}
		// 更新解析记录
		err = ayc.UpdateParseRecord(ipAddr, recordType, subDomain)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("AliDNS: " + subDomain + "." + ayc.Domain + " 已更新解析记录 " + ipAddr)
	}
	return
}

func (ayc *AliDNSConf) GetParseRecord(subDomain string) (recordIP string, err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ayc.AccessKeyId, ayc.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = ayc.Domain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return
	}

	for i := range response.DomainRecords.Record {
		if response.DomainRecords.Record[i].RR == subDomain {
			ayc.RecordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			break
		}
	}
	if ayc.RecordId == "" || recordIP == "" {
		err = errors.New("AliDNS: " + subDomain + "." + ayc.Domain + " 解析记录不存在")
	}
	return
}

func (ayc AliDNSConf) UpdateParseRecord(ipAddr, recordType, subDomain string) (err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ayc.AccessKeyId, ayc.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = ayc.RecordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ipAddr

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	return
}
