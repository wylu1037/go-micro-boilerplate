package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/model"
)

type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string) (string, error)
	SendSMS(ctx context.Context, phone, message string) (string, error)
}

type notificationService struct {
	db     *db.Pool
	logger *zap.Logger
}

func NewNotificationService(db *db.Pool, logger *zap.Logger) NotificationService {
	return &notificationService{db: db, logger: logger}
}

func (s *notificationService) SendEmail(ctx context.Context, to, subject, body string) (string, error) {
	// 1. Log to stdout (Mock Sending)
	s.logger.Info("Sending email",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("body", body),
	)

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
	s.logger.Info("Sending SMS",
		zap.String("phone", phone),
		zap.String("message", message),
	)

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
