package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
)

func NewHTTPServer(
	cfg *config.Config,
	logger *zerolog.Logger,
) *http.Server {
	ctx := context.Background()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register Identity service handler
	identityAddr := cfg.Backends["identity"].Address
	if err := identityv1.RegisterIdentityServiceHandlerFromEndpoint(ctx, mux, identityAddr, opts); err != nil {
		logger.Fatal().Err(err).Str("address", identityAddr).Msg("Failed to register Identity service handler")
	}
	logger.Info().Str("address", identityAddr).Msg("Registered Identity service")

	server := &http.Server{
		Addr:         cfg.Service.Address,
		Handler:      corsMiddleware(mux),
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-Id")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
