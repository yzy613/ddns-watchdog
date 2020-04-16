package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadAndUnmarshal(filePath string, dst interface{}) error {
	_, getErr := os.Stat(filePath)
	if getErr != nil {
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
	err = ioutil.WriteFile(filePath, jsonContent, 0664)
	if err != nil {
		return err
	}
	return nil
}
