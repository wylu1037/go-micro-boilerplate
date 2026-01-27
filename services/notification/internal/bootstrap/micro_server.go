package bootstrap

import (
	"context"

	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/middleware"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/router"
)

func NewMicroService(
	cfg *config.Config,
	logger *zap.Logger,
	microAuth auth.Auth,
) micro.Service {
	service := micro.NewService(
		micro.Name(cfg.Service.Name),
		micro.Version(cfg.Service.Version),
		micro.Address(cfg.Service.Address),
		micro.Auth(microAuth),
		micro.WrapHandler(
			opentelemetry.NewHandlerWrapper(), // Add Tracing
			middleware.NewMetricsMiddleware(), // Add Metrics
			middleware.NewRecoveryMiddleware(logger),
			middleware.AuthWrapper(microAuth, []string{}), // All notification routes potentially internal or as per logic
			middleware.NewLoggingMiddleware(logger),
			middleware.NewValidatorMiddleware(logger),
		),
	)

	service.Init()

	return service
}

type MicroServiceParams struct {
	fx.In

	Lifecycle    fx.Lifecycle
	Config       *config.Config
	Logger       *zap.Logger
	MicroService micro.Service
	Router       router.Router
}

func Start(p MicroServiceParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				p.Logger.Info("Starting Notification Micro service",
					zap.String("name", p.Config.Service.Name),
					zap.String("version", p.Config.Service.Version),
					zap.String("address", p.Config.Service.Address),
				)

				p.Router.Register()

				if err := p.MicroService.Run(); err != nil {
					p.Logger.Fatal("Notification Micro service failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			p.Logger.Info("Stopping Notification Micro service")
			return nil
		},
	})
}
