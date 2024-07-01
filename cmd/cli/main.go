package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/lithammer/shortuuid/v4"
	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/internal/http/router"
)

func main() {
	var (
		redirectURL string
		port        int
		apiKey      string
	)
	uid := shortuuid.New()

	flag.StringVar(&redirectURL, "RedirectURL", "localhost:1234", "the URL of the LM Studio Server")
	flag.IntVar(&port, "Port", 8080, "The Port to run this Redirector on")
	flag.StringVar(&apiKey, "APIKey", uid, "Need to set an API key")

	help := flag.Bool("help", false, "Display this help message")
	generatedAPIKey := false

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if apiKey == uid {
		generatedAPIKey = true
	}

	rootLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(rootLogger)

	slog.Info(fmt.Sprintf("Starting 'cheese-grater' on port: %d redirecting to: %v", port, redirectURL))
	if generatedAPIKey {
		slog.Info(fmt.Sprintf("Generated Password '%v', you will need to copy this into the 'Cursor settings - Models - OpenAI key'", apiKey))
	}

	router.Route(config.Config{
		ProxyPort:   port,
		RedirectURL: redirectURL,
		APIKey:      apiKey,
	})
}
