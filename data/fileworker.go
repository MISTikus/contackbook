package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileWorker struct {
	FileName string
}

func (file *FileWorker) Read(result interface{}) error {
	if !fileExists(file.FileName) {
		return nil;
	}

	data, err := ioutil.ReadFile(file.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func (file *FileWorker) Write(data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file.FileName, json, 0644)
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
