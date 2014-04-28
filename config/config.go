package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ApiKey   string `json:"api_key"`
	ApiUser  string `json:"api_user"`
	ApiURL   string `json:"api_url"`
	Category string `json:"category"`
}

func Parse(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var readcfg Config
	err = json.NewDecoder(f).Decode(&readcfg)
	if err != nil {
		return nil, err
	}
	return &readcfg, nil
}
