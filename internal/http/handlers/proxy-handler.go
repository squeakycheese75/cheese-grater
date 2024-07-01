package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	RedirectURL          = "localhost:1234"
	HeaderContentType    = "Content-Type"
	HeaderXForwardedHost = "X-Forwarded-Host"

	ContentTypeJSON        = "application/json"
	ContentTypeEventStream = "text/event-stream"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {

	reverseProxy := createReverseProxy(RedirectURL)
	reverseProxy.ServeHTTP(w, r)
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
		Timeout: 30 * time.Second,
	}

	r.URL.Host = target
	r.URL.Scheme = "http"
	r.Header.Set(HeaderXForwardedHost, r.Header.Get("Host"))
	r.Host = target

	for k, v := range r.Header {
		log.Printf("Header field %q, Value %q\n", k, v)
	}

	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set(HeaderContentType, ContentTypeJSON)
	req.Header.Set(HeaderXForwardedHost, r.Header.Get("Host"))

	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error proxying request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if resp.Header.Get(HeaderContentType) == ContentTypeEventStream {
		w.Header().Set(HeaderContentType, ContentTypeEventStream)
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
	} else {
		w.WriteHeader(resp.StatusCode)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error streaming response body: %v", err)
	}
}
