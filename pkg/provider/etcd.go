package provider

import (
	"context"
	"time"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
)

func NewEtcd(
	lc fx.Lifecycle,
	cfg *config.Config,
) (*clientv3.Client, error) {
	if len(cfg.Etcd.Endpoints) == 0 {
		return nil, nil
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    cfg.Etcd.Username,
		Password:    cfg.Etcd.Password,
	})
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cli.Close()
		},
	})

	return cli, nil
}
