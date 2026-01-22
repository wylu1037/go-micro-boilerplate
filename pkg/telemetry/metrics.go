package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewMeterProvider initializes an OTLP meter provider and sets the global meter provider.
func NewMeterProvider(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	endpoint := cfg.Telemetry.Endpoint
	if endpoint == "" {
		if cfg.Service.Env == "dev" {
			endpoint = "localhost:4317"
		} else {
			return nil, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

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

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(15*time.Second))),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	// Start runtime metrics collection
	if err := runtime.Start(runtime.WithMeterProvider(mp)); err != nil {
		return nil, fmt.Errorf("failed to start runtime metrics: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return mp.Shutdown(ctx)
		},
	})

	return mp, nil
}
