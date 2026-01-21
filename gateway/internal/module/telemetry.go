package module

import (
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	pkgconfig "github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/telemetry"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

func NewTracer(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*trace.TracerProvider, error) {
	return telemetry.NewTracerProvider(lc, &pkgconfig.Config{
		Service: pkgconfig.ServiceConfig{
			Name:    cfg.Service.Name,
			Version: cfg.Service.Version,
			Env:     cfg.Service.Env,
		},
		Telemetry: pkgconfig.TelemetryConfig{
			Endpoint: cfg.Telemetry.Endpoint,
			Sampling: cfg.Telemetry.Sampling,
		},
	})
}
