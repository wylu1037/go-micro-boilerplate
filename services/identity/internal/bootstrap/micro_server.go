package bootstrap

import (
	"context"

	"github.com/rs/zerolog"
	"go-micro.dev/v4"
	"go-micro.dev/v4/auth"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/middleware"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/router"
)

func NewMicroService(
	cfg *config.Config,
	logger *zerolog.Logger,
	microAuth auth.Auth,
) micro.Service {
	service := micro.NewService(
		micro.Name(cfg.Service.Name),
		micro.Version(cfg.Service.Version),
		micro.Address(cfg.Service.Address),
		micro.Auth(microAuth),
		micro.WrapHandler(
			middleware.NewRecoveryMiddleware(),
			middleware.AuthWrapper(microAuth, []string{
				"IdentityService.Register",
				"IdentityService.Login",
				"IdentityService.RefreshToken",
			}),
			middleware.NewLoggingMiddleware(logger),
			middleware.NewValidatorMiddleware(),
		),
	)

	service.Init() // Parse command line flags and environment variables

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
					Msg("Starting Identity Micro service")

				p.Router.Register()

				if err := p.MicroService.Run(); err != nil {
					p.Logger.Fatal().Err(err).Msg("Identity Micro service failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info().Msg("Stopping Identity Micro service")
			return nil
		},
	})
}
