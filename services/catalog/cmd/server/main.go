package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/infra"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/bootstrap"
	catalog "github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/module"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(infra.NewConfig("catalog", "catalog")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		infra.Module,
		catalog.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
