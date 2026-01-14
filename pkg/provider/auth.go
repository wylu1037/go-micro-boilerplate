package provider

import (
	"github.com/wylu1037/go-micro-boilerplate/pkg/auth"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

func NewJWTManager(
	cfg *config.Config,
) *auth.JWTManager {
	return auth.NewJWTManager(cfg.JWT)
}
