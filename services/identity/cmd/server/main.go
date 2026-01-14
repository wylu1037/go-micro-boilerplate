package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/bootstrap"
	identityprovider "github.com/wylu1037/go-micro-boilerplate/services/identity/internal/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(provider.NewConfig("identity", "identity")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		provider.InfraModule,
		identityprovider.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
