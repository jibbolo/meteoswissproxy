package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func startHttpServer(port string) error {
	responseCache := newCache(5 * time.Minute)

	// Apply middleware chain: validate -> cache -> handler
	handlerWithMiddleware := validateCodeMiddleware(cacheMiddleware(responseCache)(handler()))

	mux := http.NewServeMux()
	mux.Handle("GET /{code}", handlerWithMiddleware)

	slog.Info("Server listening", "port", port)
	// start new server with custom timeout
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return server.ListenAndServe()
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get code from context (set by validateCodeMiddleware)
		code, ok := r.Context().Value(contextCodeKey{}).(string)
		if !ok {
			http.Error(w, "code not found in context", http.StatusInternalServerError)
			slog.Error("code not found in context")
			return
		}

		api, err := fetchAll(code)
		if err != nil {
			var apiErr MSError
			if errors.As(err, &apiErr) {
				http.Error(w, apiErr.Message, apiErr.Status)
				slog.Error("go meteoswiss error", "error", err)
				return
			}
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			slog.Error("can't fetch all", "error", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(api); err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			slog.Error("write response", "error", err)
			return
		}
	}
}
