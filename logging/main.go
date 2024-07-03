package logging

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	rootLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(rootLogger)
}
