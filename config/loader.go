package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfig(filePath string) *Config {
	fileByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("load config file fail")
		log.Fatal(err)
	}

	c := Config{}

	err = json.Unmarshal(fileByte, &c)
	if err != nil {
		log.Println("fail to parse config file")
		log.Fatal(err)
	}

	return &c
}
