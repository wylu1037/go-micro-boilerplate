package infra

import (
	"github.com/wylu1037/go-micro-boilerplate/pkg/telemetry"
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
		telemetry.NewTracerProvider,
	),
	// Force initialization of TracerProvider to ensure global tracer is set
	fx.Invoke(func(_ *sdktrace.TracerProvider) {}),
)
