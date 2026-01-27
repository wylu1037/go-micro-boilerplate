package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusPaid      BookingStatus = "PAID"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusFailed    BookingStatus = "FAILED"
)

type Booking struct {
	ID         string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string          `gorm:"type:uuid;not null;index"`
	ShowID     string          `gorm:"type:uuid;not null"`
	SessionID  string          `gorm:"type:uuid;not null;index"`
	SeatAreaID string          `gorm:"type:uuid;not null"`
	Quantity   int32           `gorm:"not null"`
	TotalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Status     BookingStatus   `gorm:"type:varchar(20);not null;default:'PENDING';index"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
}

func (Booking) TableName() string {
	return "bookings"
}
