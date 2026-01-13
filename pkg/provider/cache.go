package provider

import (
	"context"

	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/cache"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

func NewRedis(lc fx.Lifecycle, cfg *config.Config) (*cache.Client, error) {
	client, err := cache.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
