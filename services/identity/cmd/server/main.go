package main

import (
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/bootstrap"
	identityprovider "github.com/wylu1037/go-micro-boilerplate/services/identity/internal/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/router"
)

func main() {
	fx.New(
		fx.Provide(provider.NewConfig("identity", "identity")),
		fx.Provide(bootstrap.NewGRPCServer),
		fx.Provide(router.NewRouter),
		provider.InfraModule,
		identityprovider.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
