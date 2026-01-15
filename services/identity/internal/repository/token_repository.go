package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	identityerrors "github.com/wylu1037/go-micro-boilerplate/services/identity/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/model"
)

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteRefreshTokensByUserID(ctx context.Context, userID string) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	CreatePasswordResetToken(ctx context.Context, token *model.PasswordResetToken) error
	GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error
}

type tokenRepository struct {
	db *db.Pool
}

func NewTokenRepository(db *db.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	query := `
		INSERT INTO identity.refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *tokenRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM identity.refresh_tokens
		WHERE token_hash = $1
	`

	token := &model.RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, identityerrors.ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, identityerrors.ErrTokenExpired
	}

	return token, nil
}

func (r *tokenRepository) DeleteRefreshTokensByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM identity.refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM identity.refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func (r *tokenRepository) CreatePasswordResetToken(ctx context.Context, token *model.PasswordResetToken) error {
	query := `
		INSERT INTO identity.password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *tokenRepository) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used, created_at
		FROM identity.password_reset_tokens
		WHERE token_hash = $1
	`

	token := &model.PasswordResetToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Used,
		&token.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, identityerrors.ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	if token.Used {
		return nil, identityerrors.ErrTokenUsed
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, identityerrors.ErrTokenExpired
	}

	return token, nil
}

func (r *tokenRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error {
	query := `UPDATE identity.password_reset_tokens SET used = true WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}
