package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/pkg/tools"
)

func Logging(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			requestID := chimiddleware.GetReqID(r.Context())
			traceID, spanID := tools.ExtractTraceInfo(r.Context())

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			logger.Info("Gateway HTTP request",
				zap.String("trace_id", traceID),
				zap.String("span_id", spanID),
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}
