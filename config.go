package main

import (
	"encoding/json"
	"io/ioutil"
)

var (
	config map[string]Config
)

type Config struct {
	Logfile  string
	Secret   string
	Commands []string
}

// LoadConfig load the config
func LoadConfig() error {
	result, err := ioutil.ReadFile("data/config.json")
	if err != nil {
		return err
	}
	json.Unmarshal(result, &config)
	return nil
}

// 获取配置中的log文件名称
func GetLogName(id string) string {
	logName := config[id].Logfile
	if logName == "" || len(logName) <= 0 {
		return id
	}
	return logName
}
