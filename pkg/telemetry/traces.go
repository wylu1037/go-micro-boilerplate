package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewTracerProvider initializes an OTLP tracer provider and sets the global tracer provider.
func NewTracerProvider(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	endpoint := cfg.Telemetry.Endpoint
	if endpoint == "" {
		if cfg.Service.Env == "dev" {
			endpoint = "localhost:4317"
		} else {
			return nil, nil
		}
	}

	// Set 5s timeout for connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create OTLP exporter
	// We use insecure here for internal communication, verify if TLS is needed
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithBlock()), // Wait for connection
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create Resource
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

	// Determine sampler
	var sampler sdktrace.Sampler
	if cfg.Telemetry.Sampling > 0 {
		sampler = sdktrace.TraceIDRatioBased(cfg.Telemetry.Sampling)
	} else {
		// Default to AlwaysSample for Dev, otherwise ParentBased(AlwaysOff) maybe?
		// Let's stick to AlwaysSample if 0/undefined in dev
		if cfg.Service.Env == "dev" {
			sampler = sdktrace.AlwaysSample()
		} else {
			sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1)) // Default 10% in non-dev?
		}
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global provider
	otel.SetTracerProvider(tp)

	// Set global propagator to W3C Trace Context (standard for distributed tracing)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Register cleanup hook
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return tp.Shutdown(ctx)
		},
	})

	return tp, nil
}
