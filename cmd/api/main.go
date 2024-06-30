package main

import (
	"log/slog"
	"os"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/server"
	"github.com/squeakycheese75/cheese-grater/logging"
)

func loadConfig() (*config.Config, error) {
	slog.Debug("Loading config ...")
	env := config.EnvSource{}
	cfg, err := env.Load()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func run() error {
	logging.SetupLogger()

	cfg, err := loadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", slog.String("error", err.Error()))
		return err
	}

	if err := server.Start(*cfg); err != nil {
		slog.Error("Error starting server", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("Application failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
