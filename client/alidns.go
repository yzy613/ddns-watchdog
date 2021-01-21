package client

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yzy613/ddns-watchdog/common"
	"strings"
)

func AliDNS(adc AliDNSConf, ipAddr string) (msg []string, err []error) {
	recordType := ""
	if strings.Contains(ipAddr, ":") {
		recordType = "AAAA"
		ipAddr = common.DecodeIPv6(ipAddr)
	} else {
		recordType = "A"
	}

	for _, subDomain := range adc.SubDomain {
		// 获取解析记录
		recordIP, currentErr := adc.GetParseRecord(subDomain, recordType)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		if recordIP == ipAddr {
			continue
		}
		// 更新解析记录
		currentErr = adc.UpdateParseRecord(ipAddr, recordType, subDomain)
		if currentErr != nil {
			err = append(err, currentErr)
			continue
		}
		msg = append(msg, "AliDNS: "+subDomain+"."+adc.Domain+" 已更新解析记录 "+ipAddr)
	}
	return
}

func (adc *AliDNSConf) GetParseRecord(subDomain, recordType string) (recordIP string, err error) {
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

func (adc AliDNSConf) UpdateParseRecord(ipAddr, recordType, subDomain string) (err error) {
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
