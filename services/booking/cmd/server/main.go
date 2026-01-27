package main

import (
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/wylu1037/go-micro-boilerplate/pkg/infra"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/bootstrap"
	booking "github.com/wylu1037/go-micro-boilerplate/services/booking/internal/module"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(infra.NewConfig("booking", "ticketing.booking")),
		fx.Provide(bootstrap.NewMicroService),
		fx.Provide(router.NewRouter),
		infra.Module,
		booking.Module,
		fx.Invoke(bootstrap.Start),
	).Run()
}
