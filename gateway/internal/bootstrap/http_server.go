package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"go-micro.dev/v4"
	"go-micro.dev/v4/api/handler"
	"go-micro.dev/v4/api/handler/rpc"
	"go-micro.dev/v4/api/router"
	"go-micro.dev/v4/api/router/registry"
	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/middleware"
)

func NewHTTPServer(
	cfg *config.Config,
	logger *zerolog.Logger,
	microService micro.Service,
) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.Recovery(logger))
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestIDInjector)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: false,
		MaxAge:           3600,
	}))
	r.Use(chimiddleware.Timeout(60 * time.Second))

	microRouter := registry.NewRouter(
		router.WithRegistry(microService.Options().Registry),
	)
	microHandler := rpc.NewHandler(
		handler.WithClient(microService.Client()),
		handler.WithRouter(microRouter),
	)
	r.Route("/api", func(router chi.Router) {
		router.Use(middleware.RateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst))
		router.Mount("/", microHandler)
	})

	server := &http.Server{
		Addr:         cfg.Service.Address,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func Start(
	lc fx.Lifecycle,
	server *http.Server,
	cfg *config.Config,
	logger *zerolog.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info().
					Str("name", cfg.Service.Name).
					Str("version", cfg.Service.Version).
					Str("address", cfg.Service.Address).
					Msg("Starting API Gateway")

				logger.Info().Str("address", cfg.Service.Address).Msg("HTTP server listening")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal().Err(err).Msg("Failed to start gateway server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Shutting down gateway server...")
			if err := server.Shutdown(ctx); err != nil {
				logger.Error().Err(err).Msg("Gateway server forced to shutdown")
				return err
			}
			logger.Info().Msg("Gateway server exited")
			return nil
		},
	})
}
