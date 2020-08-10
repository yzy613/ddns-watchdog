package client

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"strings"
)

func Aliyun(ayc AliyunConf, ipAddr string) (err error) {
	// 获取解析记录
	recordIP, err := ayc.GetParseRecord()
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
		err = errors.New("阿里云记录的 IP 和当前获取的 IP 一致")
		return
	}
	// 更新解析记录
	err = ayc.UpdateParseRecord(ipAddr, recordType)
	if err != nil {
		return
	}
	return
}

func (ayc *AliyunConf) GetParseRecord() (recordIP string, err error) {
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
		if response.DomainRecords.Record[i].RR == ayc.SubDomain {
			ayc.RecordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			break
		}
	}
	if ayc.RecordId == "" || recordIP == "" {
		err = errors.New("阿里云: " + ayc.SubDomain + "." + ayc.Domain + " 解析记录不存在")
	}
	return
}

func (ayc AliyunConf) UpdateParseRecord(ipAddr string, recordType string) (err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", ayc.AccessKeyId, ayc.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = ayc.RecordId
	request.RR = ayc.SubDomain
	request.Type = recordType
	request.Value = ipAddr

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	return
}
