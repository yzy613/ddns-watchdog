package client

import (
	"ddns-watchdog/internal/common"
	"errors"
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
			A:    "A记录子域名.example.com.",
			AAAA: "AAAA记录子域名.example.com.",
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

func (hc *HuaweiCloud) Run(enabled common.Enable, ipv4, ipv6 string) (msg []string, errs []error) {
	if hc.ZoneId == "" && (enabled.IPv4 || enabled.IPv6) {
		err := hc.getZoneId()
		if err != nil {
			errs = append(errs, err)
			return
		}
	}
	if enabled.IPv4 && hc.Domain.A != "" {
		recordSetId, recordIP, err := hc.getParseRecord(hc.Domain.A, "A")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv4 {
			err = hc.updateParseRecord(ipv4, recordSetId, "A", hc.Domain.A)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "HuaweiCloud: "+hc.Domain.A+" 已更新解析记录 "+ipv4)
			}
		}
	}
	if enabled.IPv6 && hc.Domain.AAAA != "" {
		recordSetId, recordIP, err := hc.getParseRecord(hc.Domain.AAAA, "AAAA")
		if err != nil {
			errs = append(errs, err)
		} else if recordIP != ipv6 {
			err = hc.updateParseRecord(ipv6, recordSetId, "AAAA", hc.Domain.AAAA)
			if err != nil {
				errs = append(errs, err)
			} else {
				msg = append(msg, "HuaweiCloud: "+hc.Domain.AAAA+" 已更新解析记录 "+ipv6)
			}
		}
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
