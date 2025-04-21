package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DbURL         string
	SecretWord    string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DbURL:         os.Getenv("DB_URL"),
		SecretWord:    os.Getenv("SECRET_WORD"),
	}

	return config, nil
}
