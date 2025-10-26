package main

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"
)

var validCode = regexp.MustCompile(`^\d{4,6}$`)

type contextCodeKey struct{}

// validateCodeMiddleware validates the code parameter from the URL path
func validateCodeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		if !validCode.MatchString(code) {
			http.Error(w, "404 page not found", http.StatusNotFound)
			slog.Error("invalid code", "code", code)
			return
		}
		// Store validated code in context
		ctx := context.WithValue(r.Context(), contextCodeKey{}, code)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// cacheMiddleware caches responses based on the code parameter
func cacheMiddleware(cache *cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			code, ok := r.Context().Value(contextCodeKey{}).(string)
			if !ok {
				http.Error(w, "code not found in context", http.StatusInternalServerError)
				slog.Error("code not found in context")
				return
			}

			// Check cache first
			if cached, ok := cache.get(code); ok {
				slog.Info("cache hit", "code", code)
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write(cached); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					slog.Error("write cached response", "error", err)
					return
				}
				return
			}

			slog.Info("cache miss", "code", code)

			// Create a response writer wrapper to capture the response
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				code:           code,
				cache:          cache,
			}

			next.ServeHTTP(wrapper, r)
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to cache successful responses
type responseWriterWrapper struct {
	http.ResponseWriter
	code       string
	cache      *cache
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	// If WriteHeader was not called, status is 200 OK by default
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	// Only cache successful responses
	if w.statusCode < 300 {
		w.cache.set(w.code, b)
		slog.Info("response cached", "code", w.code, "status", w.statusCode)
	}

	return w.ResponseWriter.Write(b)
}
