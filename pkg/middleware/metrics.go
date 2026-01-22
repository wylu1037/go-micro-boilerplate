package middleware

import (
	"context"
	"time"

	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NewMetricsMiddleware returns a go-micro server.HandlerWrapper that records metrics.
// It records:
// - rpc_server_request_total: Counter of total requests
// - rpc_server_request_duration_seconds: Histogram of request duration
func NewMetricsMiddleware() server.HandlerWrapper {
	meter := otel.GetMeterProvider().Meter("rpc_server")

	requestCounter, _ := meter.Int64Counter(
		"rpc_server_request_total",
		metric.WithDescription("Total number of RPC requests"),
		metric.WithUnit("1"),
	)

	requestDuration, _ := meter.Float64Histogram(
		"rpc_server_request_duration_seconds",
		metric.WithDescription("Duration of RPC requests"),
		metric.WithUnit("s"),
	)

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) error {
			start := time.Now()

			// Call handler
			err := fn(ctx, req, rsp)

			duration := time.Since(start).Seconds()

			attrs := metric.WithAttributes(
				attribute.String("service", req.Service()),
				attribute.String("method", req.Endpoint()),
				attribute.String("status", getStatus(err)),
			)

			requestCounter.Add(ctx, 1, attrs)
			requestDuration.Record(ctx, duration, attrs)

			return err
		}
	}
}

func getStatus(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}
