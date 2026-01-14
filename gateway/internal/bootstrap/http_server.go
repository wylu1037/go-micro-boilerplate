package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/middleware"
	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
)

func NewHTTPServer(
	cfg *config.Config,
	logger *zerolog.Logger,
) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.Recovery(logger))
	r.Use(chimiddleware.RequestID)
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

	r.Get("/health", healthCheckHandler)
	r.Get("/ready", readinessCheckHandler())
	r.Get("/version", versionHandler(cfg))

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RateLimiter(100, 200))

		gwMux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		)

		ctx := context.Background()
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		// 注册 Identity service
		identityAddr := cfg.Backends["identity"].Address
		if err := identityv1.RegisterIdentityServiceHandlerFromEndpoint(ctx, gwMux, identityAddr, opts); err != nil {
			logger.Fatal().Err(err).Str("address", identityAddr).Msg("Failed to register Identity service handler")
		}
		logger.Info().Str("address", identityAddr).Msg("Registered Identity service")

		// Mount grpc-gateway to /api/*
		r.Mount("/", gwMux)
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

func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization", "X-Request-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func readinessCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	}
}

func versionHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{"name":"` + cfg.Service.Name + `","version":"` + cfg.Service.Version + `","env":"` + cfg.Service.Env + `"}`
		w.Write([]byte(response))
	}
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
					logger.Fatal().Err(err).Msg("Failed to start server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Shutting down server...")
			if err := server.Shutdown(ctx); err != nil {
				logger.Error().Err(err).Msg("Server forced to shutdown")
				return err
			}
			logger.Info().Msg("Server exited")
			return nil
		},
	})
}
