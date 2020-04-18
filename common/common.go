package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const LocalVersion = "0.1.1"

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

func LoadAndUnmarshal(filePath string, dst interface{}) error {
	_, getErr := os.Stat(filePath)
	if getErr != nil || os.IsExist(getErr) {
		_, getErr = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
		if getErr != nil {
			return getErr
		}
	}
	jsonContent, getErr := ioutil.ReadFile(filePath)
	if getErr != nil {
		return getErr
	}
	getErr = json.Unmarshal(jsonContent, &dst)
	if getErr != nil {
		return getErr
	}
	return nil
}

func MarshalAndSave(content interface{}, filePath string) error {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, jsonContent, 0666)
	if err != nil {
		return err
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

func CompareVersionString(remoteVersion string, localVersion string) bool {
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