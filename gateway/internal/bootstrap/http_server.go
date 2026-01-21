package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	"go-micro.dev/v4"
	"go-micro.dev/v4/api/handler"
	"go-micro.dev/v4/api/handler/rpc"
	"go-micro.dev/v4/api/router"
	"go-micro.dev/v4/api/router/registry"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/middleware"
)

func NewHTTPServer(
	cfg *config.Config,
	logger *zap.Logger,
	microService micro.Service,
) *http.Server {
	r := chi.NewRouter()
	r.Use(otelchi.Middleware(cfg.Service.Name, otelchi.WithChiRoutes(r)))

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

	// Custom middleware to improve Span Name
	// Default otelchi uses route pattern which is often just /api/* for go-micro wrapper
	// Pre-compile regexes for performance & safety
	uuidRegex := regexp.MustCompile(`/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}(/|$)`)
	mongoIDRegex := regexp.MustCompile(`/[a-f0-9]{24}(/|$)`)
	numIDRegex := regexp.MustCompile(`/\d+(/|$)`)

	spanNameFormatter := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.SpanFromContext(r.Context())
			if span.IsRecording() {
				path := r.URL.Path

				// Normalize variations
				// Order matters: specific structures first
				path = uuidRegex.ReplaceAllString(path, "/{uuid}$1")
				path = mongoIDRegex.ReplaceAllString(path, "/{id}$1")
				path = numIDRegex.ReplaceAllString(path, "/{id}$1")

				span.SetName(fmt.Sprintf("%s %s", r.Method, path))
			}
			next.ServeHTTP(w, r)
		})
	}

	r.Route("/api", func(router chi.Router) {
		router.Use(spanNameFormatter)
		router.Use(middleware.TraceContextInjector) // Bridge OTel context to go-micro metadata
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
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting API Gateway",
					zap.String("name", cfg.Service.Name),
					zap.String("version", cfg.Service.Version),
					zap.String("address", cfg.Service.Address),
				)

				logger.Info("HTTP server listening", zap.String("address", cfg.Service.Address))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start gateway server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down gateway server...")
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Gateway server forced to shutdown", zap.Error(err))
				return err
			}
			logger.Info("Gateway server exited")
			return nil
		},
	})
}
