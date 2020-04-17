package server

type ServerConf struct {
	Port           string `json:"port"`
	IsRoot         bool   `json:"is_root"`
	RootServerAddr string `json:"root_server_addr"`
}
