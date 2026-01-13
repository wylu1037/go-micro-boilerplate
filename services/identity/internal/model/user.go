package model

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID            string
	Email         string
	PasswordHash  string
	Name          string
	Phone         string
	AvatarURL     string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
