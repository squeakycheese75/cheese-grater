package router

import (
	"net/http"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/handlers"
	"github.com/squeakycheese75/cheese-grater/internal/http/middleware"
)

func NewRouter(cfg config.Config) error {
	// Initialize the handler with middleware
	proxyHandler := handlers.ProxyHandler()
	authenticatedHandler := middleware.AuthWithAPIKey(proxyHandler, cfg)

	// Set up the HTTP server
	http.Handle("/", authenticatedHandler)

	// serverAddress := fmt.Sprintf(":%d", cfg.ProxyPort)
	// server := &http.Server{
	// 	Addr:    serverAddress,
	// 	Handler: middleware.LoggingMiddleware(http.DefaultServeMux),
	// }

	// Start the HTTP server
	// log.Printf("Starting server on %s", serverAddress)
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }

	// go func() {
	// 	// logger.Debug("listening at : " + fmt.Sprint(conf.AppPort))

	// 	// if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
	// 	// 	logger.Error(errors.Wrap(err, "listenAndServe error").Error())
	// 	// 	panic(fmt.Sprintf("listen: %s\n", err))
	// 	// }
	// 	// Start the HTTP server
	// 	log.Printf("Starting server on %s", serverAddress)
	// 	if err := server.ListenAndServe(); err != nil {
	// 		log.Fatalf("Failed to start server: %v", err)
	// 	}
	// }()
	return nil
}
