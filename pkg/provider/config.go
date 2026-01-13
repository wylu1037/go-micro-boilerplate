package provider

import (
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewConfig creates a new config with the specified service name and schema
func NewConfig(serviceName, schema string) func() (*config.Config, error) {
	return func() (*config.Config, error) {
		cfg, err := config.Load(serviceName)
		if err != nil {
			return nil, err
		}
		cfg.Database.Schema = schema
		return cfg, nil
	}
}
