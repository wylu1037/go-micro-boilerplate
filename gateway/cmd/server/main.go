package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/bootstrap"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/module"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	fx.New(
		fx.Provide(module.NewConfig()),
		fx.Provide(module.NewLoggerProvider),
		fx.Provide(module.NewLogger),
		fx.Provide(module.NewTracer),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(bootstrap.NewHTTPServer),
		fx.Invoke(bootstrap.Start),
		fx.Invoke(func(_ *sdktrace.TracerProvider, _ *sdklog.LoggerProvider) {}),
	).Run()
}
