package mux

import (
	"net/http"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/handlers"
	"github.com/squeakycheese75/cheese-grater/internal/http/middleware"
)

func NewRouter(cfg config.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.Handle("/", middleware.AuthWithAPIKey(handlers.ProxyHandler(), cfg))

	return mux, nil
}
