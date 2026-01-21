package bootstrap

import (
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/rs/zerolog"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	"go-micro.dev/v4"
)

func NewMicroService(
	cfg *config.Config,
	logger *zerolog.Logger,
) micro.Service {
	service := micro.NewService(
		micro.Name(cfg.Service.Name),
		micro.Version(cfg.Service.Version),
		micro.Address(cfg.Service.Address),
		micro.WrapClient(opentelemetry.NewClientWrapper()),
	)

	service.Init() // Parse command line flags and environment variables

	return service
}
