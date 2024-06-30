package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/middleware"
	"github.com/squeakycheese75/cheese-grater/internal/http/mux"
)

func Start(cfg config.Config) error {
	slog.Info("Starting cheese-grater...")
	mux, err := mux.NewRouter(cfg)
	if err != nil {
		slog.Error(errors.Wrap(err, "server start error").Error())
		return err
	}

	serverAddress := fmt.Sprintf(":%d", cfg.ProxyPort)
	server := &http.Server{
		Addr:    serverAddress,
		Handler: middleware.LoggingMiddleware(mux),
	}

	go func() {
		slog.Info("listening at : " + fmt.Sprint(cfg.ProxyPort))
		slog.Info("redirecting to : " + fmt.Sprint(cfg.RedirectURL))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(errors.Wrap(err, "listenAndServe error").Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Received request to shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	slog.Warn("Server exiting ...")

	return nil
}
