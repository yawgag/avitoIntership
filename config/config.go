package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	ServerAddress string `json:"ServerAddress"`
	DbURL         string `json:"DbURL"`
	SecretWord    string `json:"SecretWord"`
}

func LoadConfig() (*Config, error) {
	configFile, err := os.Open("../config/config.json")

	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	byteValue, _ := io.ReadAll(configFile)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
