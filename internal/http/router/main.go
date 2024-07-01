package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/handlers"
	"github.com/squeakycheese75/cheese-grater/internal/http/middleware"
)

func Route(cfg config.Config) {
	http.HandleFunc("/", handlers.ProxyHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ProxyPort), middleware.LoggingMiddleware(http.DefaultServeMux)); err != nil {
		slog.Info("Failed to start server: %v", err)
		panic(err)
	}
}
