package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/infra"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/bootstrap"
	notification "github.com/wylu1037/go-micro-boilerplate/services/notification/internal/module"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(infra.NewConfig("notification", "ticketing.notification")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		infra.Module,
		notification.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
