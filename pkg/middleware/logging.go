package middleware

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/pkg/tools"
)

// NewLoggingMiddleware returns a go-micro server.HandlerWrapper that logs requests.
// It automatically extracts trace_id and span_id from the context for log correlation.
func NewLoggingMiddleware(logger *zap.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) error {
			start := time.Now()

			userID := lo.TernaryF(ctx.Value("userId") != nil, func() string {
				return ctx.Value("userId").(string)
			}, func() string {
				return ""
			})

			traceID, spanID := tools.ExtractTraceInfo(ctx)

			logger.Info("request started",
				zap.String("trace_id", traceID),
				zap.String("span_id", spanID),
				zap.String("request_id", req.Header()["X-Request-Id"]),
				zap.String("user_id", userID),
				zap.String("service", req.Service()),
				zap.String("endpoint", req.Endpoint()),
				zap.String("method", req.Method()),
			)

			err := fn(ctx, req, rsp)

			duration := time.Since(start)
			if err != nil {
				logger.Info("request failed",
					zap.String("trace_id", traceID),
					zap.String("span_id", spanID),
					zap.String("request_id", req.Header()["X-Request-Id"]),
					zap.String("user_id", userID),
					zap.String("service", req.Service()),
					zap.String("endpoint", req.Endpoint()),
					zap.String("method", req.Method()),
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				logger.Info("request completed",
					zap.String("trace_id", traceID),
					zap.String("span_id", spanID),
					zap.String("request_id", req.Header()["X-Request-Id"]),
					zap.String("user_id", userID),
					zap.String("service", req.Service()),
					zap.String("endpoint", req.Endpoint()),
					zap.String("method", req.Method()),
					zap.Duration("duration", duration),
				)
			}

			return err
		}
	}
}
