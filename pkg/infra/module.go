package infra

import (
	"github.com/wylu1037/go-micro-boilerplate/pkg/telemetry"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

// Module provides common infrastructure dependencies (logger, database, cache, auth)
// Note: Config must be provided separately by each service using NewConfig(serviceName, schema)
var Module = fx.Options(
	fx.Provide(
		NewLogger,
		NewDatabase,
		NewMicroAuth,
		NewRedis,
		NewEtcd,
		NewDistributedLocker,
		telemetry.NewLoggerProvider,
		telemetry.NewTracerProvider,
	),
	// Force initialization of providers to ensure global tracer/logger are set
	fx.Invoke(func(_ *sdktrace.TracerProvider, _ *sdklog.LoggerProvider) {}),
)
