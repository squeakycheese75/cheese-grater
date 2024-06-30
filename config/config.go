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
	ProxyPort   int    `envconfig:"PROXY_PORT"`
	RedirectURL string `envconfig:"REDIRECT_URL"`
}

func (es EnvSource) Load() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var cfg Config

	err = envconfig.Process(es.Prefix, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
