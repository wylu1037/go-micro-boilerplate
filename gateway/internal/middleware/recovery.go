package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Recovered from panic",
						zap.Any("panic", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("stack", string(debug.Stack())),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error","code":500}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
