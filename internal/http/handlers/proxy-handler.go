package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/squeakycheese75/cheese-grater/entities"
)

func ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := r.Context().Value(entities.RedirectURL).(string)

		reverseProxy := http.NewServeMux()
		reverseProxy.Handle("/", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpClient := &http.Client{
				Timeout: 10 * time.Second,
			}

			r.URL.Host = redirectURL
			r.URL.Scheme = "http"
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = redirectURL

			for k, v := range r.Header {
				log.Printf("Header field %q, Value %q\n", k, v)
			}

			req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

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

			if resp.Header.Get("Content-Type") == "text/event-stream" {
				w.Header().Set("Content-Type", "text/event-stream")
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
		})))
		reverseProxy.ServeHTTP(w, r)
	}
}
