package bootstrap

import (
	"context"

	"github.com/rs/zerolog"
	"go-micro.dev/v4"
	"go-micro.dev/v4/auth"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/logger"
	"github.com/wylu1037/go-micro-boilerplate/pkg/middleware"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/router"
)

func NewMicroService(
	logger *zerolog.Logger,
	cfg *config.Config,
	microAuth auth.Auth,
) micro.Service {
	service := micro.NewService(
		micro.Name(cfg.Service.Name),
		micro.Version(cfg.Service.Version),
		micro.Address(cfg.Service.Address),
		micro.Auth(microAuth),
		micro.WrapHandler(
			middleware.NewRecoveryMiddleware(),
			middleware.AuthWrapper(microAuth, []string{}),
			middleware.NewLoggingMiddleware(logger),
			middleware.NewValidatorMiddleware(),
		),
	)

	service.Init()

	return service
}

type MicroServiceParams struct {
	fx.In

	Lifecycle    fx.Lifecycle
	Config       *config.Config
	Logger       *zerolog.Logger
	MicroService micro.Service
	Router       router.Router
}

func Start(p MicroServiceParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				p.Logger.Info().
					Str("name", p.Config.Service.Name).
					Str("version", p.Config.Service.Version).
					Str("address", p.Config.Service.Address).
					Msg("Starting Catalog Micro service")

				p.Router.Register()

				if err := p.MicroService.Run(); err != nil {
					p.Logger.Fatal().Err(err).Msg("Catalog Micro service failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Stopping Catalog Micro service")
			return nil
		},
	})
}
