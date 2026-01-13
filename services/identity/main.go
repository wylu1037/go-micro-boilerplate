package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-micro.dev/v5"
	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/pkg/cache"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/pkg/logger"
)

const serviceName = "ticketing.identity"

func main() {
	// Load configuration
	cfg, err := config.Load("identity")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Set schema for this service
	cfg.Database.Schema = "identity"

	// Initialize logger
	if err := logger.Init(cfg.Log); err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("Starting identity service", zap.String("version", cfg.Service.Version))

	// Initialize database
	ctx := context.Background()
	dbPool, err := db.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	// Initialize Redis
	redisClient, err := cache.NewClient(cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Create go-micro service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(cfg.Service.Version),
		micro.Address(cfg.Service.Address),
	)

	// Initialize the service
	service.Init()

	// TODO: Register gRPC handler
	// identityv1.RegisterIdentityServiceHandler(service.Server(), handler)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Info("Shutting down...")
	}()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal("service run failed", zap.Error(err))
	}
}
