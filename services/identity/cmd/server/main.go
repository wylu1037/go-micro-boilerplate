package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/infra"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/bootstrap"
	identity "github.com/wylu1037/go-micro-boilerplate/services/identity/internal/module"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(infra.NewConfig("identity", "identity")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		infra.Module,
		identity.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
