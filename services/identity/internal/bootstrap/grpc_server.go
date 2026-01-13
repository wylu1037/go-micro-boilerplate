package bootstrap

import (
	"context"
	"net"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/middleware"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/router"
)

func NewGRPCServer(cfg *config.Config, logger *zerolog.Logger) *grpc.Server {
	server := grpc.NewServer(
		middleware.RegisterInterceptors(
			[]grpc.UnaryServerInterceptor{
				middleware.NewLoggingInterceptor(logger),
				middleware.NewRecoveryInterceptor(),
				middleware.NewValidatorInterceptor(),
			}),
	)

	return server
}

type GRPCServerParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     *config.Config
	Logger     *zerolog.Logger
	Router     router.Router
	GRPCServer *grpc.Server
}

func Start(p GRPCServerParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", p.Config.Service.Address)
			if err != nil {
				return err
			}

			p.Router.Register()

			p.Logger.Info().
				Str("name", p.Config.Service.Name).
				Str("address", p.Config.Service.Address).
				Msg("Starting gRPC server")

			go func() {
				if err := p.GRPCServer.Serve(lis); err != nil {
					p.Logger.Fatal().Err(err).Msg("gRPC server failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info().Msg("Stopping gRPC server")
			p.GRPCServer.GracefulStop()
			return nil
		},
	})
}
