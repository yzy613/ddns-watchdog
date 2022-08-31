package server

import (
	"ddns-watchdog/internal/client"
	"ddns-watchdog/internal/common"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RespGetIPReq(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")

	// 判断请求方法
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	info := common.GetIPResp{
		IP:      GetClientIP(req),
		Version: Srv.GetLatestVersion(),
	}
	sendJson, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(sendJson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RespCenterReq(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	
	// 判断请求方法
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 读取并解码 POST 请求
	bodyJson, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var body common.CenterReq
	err = json.Unmarshal(bodyJson, &body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 鉴权
	if len(body.Token) > 127 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if whitelist[body.Token] == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 模拟客户端
	var msg []string
	var errs []error
	switch body.Service {
	case common.DNSPod:
		if !Service.DNSPod.Enable {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		dp := client.DNSPod{}
		err = json.Unmarshal(body.Data, &dp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 补全
		dp.ID = Service.DNSPod.ID
		dp.Token = Service.DNSPod.Token
		msg, errs = dp.Run(body.Enable, body.IP.IPv4, body.IP.IPv6)
	case common.AliDNS:
		if !Service.AliDNS.Enable {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ad := client.AliDNS{}
		err = json.Unmarshal(body.Data, &ad)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 补全
		ad.AccessKeyId = Service.AliDNS.AccessKeyId
		ad.AccessKeySecret = Service.AliDNS.AccessKeySecret
		msg, errs = ad.Run(body.Enable, body.IP.IPv4, body.IP.IPv6)
	case common.Cloudflare:
		if !Service.Cloudflare.Enable {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		cf := client.Cloudflare{}
		err = json.Unmarshal(body.Data, &cf)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 补全
		cf.ZoneID = Service.Cloudflare.ZoneID
		cf.APIToken = Service.Cloudflare.APIToken
		msg, errs = cf.Run(body.Enable, body.IP.IPv4, body.IP.IPv6)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody := common.GeneralResp{}
	for _, v := range msg {
		respBody.Message += v + "\n"
	}
	for _, v := range errs {
		respBody.Message += v.Error() + "\n"
	}
	respJson, err := json.Marshal(respBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respJson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
