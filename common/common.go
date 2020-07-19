package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	LocalVersion = "0.2.3"
	RootServer   = "https://yzyweb.cn/ddns"
	ProjectAddr  = "https://github.com/yzy613/ddns/releases"
)

func IsDirExistAndCreate(dirPath string) (err error) {
	_, err = os.Stat(dirPath)
	if err != nil || !os.IsExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return
}

func CopyFile(srcPath, dstPath string) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dirSplit := strings.Split(dstPath, "/")
	dirPath := ""
	if dirPathLen := len(dirSplit); dirPathLen > 1 {
		switch dirSplit[0] {
		case ".":
			dirSplit = dirSplit[1:]
			dirPath = "./"
		case "":
			dirSplit = dirSplit[1:]
			dirPath = "/"
		}
		if dirPathLen := len(dirSplit); dirPathLen > 1 {
			for i := 0; i < dirPathLen-1; i++ {
				dirPath = dirPath + dirSplit[i] + "/"
			}
			err = os.MkdirAll(dirPath, 0750)
			if err != nil {
				return
			}
		}
	}
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		_, err = dstFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return
}

func LoadAndUnmarshal(filePath string, dst interface{}) error {
	_, err := os.Stat(filePath)
	if err != nil || os.IsExist(err) {
		_, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
	}
	jsonContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonContent, &dst)
	if err != nil {
		return err
	}
	return nil
}

func MarshalAndSave(content interface{}, filePath string) (err error) {
	jsonContent, err := json.Marshal(content)
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
			if value > lv[key] {
				return true
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
