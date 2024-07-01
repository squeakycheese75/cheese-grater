package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/router"
)

func main() {
	rootLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(rootLogger)

	slog.Debug("Loading config ...")
	env := config.EnvSource{}
	cfg, err := env.Load()
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Starting 'cheese-grater' on port: %d redirecting to: %v", cfg.ProxyPort, cfg.RedirectURL))

	router.Route(cfg)
}
