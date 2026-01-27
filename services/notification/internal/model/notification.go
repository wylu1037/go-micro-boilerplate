package model

import (
	"time"
)

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypeSMS   NotificationType = "SMS"
)

type NotificationLog struct {
	ID        string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Type      NotificationType `gorm:"type:varchar(10);not null"`
	Recipient string           `gorm:"not null"`
	Subject   string           `gorm:"type:varchar(255)"` // Only for Email
	Content   string           `gorm:"type:text;not null"`
	Info      string           `gorm:"type:text"` // JSON or extra info
	CreatedAt time.Time        `gorm:"autoCreateTime"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}
