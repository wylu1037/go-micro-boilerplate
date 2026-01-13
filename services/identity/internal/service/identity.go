package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/wylu1037/go-micro-boilerplate/pkg/auth"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/repository"
)

type IdentityService interface {
	Register(ctx context.Context, email, password, name, phone string) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResult, error)
	GetProfile(ctx context.Context, userID string) (*model.User, error)
	UpdateProfile(ctx context.Context, userID, name, phone, avatarURL string) (*model.User, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ValidateToken(ctx context.Context, accessToken string) (*auth.Claims, error)
}

type identityService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *auth.JWTManager
	config     *config.Config
	logger     *zerolog.Logger
}

func NewIdentityService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager *auth.JWTManager,
	cfg *config.Config,
	logger *zerolog.Logger,
) IdentityService {
	return &identityService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		config:     cfg,
		logger:     logger,
	}
}

func (s *identityService) Register(ctx context.Context, email, password, name, phone string) (*model.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Name:         name,
		Phone:        phone,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.logger.Info().Str("user_id", user.ID).Str("email", email).Msg("User registered")

	return user, nil
}

func (s *identityService) Login(ctx context.Context, email, password string) (*model.LoginResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenTTL),
	}

	if err := s.tokenRepo.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	s.logger.Info().Str("user_id", user.ID).Msg("User logged in")

	return &model.LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *identityService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResult, error) {
	tokenHash := model.HashToken(refreshToken)

	token, err := s.tokenRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenTTL),
	}

	if err := s.tokenRepo.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &model.TokenResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *identityService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *identityService) UpdateProfile(ctx context.Context, userID, name, phone, avatarURL string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if phone != "" {
		user.Phone = phone
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *identityService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			// Don't reveal if user exists
			return nil
		}
		return err
	}

	token, err := s.generateRefreshToken()
	if err != nil {
		return err
	}

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(token),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.tokenRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return err
	}

	// TODO: Send email with reset link containing token
	s.logger.Info().Str("user_id", user.ID).Msg("Password reset requested")

	return nil
}

func (s *identityService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := model.HashToken(token)

	resetToken, err := s.tokenRepo.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)

	// Mark token as used
	if err := s.tokenRepo.MarkPasswordResetTokenUsed(ctx, tokenHash); err != nil {
		return err
	}

	// Invalidate all refresh tokens
	if err := s.tokenRepo.DeleteRefreshTokensByUserID(ctx, user.ID); err != nil {
		return err
	}

	s.logger.Info().Str("user_id", user.ID).Msg("Password reset completed")

	return nil
}

func (s *identityService) ValidateToken(ctx context.Context, accessToken string) (*auth.Claims, error) {
	return s.jwtManager.ValidateAccessToken(accessToken)
}

func (s *identityService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
