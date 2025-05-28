package config

import (
	"fmt"

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
	_ = godotenv.Load(".env") // Optional: log warning if you like

	var cfg Config
	if err := envconfig.Process(es.Prefix, &cfg); err != nil {
		return cfg, fmt.Errorf("error processing environment variables: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	var missing []string
	if c.APIKey == "" {
		missing = append(missing, "APIKey")
	}
	if c.ProxyPort == 0 {
		missing = append(missing, "ProxyPort")
	}
	if c.RedirectURL == "" {
		missing = append(missing, "RedirectURL")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing or invalid config fields: %v", missing)
	}
	return nil
}
