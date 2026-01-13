package provider

import (
	"context"

	"go.uber.org/fx"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
)

func NewDatabase(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*db.Pool, error) {
	pool, err := db.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
