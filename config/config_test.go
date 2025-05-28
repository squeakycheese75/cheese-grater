package config

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantError bool
	}{
		{
			name: "valid config",
			cfg: Config{
				APIKey:      "some-api-key",
				ProxyPort:   8080,
				RedirectURL: "http://localhost/callback",
			},
			wantError: false,
		},
		{
			name: "missing APIKey",
			cfg: Config{
				APIKey:      "",
				ProxyPort:   8080,
				RedirectURL: "http://localhost/callback",
			},
			wantError: true,
		},
		{
			name: "missing ProxyPort",
			cfg: Config{
				APIKey:      "some-api-key",
				ProxyPort:   0,
				RedirectURL: "http://localhost/callback",
			},
			wantError: true,
		},
		{
			name: "missing RedirectURL",
			cfg: Config{
				APIKey:      "some-api-key",
				ProxyPort:   8080,
				RedirectURL: "",
			},
			wantError: true,
		},
		{
			name: "all fields missing",
			cfg: Config{
				APIKey:      "",
				ProxyPort:   0,
				RedirectURL: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}
