package infra

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"go.uber.org/fx"
)

func NewRedis(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*redis.Client, error) {
	if cfg.Redis.URL == "" {
		return nil, nil
	}

	opt, _ := redis.ParseURL("rediss://default:AVlDAAIncDI3MGExZGEwYjVlZmU0ZDkwYWRiM2U0ODIwNDhjNWRiZnAyMjI4NTE@civil-cobra-22851.upstash.io:6379")
	client := redis.NewClient(opt)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
