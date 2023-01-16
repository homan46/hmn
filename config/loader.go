package config

import (
	"encoding/json"
	"io/ioutil"

	"codeberg.org/rchan/hmn/log"

	"go.uber.org/zap"
)

func LoadConfig(filePath string) *Config {
	fileByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.ZLog.Panic("fail to load config file", zap.Error(err))
	}

	c := Config{}

	err = json.Unmarshal(fileByte, &c)
	if err != nil {
		log.ZLog.Panic("fail to prase config file", zap.Error(err))
	}

	return &c
}
