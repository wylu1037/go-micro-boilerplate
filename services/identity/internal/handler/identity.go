package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/service"
)

func NewIdentityHandler(
	svc service.IdentityService,
) identityv1.IdentityServiceServer {
	return &IdentityHandler{svc: svc}
}

type IdentityHandler struct {
	identityv1.UnimplementedIdentityServiceServer
	svc service.IdentityService
}

func (h *IdentityHandler) Register(ctx context.Context, req *identityv1.RegisterRequest) (*identityv1.RegisterResponse, error) {
	user, err := h.svc.Register(ctx, req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &identityv1.RegisterResponse{
		UserId:  user.ID,
		Message: "Registration successful",
	}, nil
}

func (h *IdentityHandler) Login(ctx context.Context, req *identityv1.LoginRequest) (*identityv1.LoginResponse, error) {
	result, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &identityv1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: &identityv1.UserProfile{
			UserId:    result.User.ID,
			Email:     result.User.Email,
			Name:      result.User.Name,
			Phone:     result.User.Phone,
			AvatarUrl: result.User.AvatarURL,
		},
	}, nil
}

func (h *IdentityHandler) RefreshToken(ctx context.Context, req *identityv1.RefreshTokenRequest) (*identityv1.RefreshTokenResponse, error) {
	result, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, model.ErrTokenNotFound) || errors.Is(err, model.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
		}
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &identityv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *IdentityHandler) GetProfile(ctx context.Context, req *identityv1.GetProfileRequest) (*identityv1.GetProfileResponse, error) {
	user, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get profile")
	}

	return &identityv1.GetProfileResponse{
		User: &identityv1.UserProfile{
			UserId:    user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			AvatarUrl: user.AvatarURL,
		},
	}, nil
}

func (h *IdentityHandler) UpdateProfile(ctx context.Context, req *identityv1.UpdateProfileRequest) (*identityv1.UpdateProfileResponse, error) {
	user, err := h.svc.UpdateProfile(ctx, req.UserId, req.Name, req.Phone, req.AvatarUrl)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to update profile")
	}

	return &identityv1.UpdateProfileResponse{
		User: &identityv1.UserProfile{
			UserId:    user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			AvatarUrl: user.AvatarURL,
		},
	}, nil
}

func (h *IdentityHandler) RequestPasswordReset(ctx context.Context, req *identityv1.RequestPasswordResetRequest) (*identityv1.RequestPasswordResetResponse, error) {
	if err := h.svc.RequestPasswordReset(ctx, req.Email); err != nil {
		return nil, status.Error(codes.Internal, "failed to request password reset")
	}

	return &identityv1.RequestPasswordResetResponse{
		Message: "If the email exists, a password reset link has been sent",
	}, nil
}

func (h *IdentityHandler) ResetPassword(ctx context.Context, req *identityv1.ResetPasswordRequest) (*identityv1.ResetPasswordResponse, error) {
	if err := h.svc.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		if errors.Is(err, model.ErrTokenNotFound) || errors.Is(err, model.ErrTokenExpired) || errors.Is(err, model.ErrTokenUsed) {
			return nil, status.Error(codes.InvalidArgument, "invalid or expired token")
		}
		return nil, status.Error(codes.Internal, "failed to reset password")
	}

	return &identityv1.ResetPasswordResponse{
		Message: "Password has been reset successfully",
	}, nil
}

func (h *IdentityHandler) ValidateToken(ctx context.Context, req *identityv1.ValidateTokenRequest) (*identityv1.ValidateTokenResponse, error) {
	claims, err := h.svc.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &identityv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &identityv1.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}
