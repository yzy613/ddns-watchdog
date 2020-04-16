package client

import (
	"bytes"
	"ddns/common"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PublicParameterInit(dps DNSPodConf) (pp PublicParameter) {
	pp = PublicParameter{
		LoginToken:   dps.Id + "," + dps.Token,
		Format:       "json",
		Lang:         "cn",
		ErrorOnEmpty: "no",
	}
	return
}

func RecordListInit(dpc DNSPodConf) (ri RecordList) {
	ri = RecordList{
		Domain:    dpc.Domain,
		SubDomain: dpc.SubDomain,
	}
	return
}

func RecordModifyInit(dpc DNSPodConf, ipAddr string) (rm RecordModify) {
	rm = RecordModify{
		Domain:     dpc.Domain,
		RecordId:   dpc.RecordId,
		RecordType: "A",
		RecordLine: "默认",
		Value:      ipAddr,
	}
	return
}

func DNSPod(dpc DNSPodConf, ipAddr string) ([]byte, error) {
	postContent := common.Struct2Map(PublicParameterInit(dpc))
	tmpMap := common.Struct2Map(RecordListInit(dpc))
	for key, value := range tmpMap {
		postContent[key] = value
	}
/*
	stringContent := ""
	for key, value := range postContent {
		stringContent += fmt.Sprint()
	}
*/
	// 查询解析记录列表
	postJson, err := json.Marshal(postContent)
	if err != nil {
		return nil, err
	}
	res, err := http.Post("https://dnsapi.cn/Record.List", "application/json;charset=utf-8", bytes.NewReader(postJson))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 发送新的解析
	return recvJson, nil
}

func GetOwnIP(webAddr string) (ipAddr string, err error) {
	if webAddr == "" {
		webAddr = "https://yzyweb.cn/ddns"
	}
	res, err := http.Get(webAddr)
	defer res.Body.Close()
	if err != nil {
		return
	}
	recvJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var ipInfo common.IpInfoFormat
	err = json.Unmarshal(recvJson, &ipInfo)
	if err != nil {
		return
	}
	ipAddr = ipInfo.Ip
	return
}
