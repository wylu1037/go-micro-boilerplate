package provider

import (
	"github.com/rs/zerolog"

	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
	pkgconfig "github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/logger"
)

func NewLogger(cfg *config.Config) (*zerolog.Logger, error) {
	if err := logger.Init(pkgconfig.LogConfig{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	}); err != nil {
		return nil, err
	}

	return logger.Get(), nil
}
