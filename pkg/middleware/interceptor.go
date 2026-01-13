package middleware

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wylu1037/go-micro-boilerplate/pkg/auth"
	"github.com/wylu1037/go-micro-boilerplate/pkg/logger"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	EmailKey  ContextKey = "email"
)

// LoggingInterceptor logs gRPC requests
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		logger.Get().Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("code", code.String()),
		)

		return resp, err
	}
}

// RecoveryInterceptor recovers from panics
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Get().Error("panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// AuthInterceptor validates JWT tokens
func AuthInterceptor(jwtManager *auth.JWTManager, publicMethods []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if method is public
		for _, method := range publicMethods {
			if info.FullMethod == method {
				return handler(ctx, req)
			}
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add user info to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		return handler(ctx, req)
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetEmail extracts email from context
func GetEmail(ctx context.Context) string {
	if email, ok := ctx.Value(EmailKey).(string); ok {
		return email
	}
	return ""
}
