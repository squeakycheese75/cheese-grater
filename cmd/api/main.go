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

	HeaderAuthorization    = "Authorization"
	HeaderContentType      = "Content-Type"
	HeaderXForwardedHost   = "X-Forwarded-Host"
	ContentTypeJSON        = "application/json"
	ContentTypeEventStream = "text/event-stream"
	MethodOptions          = "OPTIONS"
	BearerPrefix           = "Bearer "
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
	if r.Method == MethodOptions {
		setCORSHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	bearerToken := r.Header.Get(HeaderAuthorization)
	if !strings.HasPrefix(bearerToken, BearerPrefix) || bearerToken[len(BearerPrefix):] != ValidAPIKey {
		http.Error(w, "Unauthorized: Missing or invalid API-Key", http.StatusUnauthorized)
		return
	}

	reverseProxy := createReverseProxy(RedirectURL)
	reverseProxy.ServeHTTP(w, r)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,PATCH,DELETE,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, User-Agent, X-Api-Key, X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, HTTP-Referer, X-Windowai-Title, X-Openrouter-Title, X-Title, X-Stainless-Lang, X-Stainless-Package-Version, X-Stainless-OS, X-Stainless-Arch, X-Stainless-Runtime, X-Stainless-Runtime-Version")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

func createReverseProxy(target string) *http.ServeMux {
	proxy := http.NewServeMux()
	proxy.Handle("/", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, target)
	})))
	return proxy
}

func proxyRequest(w http.ResponseWriter, r *http.Request, target string) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second, // Set the timeout duration here
	}

	r.URL.Host = target
	r.URL.Scheme = "http"
	r.Header.Set(HeaderXForwardedHost, r.Header.Get("Host"))
	r.Host = target

	req, err := createProxyRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error proxying request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyResponseHeaders(w, resp)

	if resp.Header.Get(HeaderContentType) == ContentTypeEventStream {
		setEventStreamHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
	} else {
		w.WriteHeader(resp.StatusCode)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error streaming response body: %v", err)
	}
}

func createProxyRequest(r *http.Request) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	req.Header.Set(HeaderXForwardedHost, r.Header.Get("Host"))
	return req, nil
}

func copyResponseHeaders(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

func setEventStreamHeaders(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, ContentTypeEventStream)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func main() {
	http.HandleFunc("/", proxyHandler)

	if err := http.ListenAndServe(ProxyEndpoint, loggingMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
