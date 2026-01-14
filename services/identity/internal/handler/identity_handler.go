package handler

import (
	"context"

	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/service"
)

func NewMicroIdentityHandler(
	svc service.IdentityService,
) identityv1.IdentityServiceHandler {
	return &microIdentityHandler{svc: svc}
}

type microIdentityHandler struct {
	svc service.IdentityService
}

func (h *microIdentityHandler) Register(ctx context.Context, req *identityv1.RegisterRequest, rsp *identityv1.RegisterResponse) error {
	result, err := h.svc.Register(ctx, req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		return err
	}

	rsp.UserId = result.ID
	rsp.Message = "Registration successful"
	return nil
}

func (h *microIdentityHandler) Login(ctx context.Context, req *identityv1.LoginRequest, rsp *identityv1.LoginResponse) error {
	result, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}

	rsp.AccessToken = result.AccessToken
	rsp.RefreshToken = result.RefreshToken
	rsp.ExpiresIn = result.ExpiresIn
	rsp.User = &identityv1.UserProfile{
		UserId:    result.User.ID,
		Email:     result.User.Email,
		Name:      result.User.Name,
		Phone:     result.User.Phone,
		AvatarUrl: result.User.AvatarURL,
	}
	return nil
}

func (h *microIdentityHandler) RefreshToken(ctx context.Context, req *identityv1.RefreshTokenRequest, rsp *identityv1.RefreshTokenResponse) error {
	result, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return err
	}

	rsp.AccessToken = result.AccessToken
	rsp.RefreshToken = result.RefreshToken
	rsp.ExpiresIn = result.ExpiresIn
	return nil
}

func (h *microIdentityHandler) GetProfile(ctx context.Context, req *identityv1.GetProfileRequest, rsp *identityv1.GetProfileResponse) error {
	user, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return err
	}

	rsp.User = &identityv1.UserProfile{
		UserId:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarUrl: user.AvatarURL,
	}
	return nil
}

func (h *microIdentityHandler) UpdateProfile(ctx context.Context, req *identityv1.UpdateProfileRequest, rsp *identityv1.UpdateProfileResponse) error {
	user, err := h.svc.UpdateProfile(ctx, req.UserId, req.Name, req.Phone, req.AvatarUrl)
	if err != nil {
		return err
	}

	rsp.User = &identityv1.UserProfile{
		UserId:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarUrl: user.AvatarURL,
	}
	return nil
}

func (h *microIdentityHandler) RequestPasswordReset(ctx context.Context, req *identityv1.RequestPasswordResetRequest, rsp *identityv1.RequestPasswordResetResponse) error {
	if err := h.svc.RequestPasswordReset(ctx, req.Email); err != nil {
		return err
	}

	rsp.Message = "If the email exists, a password reset link has been sent"
	return nil
}

func (h *microIdentityHandler) ResetPassword(ctx context.Context, req *identityv1.ResetPasswordRequest, rsp *identityv1.ResetPasswordResponse) error {
	if err := h.svc.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		return err
	}

	rsp.Message = "Password has been reset successfully"
	return nil
}

func (h *microIdentityHandler) ValidateToken(ctx context.Context, req *identityv1.ValidateTokenRequest, rsp *identityv1.ValidateTokenResponse) error {
	claims, err := h.svc.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		rsp.Valid = false
		return nil
	}

	rsp.Valid = true
	rsp.UserId = claims.UserID
	rsp.Email = claims.Email
	return nil
}
