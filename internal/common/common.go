package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetRunningPath() (path string) {
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	tempStr := strings.ReplaceAll(path, "\\", "/")
	if tempStr[len(tempStr)-1:] != "/" {
		tempStr = tempStr + "/"
	}
	return tempStr
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
func LoadAndUnmarshal(filePath string, dst interface{}) (err error) {
	_, err = os.Stat(filePath)
	if err != nil {
		return
	}
	jsonContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonContent, &dst)
	if err != nil {
		return
	}
	return
}

func MarshalAndSave(content interface{}, filePath string) (err error) {
	err = IsDirExistAndCreate(filepath.Dir(filePath))
	if err != nil {
		return
	}
	jsonContent, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filePath, jsonContent, 0600)
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
