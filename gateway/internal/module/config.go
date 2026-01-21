package module

import (
	"github.com/wylu1037/go-micro-boilerplate/gateway/internal/config"
)

func NewConfig() func() (*config.Config, error) {
	return func() (*config.Config, error) {
		return config.Load()
	}
}
