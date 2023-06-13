package client

import (
	"ddns-watchdog/internal/common"
	"errors"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const AliDNSConfFileName = "alidns.json"

type AliDNS struct {
	AccessKeyId     string           `json:"access_key_id"`
	AccessKeySecret string           `json:"access_key_secret"`
	Domain          string           `json:"domain"`
	SubDomain       common.Subdomain `json:"sub_domain"`
}

func (ad *AliDNS) InitConf() (msg string, err error) {
	*ad = AliDNS{
		AccessKeyId: "在 https://ram.console.aliyun.com/users 获取",
		Domain:      "example.com",
		SubDomain: common.Subdomain{
			A:    "A记录子域名1,A记录子域名2,……",
			AAAA: "AAAA记录子域名1,AAAA记录子域名2,……",
		},
	}
	ad.AccessKeySecret = ad.AccessKeyId

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
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + AliDNSConfFileName + " 检查你的 access_key_id, access_key_secret, domain, sub_domain 并重新启动")
	}
	return
}

func (ad *AliDNS) processRecords(records []string, recordType string, ip string, msgChan chan<- string, errChan chan<- error, wg *sync.WaitGroup) {
	for _, record := range records {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			// 获取解析记录
			recordId, recordIP, err := ad.getParseRecord(item, recordType)
			if err != nil {
				errChan <- err
			} else if recordIP != ip {
				// 更新解析记录
				err = ad.updateParseRecord(ip, recordId, recordType, item)
				if err != nil {
					errChan <- err
				} else {
					msgChan <- "AliDNS: " + item + "." + ad.Domain + " 已更新解析记录 " + ip
				}
			}
		}(record)
	}
}

func (ad *AliDNS) Run(enabled common.Enable, ipv4, ipv6 string) (msgs []string, errs []error) {
	AArr := common.DomainStr2Arr(ad.SubDomain.A)
	AAAAArr := common.DomainStr2Arr(ad.SubDomain.AAAA)

	var wg sync.WaitGroup
	msgChan := make(chan string)
	errChan := make(chan error)

	if enabled.IPv4 {
		ad.processRecords(AArr, "A", ipv4, msgChan, errChan, &wg)
	}

	if enabled.IPv6 {
		ad.processRecords(AAAAArr, "AAAA", ipv6, msgChan, errChan, &wg)
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

func (ad *AliDNS) getParseRecord(subDomain, recordType string) (recordId, recordIP string, err error) {
	dnsClient, err := alidns.NewClientWithAccessKey("cn-hangzhou", ad.AccessKeyId, ad.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = ad.Domain

	response, err := dnsClient.DescribeDomainRecords(request)
	if err != nil {
		return
	}

	for i := range response.DomainRecords.Record {
		if response.DomainRecords.Record[i].RR == subDomain &&
			response.DomainRecords.Record[i].Type == recordType {
			recordId = response.DomainRecords.Record[i].RecordId
			recordIP = response.DomainRecords.Record[i].Value
			break
		}
	}
	if recordId == "" || recordIP == "" {
		err = errors.New("AliDNS: " + subDomain + "." + ad.Domain + " 的 " + recordType + " 解析记录不存在")
	}
	return
}

func (ad *AliDNS) updateParseRecord(ipAddr, recordId, recordType, subDomain string) (err error) {
	dnsClient, err := alidns.NewClientWithAccessKey("cn-hangzhou", ad.AccessKeyId, ad.AccessKeySecret)
	if err != nil {
		return
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = recordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ipAddr

	_, err = dnsClient.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	return
}
