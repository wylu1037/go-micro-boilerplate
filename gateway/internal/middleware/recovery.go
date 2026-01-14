package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

func Recovery(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().
						Interface("panic", err).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("stack", string(debug.Stack())).
						Msg("Recovered from panic")

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error","code":500}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
