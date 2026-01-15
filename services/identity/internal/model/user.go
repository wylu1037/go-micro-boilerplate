package model

import (
	"time"
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
