package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewLoggerProvider initializes an OTLP log provider and sets the global log provider.
func NewLoggerProvider(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*sdklog.LoggerProvider, error) {
	ctx := context.Background()
	endpoint := cfg.Telemetry.Endpoint
	if endpoint == "" {
		if cfg.Service.Env == "dev" {
			endpoint = "localhost:4317"
		} else {
			// No telemetry endpoint configured for non-dev, skip log provider
			return nil, nil
		}
	}

	// Set 5s timeout for connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create OTLP log exporter
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Create Resource (reuse the same resource attributes as tracing)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.Service.Name),
			semconv.ServiceVersion(cfg.Service.Version),
			semconv.DeploymentEnvironment(cfg.Service.Env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create LoggerProvider
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)

	// Set global log provider
	global.SetLoggerProvider(lp)

	// Register cleanup hook
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return lp.Shutdown(ctx)
		},
	})

	return lp, nil
}
