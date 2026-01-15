package middleware

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"go-micro.dev/v4/server"
)

// NewLoggingMiddleware returns a go-micro server.HandlerWrapper that logs requests.
func NewLoggingMiddleware(logger *zerolog.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) error {
			start := time.Now()

			userId := lo.TernaryF(ctx.Value("userId") != nil, func() string {
				return ctx.Value("userId").(string)
			}, func() string {
				return ""
			})

			// Log request start
			logger.Info().
				Str("requestId", req.Header()["X-Request-Id"]).
				Str("userId", userId).
				Str("service", req.Service()).
				Str("endpoint", req.Endpoint()).
				Str("method", req.Method()).
				Msg("request started")

			// Call the handler
			err := fn(ctx, req, rsp)

			// Log request completion
			duration := time.Since(start)
			event := logger.Info().
				Str("requestId", req.Header()["X-Request-Id"]).
				Str("userId", userId).
				Str("service", req.Service()).
				Str("endpoint", req.Endpoint()).
				Str("method", req.Method()).
				Dur("duration", duration)

			if err != nil {
				event.Err(err).Msg("request failed")
			} else {
				event.Msg("request completed")
			}

			return err
		}
	}
}
