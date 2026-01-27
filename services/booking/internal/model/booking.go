package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingStatus string

const (
	BookingStatusPendingPayment BookingStatus = "pending_payment"
	BookingStatusPaid           BookingStatus = "paid"
	BookingStatusCancelled      BookingStatus = "cancelled"
	BookingStatusRefunded       BookingStatus = "refunded"
	BookingStatusCompleted      BookingStatus = "completed"
)

type Booking struct {
	ID          string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNo     string          `gorm:"type:varchar(32);not null;uniqueIndex"`
	UserID      string          `gorm:"type:uuid;not null;index"`
	SessionID   string          `gorm:"type:uuid;not null;index"`
	SeatAreaID  string          `gorm:"type:uuid;not null"`
	Quantity    int32           `gorm:"not null"`
	UnitPrice   decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	TotalAmount decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Status      BookingStatus   `gorm:"type:varchar(20);not null;default:'pending_payment';index"`
	ExpiresAt   *time.Time      `gorm:"type:timestamptz"`
	PaidAt      *time.Time      `gorm:"type:timestamptz"`
	CancelledAt *time.Time      `gorm:"type:timestamptz"`
	CreatedAt   time.Time       `gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (Booking) TableName() string {
	return "booking.orders"
}
