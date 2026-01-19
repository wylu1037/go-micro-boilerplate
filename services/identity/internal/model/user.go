package model

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"passwordHash"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	AvatarURL     string    `json:"avatarUrl"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
