package middleware

import (
	"runtime"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRecoveryInterceptor() grpc.UnaryServerInterceptor {
	// Define recoveryFunc to handle panic
	recoveryFunc := func(p any) (err error) {
		stack := make([]byte, 64<<10)
		stack = stack[:runtime.Stack(stack, false)]
		log.Error().Msgf("panic triggered: %v, stack: %s", p, stack)
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(recoveryFunc),
	}

	return recovery.UnaryServerInterceptor(opts...)
}
