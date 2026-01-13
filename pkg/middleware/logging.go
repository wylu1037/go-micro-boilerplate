package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// InterceptorLogger adapts zerolog logger to interceptor logger.
// This code is based on the example from https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/interceptors/logging/examples/zerolog/example_test.go
type InterceptorLogger struct {
	logger *zerolog.Logger
}

func (l *InterceptorLogger) Log(ctx context.Context, level logging.Level, msg string, fields ...any) {
	logger := l.logger.With().Fields(fields).Logger()
	switch level {
	case logging.LevelDebug:
		logger.Debug().Msg(msg)
	case logging.LevelInfo:
		logger.Info().Msg(msg)
	case logging.LevelWarn:
		logger.Warn().Msg(msg)
	case logging.LevelError:
		logger.Error().Msg(msg)
	default:
		logger.Info().Msg(msg)
	}
}

func NewLoggingInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	interceptorLogger := &InterceptorLogger{logger: logger}

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		logging.WithTimestampFormat(zerolog.TimeFormatUnix),
	}

	return logging.UnaryServerInterceptor(interceptorLogger, opts...)
}
