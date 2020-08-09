package client

import (
	"ddns/common"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"strings"
)

func Aliyun(ipAddr string) (err error) {
	ayc := AliyunConf{}
	// 获取配置
	err = common.LoadAndUnmarshal(ConfPath+"/aliyun.json", &ayc)
	if err != nil {
		return
	}

	if ayc.AccessKeyId == "" || ayc.AccessKeySecret == "" || ayc.Domain == "" || ayc.SubDomain == "" {
		err = errors.New("请打开配置文件 " + ConfPath + "/aliyun.json 核对你的 accesskey_id, accesskey_secret, domain, sub_domain 并重新启动")
		return
	}
	// 获取解析记录
	recordId, recordIP, err := ayc.GetParseRecord()
	if err != nil {
		return
	}
	ayc.RecordId = recordId
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

func (ayc AliyunConf) GetParseRecord() (recordId, recordIP string, err error) {
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
			recordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			break
		}
	}
	if recordId == "" || recordIP == "" {
		err = errors.New("阿里云: "+ayc.SubDomain + "." + ayc.Domain + " 解析记录不存在")
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
