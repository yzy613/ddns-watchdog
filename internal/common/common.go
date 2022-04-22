package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	LocalVersion      = "1.4.5"
	DefaultAPIUrl     = "https://yzyweb.cn/ddns-watchdog"
	DefaultIPv6APIUrl = "https://yzyweb.cn/ddns-watchdog6"
	ProjectUrl        = "https://github.com/yzy613/ddns-watchdog"
)

type PublicInfo struct {
	IP      string `json:"ip"`
	Version string `json:"latest_version"`
}

func FormatDirectoryPath(srcPath string) (dstPath string) {
	if length := len(srcPath); srcPath[length-1:] == "/" {
		dstPath = srcPath[0 : length-1]
	} else {
		dstPath = srcPath
	}
	return
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func IsDirExistAndCreate(dirPath string) (err error) {
	_, err = os.Stat(dirPath)
	if err != nil || os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0750)
		if err != nil {
			return err
		}
	}
	return
}

// LoadAndUnmarshal dst 参数要加 & 才能修改原变量
func LoadAndUnmarshal(filePath string, dst any) (err error) {
	_, err = os.Stat(filePath)
	if err != nil {
		return
	}
	jsonContent, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonContent, &dst)
	if err != nil {
		return
	}
	return
}

func MarshalAndSave(content any, filePath string) (err error) {
	err = IsDirExistAndCreate(filepath.Dir(filePath))
	if err != nil {
		return
	}
	jsonContent, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return
	}
	err = os.WriteFile(filePath, jsonContent, 0600)
	if err != nil {
		return
	}
	return nil
}

func CompareVersionString(remoteVersion, localVersion string) bool {
	rv := strings.Split(remoteVersion, ".")
	lv := strings.Split(localVersion, ".")
	if len(rv) <= len(lv) {
		for key, value := range rv {
			switch {
			case value > lv[key]:
				return true
			case value < lv[key]:
				return false
			}
		}
	}
	return false
}

func DecodeIPv6(srcIP string) (dstIP string) {
	if strings.Contains(srcIP, "::") {
		splitArr := strings.Split(srcIP, "::")
		decode := ""
		switch {
		case srcIP == "::":
			dstIP = "0:0:0:0:0:0:0:0"
		case splitArr[0] == "" && splitArr[1] != "":
			for i := 0; i < 8-len(strings.Split(splitArr[1], ":")); i++ {
				decode = "0:" + decode
			}
			dstIP = decode + splitArr[1]
		case splitArr[0] != "" && splitArr[1] == "":
			for i := 0; i < 8-len(strings.Split(splitArr[0], ":")); i++ {
				decode = decode + ":0"
			}
			dstIP = splitArr[0] + decode
		default:
			for i := 0; i < 8-len(strings.Split(splitArr[0], ":"))-len(strings.Split(splitArr[1], ":")); i++ {
				decode = decode + ":0"
			}
			decode = decode + ":"
			dstIP = splitArr[0] + decode + splitArr[1]
		}
	} else {
		dstIP = srcIP
	}
	return
}

func VersionTips(LatestVersion string) {
	fmt.Println("当前版本 ", LocalVersion)
	fmt.Println("最新版本 ", LatestVersion)
	fmt.Println("项目地址 ", ProjectUrl)
	switch {
	case strings.Contains(LatestVersion, "N/A"):
		fmt.Println("\n" + LatestVersion + "\n需要手动检查更新，请前往 项目地址 查看")
	case CompareVersionString(LatestVersion, LocalVersion):
		fmt.Println("\n发现新版本，请前往 项目地址 下载")
	}
}
