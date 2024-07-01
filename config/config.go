package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvSource struct {
	Prefix string
}

type Config struct {
	APIKey      string
	ProxyPort   int
	RedirectURL string
}

func (es EnvSource) Load() (Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config

	err := envconfig.Process(es.Prefix, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
