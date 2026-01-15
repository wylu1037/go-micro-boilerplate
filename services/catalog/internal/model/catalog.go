package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Show struct {
	ID          string
	Title       string
	Description string
	Artist      string
	Category    string
	PosterURL   string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Venue struct {
	ID        string
	Name      string
	City      string
	Address   string
	Capacity  int32
	CreatedAt time.Time
}

type Session struct {
	ID            string
	ShowID        string
	VenueID       string
	Venue         *Venue
	StartTime     time.Time
	EndTime       *time.Time
	SaleStartTime *time.Time
	SaleEndTime   *time.Time
	Status        string
	CreatedAt     time.Time
}

type SeatArea struct {
	ID             string
	SessionID      string
	Name           string
	Price          decimal.Decimal
	TotalSeats     int32
	AvailableSeats int32
	CreatedAt      time.Time
}
