package middleware

import (
	"context"
	"runtime"

	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRecoveryMiddleware returns a go-micro server.HandlerWrapper that recovers from panics.
func NewRecoveryMiddleware(logger *zap.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := make([]byte, 64<<10)
					stack = stack[:runtime.Stack(stack, false)]
					logger.Error("panic recovered in handler",
						zap.Any("panic", r),
						zap.ByteString("stack", stack),
					)
					err = status.Errorf(codes.Internal, "internal server error: %v", r)
				}
			}()

			return fn(ctx, req, rsp)
		}
	}
}
