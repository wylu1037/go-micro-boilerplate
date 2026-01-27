package service

import (
	"context"
	"fmt"
	"log"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/model"
)

type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string) (string, error)
	SendSMS(ctx context.Context, phone, message string) (string, error)
}

type notificationService struct {
	db *db.Pool
}

func NewNotificationService(db *db.Pool) NotificationService {
	return &notificationService{db: db}
}

func (s *notificationService) SendEmail(ctx context.Context, to, subject, body string) (string, error) {
	// 1. Log to stdout (Mock Sending)
	log.Printf("📧 [EMAIL] To: %s | Subject: %s | Body: %s", to, subject, body)

	// 2. Save log to DB
	query := `
		INSERT INTO notification.notification_logs (type, recipient, subject, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id string
	err := s.db.QueryRow(ctx, query,
		model.NotificationTypeEmail,
		to,
		subject,
		body,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to save notification log: %w", err)
	}

	return id, nil
}

func (s *notificationService) SendSMS(ctx context.Context, phone, message string) (string, error) {
	// 1. Log to stdout (Mock Sending)
	log.Printf("📱 [SMS] To: %s | Message: %s", phone, message)

	// 2. Save log to DB
	query := `
		INSERT INTO notification.notification_logs (type, recipient, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id string
	err := s.db.QueryRow(ctx, query,
		model.NotificationTypeSMS,
		phone,
		message,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to save notification log: %w", err)
	}

	return id, nil
}
