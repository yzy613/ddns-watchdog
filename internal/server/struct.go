package server

import (
	"ddns-watchdog/internal/common"
)

const (
	RunningName  = "ddns-watchdog-server"
	ConfFileName = "server.json"
)

var (
	RunningPath = common.GetRunningPath()
	InstallPath = "/etc/systemd/system/" + RunningName + ".service"
	ConfPath    = RunningPath + "conf/"
)

type ServerConf struct {
	Port           string `json:"port"`
	IsRoot         bool   `json:"is_root"`
	RootServerAddr string `json:"root_server_addr"`
}
