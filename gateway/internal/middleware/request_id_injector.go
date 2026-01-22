package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestIDInjector extracts the request ID from Chi context and injects it into HTTP headers.
// This allows go-micro handlers to access the same request ID for distributed tracing.
func RequestIDInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := chimiddleware.GetReqID(r.Context()); reqID != "" {
			r.Header.Set("X-Request-Id", reqID)
		}
		next.ServeHTTP(w, r)
	})
}
