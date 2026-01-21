package infra

import (
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

// NewLogger creates a zap logger with OTLP bridge for sending logs to the collector.
// The logger outputs to both stdout and OTLP (when log provider is available).
func NewLogger(
	cfg *config.Config,
	lp *sdklog.LoggerProvider,
) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	var encoder zapcore.Encoder
	if cfg.Log.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	stdoutCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	cores := []zapcore.Core{stdoutCore}

	if lp != nil {
		otelCore := otelzap.NewCore(cfg.Service.Name, otelzap.WithLoggerProvider(global.GetLoggerProvider()))
		cores = append(cores, otelCore)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore,
		zap.AddCaller(),
		zap.AddCallerSkip(0),
	)

	return logger, nil
}
