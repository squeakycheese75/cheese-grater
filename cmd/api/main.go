package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Constants for API key validation
const (
	ValidAPIKey   = "my-secret-api-key"
	RedirectURL   = "localhost:1234" // Replace with your redirect server URL
	ProxyEndpoint = ":8080"          // Port on which the proxy server will run
)

// Middleware for logging incoming and outgoing URLs
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming URL
		log.Printf("Incoming request URL: %s", r.URL)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Handler function for the proxy server
func proxyHandler(w http.ResponseWriter, r *http.Request) {

	// Don't check for a key on the Optios
	if r.Method == http.MethodOptions {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                                                                                                                                                                                                                                                                                             // Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,PATCH,DELETE,POST,PUT")                                                                                                                                                                                                                                                                                                                            // Allow GET, POST, and OPTIONS methods
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, User-Agent, X-Api-Key, X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, HTTP-Referer, X-Windowai-Title, X-Openrouter-Title, X-Title, X-Stainless-Lang, X-Stainless-Package-Version, X-Stainless-OS, X-Stainless-Arch, X-Stainless-Runtime, X-Stainless-Runtime-Version") // Allow Content-Type and API-Key headers
		w.Header().Set("Access-Control-Max-Age", "86400")                                                                                                                                                                                                                                                                                                                                                              // Cache preflight response for 24 hours
		w.WriteHeader(http.StatusOK)
		return
	}

	bearerToken := r.Header.Get("Authorization")

	tokenString := ""

	if strings.HasPrefix(bearerToken, "Bearer ") {
		tokenString = bearerToken[7:]
	}

	if tokenString == string(ValidAPIKey) {
		reverseProxy := createReverseProxy(RedirectURL)
		reverseProxy.ServeHTTP(w, r)
	}

	http.Error(w, "Unauthorized: Missing or invalid API-Key", http.StatusUnauthorized)
}

// CreateReverseProxy creates a reverse proxy to the given target URL
func createReverseProxy(target string) *http.ServeMux {
	client := &http.Client{
		Timeout: 30 * time.Second, // Set the timeout duration here
	}

	proxy := http.NewServeMux()
	proxy.Handle("/", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Modify the request as needed (optional)
		r.URL.Host = target
		r.URL.Scheme = "http"
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target

		resp, err := client.Post("http://localhost:1234/v1/chat/completions", "application/json", r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error proxying request: %v", err), http.StatusBadGateway)
			return
		}

		fmt.Println(resp.Header.Get("Content-Type"))
		if resp.Header.Get("Content-Type") == "text/event-stream" {
			// Set SSE headers
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()

			// Stream the response body
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				log.Printf("Error streaming response body: %v", err)
			}
		} else {
			// Copy response headers
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			// Copy status code
			w.WriteHeader(resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusBadGateway)
				return
			}
			_, err = w.Write(body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error writing response body: %v", err), http.StatusInternalServerError)
				return
			}
		}
	})))
	return proxy
}

func main() {
	http.HandleFunc("/", proxyHandler)

	http.ListenAndServe(ProxyEndpoint, loggingMiddleware(http.DefaultServeMux))
}
