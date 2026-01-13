package main

import (
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/bootstrap"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/provider"
)

func main() {
	fx.New(
		fx.Provide(provider.NewConfig()),
		fx.Provide(provider.NewLogger),
		fx.Provide(bootstrap.NewHTTPServer),
		fx.Invoke(bootstrap.Start),
	).Run()
}
