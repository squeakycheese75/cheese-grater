package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/squeakycheese75/cheese-grater/config"
	"github.com/squeakycheese75/cheese-grater/entities"
)

const (
	HeaderAuthorization = "Authorization"
	BearerPrefix        = "Bearer "
)

// Middleware for logging incoming and outgoing URLs
func AuthWithAPIKey(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			setCORSHeaders(w)
			w.WriteHeader(http.StatusOK)
			return
		}

		bearerToken := r.Header.Get(HeaderAuthorization)
		if !strings.HasPrefix(bearerToken, BearerPrefix) || bearerToken[len(BearerPrefix):] != cfg.APIKey {
			http.Error(w, "Unauthorized: Missing or invalid API-Key", http.StatusUnauthorized)
			return
		}

		ctxRedirectURL := context.WithValue(r.Context(), entities.RedirectURL, cfg.RedirectURL)

		rWithRedirectURL := r.WithContext(ctxRedirectURL)

		next.ServeHTTP(w, rWithRedirectURL)
	})
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,PATCH,DELETE,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, User-Agent, X-Api-Key, X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, HTTP-Referer, X-Windowai-Title, X-Openrouter-Title, X-Title, X-Stainless-Lang, X-Stainless-Package-Version, X-Stainless-OS, X-Stainless-Arch, X-Stainless-Runtime, X-Stainless-Runtime-Version")
	w.Header().Set("Access-Control-Max-Age", "86400")
}
