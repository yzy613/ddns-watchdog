package client

import (
	"ddns-watchdog/internal/common"
	"errors"
	"sync"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
)

const HuaweiCloudConfFileName = "huaweicloud.json"

type HuaweiCloud struct {
	AccessKeyId     string           `json:"access_key_id"`
	SecretAccessKey string           `json:"secret_access_key"`
	ZoneName        string           `json:"zone_name"`
	Domain          common.Subdomain `json:"domain"`
	ZoneId          string           `json:"-"`
}

func (hc *HuaweiCloud) InitConf() (msg string, err error) {
	*hc = HuaweiCloud{
		AccessKeyId: "在 https://console.huaweicloud.com/iam/ 获取",
		ZoneName:    "example.com.",
		Domain: common.Subdomain{
			A:    "A记录子域名1.example.com,A记录子域名2.example.com,……",
			AAAA: "AAAA记录子域名1.example.com,AAAA记录子域名2.example.com,……",
		},
	}
	hc.SecretAccessKey = hc.AccessKeyId

	err = common.MarshalAndSave(hc, ConfDirectoryName+"/"+HuaweiCloudConfFileName)
	msg = "初始化 " + ConfDirectoryName + "/" + HuaweiCloudConfFileName
	return
}

func (hc *HuaweiCloud) LoadConf() (err error) {
	err = common.LoadAndUnmarshal(ConfDirectoryName+"/"+HuaweiCloudConfFileName, &hc)
	if err != nil {
		return
	}
	if hc.AccessKeyId == "" || hc.SecretAccessKey == "" || hc.ZoneName == "" || (hc.Domain.A == "" && hc.Domain.AAAA == "") {
		err = errors.New("请打开配置文件 " + ConfDirectoryName + "/" + HuaweiCloudConfFileName + " 检查你的 access_key_id, secret_access_key, domain 并重新启动")
	}
	return
}

func (hc *HuaweiCloud) processRecords(records []string, recordType string, ip string, msgChan chan<- string, errChan chan<- error, wg *sync.WaitGroup) {
	for _, record := range records {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			// 获取解析记录
			recordSetId, recordIP, err := hc.getParseRecord(item, recordType)
			if err != nil {
				errChan <- err
			} else if recordIP != ip {
				// 更新解析记录
				err = hc.updateParseRecord(ip, recordSetId, recordType, item)
				if err != nil {
					errChan <- err
				} else {
					msgChan <- "HuaweiCloud: " + item + " 已更新解析记录 " + ip
				}
			}
		}(record)
	}
}

func (hc *HuaweiCloud) Run(enabled common.Enable, ipv4, ipv6 string) (msgs []string, errs []error) {
	AArr := common.DomainStr2Arr(hc.Domain.A)
	AAAAArr := common.DomainStr2Arr(hc.Domain.AAAA)

	if hc.ZoneId == "" && (enabled.IPv4 || enabled.IPv6) {
		err := hc.getZoneId()
		if err != nil {
			errs = append(errs, err)
			return
		}
	}

	var wg sync.WaitGroup
	msgChan := make(chan string)
	errChan := make(chan error)

	if enabled.IPv4 {
		hc.processRecords(AArr, "A", ipv4, msgChan, errChan, &wg)
	}

	if enabled.IPv6 {
		hc.processRecords(AAAAArr, "AAAA", ipv6, msgChan, errChan, &wg)
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

func (hc *HuaweiCloud) getZoneId() (err error) {
	auth := basic.NewCredentialsBuilder().
		WithAk(hc.AccessKeyId).
		WithSk(hc.SecretAccessKey).
		Build()

	dnsClient := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())

	request := &model.ListPublicZonesRequest{}
	response, err := dnsClient.ListPublicZones(request)
	if err != nil {
		return
	}
	for _, v := range *response.Zones {
		if *v.Name == hc.ZoneName {
			hc.ZoneId = *v.Id
		}
	}
	if hc.ZoneId == "" {
		err = errors.New("HuaweiCloud: " + hc.ZoneName + " Zone 不存在")
	}
	return
}

func (hc *HuaweiCloud) getParseRecord(domain, recordType string) (recordSetId, recordIP string, err error) {
	auth := basic.NewCredentialsBuilder().
		WithAk(hc.AccessKeyId).
		WithSk(hc.SecretAccessKey).
		Build()

	dnsClient := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())

	request := &model.ListRecordSetsByZoneRequest{}
	request.ZoneId = hc.ZoneId
	response, err := dnsClient.ListRecordSetsByZone(request)
	if err != nil {
		return
	}
	for _, v := range *response.Recordsets {
		if *v.Name == domain && *v.Type == recordType {
			recordSetId = *v.Id
			for _, vv := range *v.Records {
				recordIP = vv
			}
			break
		}
	}
	if recordSetId == "" || recordIP == "" {
		err = errors.New("HuaweiCloud: " + domain + " 的 " + recordType + " 解析记录不存在")
	}
	return
}

func (hc *HuaweiCloud) updateParseRecord(ipAddr, recordSetId, recordType, domain string) (err error) {
	auth := basic.NewCredentialsBuilder().
		WithAk(hc.AccessKeyId).
		WithSk(hc.SecretAccessKey).
		Build()

	dnsClient := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())

	request := &model.UpdateRecordSetRequest{}
	request.ZoneId = hc.ZoneId
	request.RecordsetId = recordSetId
	var listRecordsbody = []string{
		ipAddr,
	}
	request.Body = &model.UpdateRecordSetReq{
		Records: &listRecordsbody,
		Type:    recordType,
		Name:    domain,
	}
	_, err = dnsClient.UpdateRecordSet(request)
	return
}
