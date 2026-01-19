package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go-micro.dev/v4/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/redis/go-redis/v9"
	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
	identityerrors "github.com/wylu1037/go-micro-boilerplate/services/identity/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/repository"
)

func NewIdentityService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	microAuth auth.Auth,
	cfg *config.Config,
	logger *zerolog.Logger,
	cache *redis.Client,
) IdentityService {
	return &identityService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		auth:      microAuth,
		config:    cfg,
		logger:    logger,
		cache:     cache,
	}
}

type IdentityService interface {
	Register(ctx context.Context, email, password, name, phone string) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResult, error)
	GetProfile(ctx context.Context, userID string) (*model.User, error)
	UpdateProfile(ctx context.Context, userID, name, phone, avatarURL string) (*model.User, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ValidateToken(ctx context.Context, accessToken string) (*auth.Account, error)
}

type identityService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	auth      auth.Auth
	config    *config.Config
	logger    *zerolog.Logger
	cache     *redis.Client
}

func (svc *identityService) Register(ctx context.Context, email, password, name, phone string) (*model.User, error) {
	exists, err := svc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, identityerrors.ErrUserAlreadyExists
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

	if err := svc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	svc.logger.Info().Str("user_id", user.ID).Str("email", email).Msg("User registered")

	return user, nil
}

func (svc *identityService) Login(ctx context.Context, email, password string) (*model.LoginResult, error) {
	log := svc.logger.With().Str("email", email).Logger()

	user, err := svc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by email")
		if errors.Is(err, identityerrors.ErrUserNotFound) {
			return nil, identityerrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Error().Err(err).Msg("failed to compare password")
		return nil, identityerrors.ErrInvalidCredentials
	}

	if svc.cache != nil {
		cacheKey := fmt.Sprintf("auth:token:%s", user.ID)
		if val, err := svc.cache.Get(ctx, cacheKey).Result(); err == nil {
			var cachedRes model.LoginResult
			if err := json.Unmarshal([]byte(val), &cachedRes); err == nil {
				return &cachedRes, nil
			}
		} else {
			log.Error().Err(err).Msg("failed to get cached login result")
		}
	}

	// 生成 auth.Account
	account, err := svc.auth.Generate(user.ID, auth.WithMetadata(map[string]string{
		"email": user.Email,
		"name":  user.Name,
	}))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate auth account")
		return nil, err
	}

	// 生成 access token 和 refresh token
	tokenPair, err := svc.auth.Token(
		auth.WithCredentials(account.ID, account.Secret),
		auth.WithExpiry(svc.config.JWT.AccessTokenTTL),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token pair")
		return nil, err
	}

	// 保存 refresh token
	refreshTokenEntity := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(tokenPair.RefreshToken),
		ExpiresAt: time.Now().Add(svc.config.JWT.RefreshTokenTTL),
	}

	if err := svc.tokenRepo.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		log.Error().Err(err).Msg("failed to create refresh token")
		return nil, err
	}

	log.Info().Msg("User logged in")

	res := &model.LoginResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64(svc.config.JWT.AccessTokenTTL.Seconds()),
	}

	if svc.cache != nil {
		cacheKey := fmt.Sprintf("auth:token:%s", user.ID)
		if bytes, err := json.Marshal(res); err == nil {
			svc.cache.Set(ctx, cacheKey, bytes, svc.config.JWT.AccessTokenTTL)
		} else {
			log.Error().Err(err).Msg("failed to marshal login result")
		}
	}

	return res, nil
}

func (svc *identityService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResult, error) {
	tokenHash := model.HashToken(refreshToken)

	oldToken, err := svc.tokenRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := svc.userRepo.GetByID(ctx, oldToken.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	if err := svc.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return nil, err
	}

	// 使用 refresh token 生成新的 token pair
	newTokenPair, err := svc.auth.Token(
		auth.WithToken(refreshToken),
		auth.WithExpiry(svc.config.JWT.AccessTokenTTL),
	)
	if err != nil {
		return nil, err
	}

	// 保存新的 refresh token
	refreshTokenEntity := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(newTokenPair.RefreshToken),
		ExpiresAt: time.Now().Add(svc.config.JWT.RefreshTokenTTL),
	}

	if err := svc.tokenRepo.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &model.TokenResult{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresIn:    int64(svc.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

func (svc *identityService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := svc.userRepo.GetByID(ctx, userID)
	if err != nil {
		svc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get user profile")
		return nil, err
	}
	return user, nil
}

func (svc *identityService) UpdateProfile(ctx context.Context, userID, name, phone, avatarURL string) (*model.User, error) {
	user, err := svc.userRepo.GetByID(ctx, userID)
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

	if err := svc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *identityService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := svc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, identityerrors.ErrUserNotFound) {
			// Don't reveal if user exists
			return nil
		}
		return err
	}

	token, err := svc.generateRefreshToken()
	if err != nil {
		return err
	}

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: model.HashToken(token),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := svc.tokenRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return err
	}

	// TODO: Send email with reset link containing token
	svc.logger.Info().Str("user_id", user.ID).Msg("Password reset requested")

	return nil
}

func (svc *identityService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := model.HashToken(token)

	resetToken, err := svc.tokenRepo.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := svc.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)

	// Mark token as used
	if err := svc.tokenRepo.MarkPasswordResetTokenUsed(ctx, tokenHash); err != nil {
		return err
	}

	// Invalidate all refresh tokens
	if err := svc.tokenRepo.DeleteRefreshTokensByUserID(ctx, user.ID); err != nil {
		return err
	}

	svc.logger.Info().Str("user_id", user.ID).Msg("Password reset completed")

	return nil
}

func (svc *identityService) ValidateToken(ctx context.Context, accessToken string) (*auth.Account, error) {
	return svc.auth.Inspect(accessToken)
}

func (svc *identityService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
