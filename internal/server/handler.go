package server

import (
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
	if v, ok := whitelist[body.Token]; ok {
		if !v.Enable {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 模拟客户端
	httpStatus, respBody, err := doVirtualClient(body, whitelist[body.Token])
	if httpStatus != http.StatusOK {
		if httpStatus == http.StatusInternalServerError {
			log.Println(err)
		}
		w.WriteHeader(httpStatus)
		return
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
