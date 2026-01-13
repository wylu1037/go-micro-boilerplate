package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

var defaultLogger *zap.Logger

// Init initializes the global logger
func Init(cfg config.LogConfig) error {
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		logLevel = zapcore.InfoLevel
	}

	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if cfg.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// Get returns the default logger
func Get() *zap.Logger {
	if defaultLogger == nil {
		defaultLogger, _ = zap.NewProduction()
	}
	return defaultLogger
}

// Sugar returns a sugared logger
func Sugar() *zap.SugaredLogger {
	return Get().Sugar()
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Get().Sync()
}
