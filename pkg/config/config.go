package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	NATS     NATSConfig     `mapstructure:"nats"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Address string `mapstructure:"address"`
	Env     string `mapstructure:"env"` // dev, staging, prod
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	Schema       string        `mapstructure:"schema"`
	SSLMode      string        `mapstructure:"ssl_mode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

// DSN returns the database connection string
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode, c.Schema,
	)
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// NATSConfig holds NATS configuration
type NATSConfig struct {
	URL string `mapstructure:"url"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Issuer          string        `mapstructure:"issuer"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
}

// Load loads configuration from file and environment variables
func Load(serviceName string) (*Config, error) {
	v := viper.New()

	// Set config file name
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	v.AddConfigPath("./config")
	v.AddConfigPath(fmt.Sprintf("./services/%s/config", serviceName))
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Read from environment variables
	v.SetEnvPrefix("TICKETING")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
