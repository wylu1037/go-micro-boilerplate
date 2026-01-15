package middleware

import (
	"context"
	"slices"
	"strings"

	"go-micro.dev/v4/auth"
	"go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

func AuthWrapper(a auth.Auth, publicEndpoints []string) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// 1. Check if the endpoint is in the whitelist
			// req.Method() returns "Service.Method", e.g., "IdentityService.Login"
			if slices.Contains(publicEndpoints, req.Method()) {
				// Try to parse the token but don't require it (to extract user information for Public interfaces)
				// Here we can choose to do nothing and just let it pass
				return fn(ctx, req, rsp)
			}

			// 2. Get token from metadata
			// go-micro will automatically put the Authorization header into metadata
			md, ok := metadata.FromContext(ctx)
			if !ok {
				return errors.Unauthorized(req.Service(), "no metadata found")
			}

			token, ok := md["Authorization"]
			if !ok || token == "" {
				return errors.Unauthorized(req.Service(), "no auth token provided")
			}
			token = strings.TrimPrefix(token, "Bearer ")

			// 3. Verify token
			account, err := a.Inspect(token)
			if err != nil {
				return errors.Unauthorized(req.Service(), "invalid token: %v", err)
			}

			// 4. Inject Account information into Context (optional)
			// Note: go-micro's auth plugin usually automatically handles context injection,
			// but we may need to handle it ourselves when manually inspecting
			// Here we mainly rely on the result of auth.Inspect to decide whether to let it pass

			// Can be extended here, for example, put UserID into connection context
			ctx = context.WithValue(ctx, "userId", account.ID)

			return fn(ctx, req, rsp)
		}
	}
}
