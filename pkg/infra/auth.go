package infra

import (
	"github.com/go-micro/plugins/v4/auth/jwt"
	"go-micro.dev/v4/auth"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewMicroAuth creates an authentication instance based on the go-micro JWT plugin.
// For the identity service, PrivateKey and PublicKey must be configured.
// For other services, only PublicKey is required for verification.
func NewMicroAuth(cfg *config.Config) auth.Auth {
	opts := []auth.Option{
		auth.PublicKey(cfg.JWT.PublicKey),
		auth.Namespace(cfg.JWT.Namespace),
	}

	if cfg.JWT.PrivateKey != "" {
		opts = append(opts, auth.PrivateKey(cfg.JWT.PrivateKey))
	}

	return jwt.NewAuth(opts...)
}
