package provider

import (
	"github.com/rs/zerolog"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/logger"
)

func NewLogger(cfg *config.Config) (*zerolog.Logger, error) {
	if err := logger.Init(cfg.Log); err != nil {
		return nil, err
	}
	return logger.Get(), nil
}
