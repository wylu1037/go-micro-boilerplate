package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/bootstrap"
	catalogprovider "github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/provider"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(provider.NewConfig("catalog", "catalog")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		provider.InfraModule,
		catalogprovider.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
