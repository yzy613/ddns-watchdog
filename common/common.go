package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	LocalVersion = "0.2.3"
	RootServer   = "https://yzyweb.cn/ddns"
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
			err = os.MkdirAll(dirPath, 0777)
			if err != nil {
				return
			}
		}
	}
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0744)
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
		n, err = dstFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return
}

func LoadAndUnmarshal(filePath string, dst interface{}) error {
	_, err := os.Stat(filePath)
	if err != nil || os.IsExist(err) {
		_, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
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
	err = ioutil.WriteFile(filePath, jsonContent, 0666)
	if err != nil {
		return
	}
	return nil
}

func Struct2Map(src interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	// 原始复制
	/*key := reflect.TypeOf(src)
	value := reflect.ValueOf(src)
	for i := 0; i < key.NumField(); i++ {
		if value.Field(i).Interface() == "" {
			continue
		}
		dst[key.Field(i).Name] = value.Field(i).Interface()
	}
	return dst*/

	// 以 json 格式复制
	tmpJson, getErr := json.Marshal(src)
	if getErr != nil {
		fmt.Println(getErr)
	}
	getErr = json.Unmarshal(tmpJson, &dst)
	if getErr != nil {
		fmt.Println(getErr)
	}
	return dst
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
		split := strings.Split(srcIP, "::")
		decode := ""
		switch {
		case srcIP == "::":
			dstIP = "0:0:0:0:0:0:0:0"
		case split[0] == "" && split[1] != "":
			for i := 0; i < 8-len(strings.Split(split[1], ":")); i++ {
				decode = "0:" + decode
			}
			dstIP = decode + split[1]
		case split[0] != "" && split[1] == "":
			for i := 0; i < 8-len(strings.Split(split[0], ":")); i++ {
				decode = decode + ":0"
			}
			dstIP = split[0] + decode
		default:
			for i := 0; i < 8-len(strings.Split(split[0], ":"))-len(strings.Split(split[1], ":")); i++ {
				decode = decode + ":0"
			}
			decode = decode + ":"
			dstIP = split[0] + decode + split[1]
		}
	} else {
		dstIP = srcIP
	}
	return
}
