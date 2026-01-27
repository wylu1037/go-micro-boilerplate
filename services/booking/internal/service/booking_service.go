package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	catalogv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/catalog/v1"
	notificationv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/notification/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/repository"
)

var (
	ErrShowNotFound        = errors.New("show not found") // Although we might just rely on catalog error
	ErrNotEnoughSeats      = errors.New("not enough available seats")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrInvalidBookingState = errors.New("invalid booking state for payment")
)

type BookingService interface {
	CreateBooking(ctx context.Context, userID, showID, sessionID, seatAreaID string, quantity int32) (*model.Booking, error)
	GetBooking(ctx context.Context, bookingID string) (*model.Booking, error)
	ListBookings(ctx context.Context, userID string, page, pageSize int, status *model.BookingStatus) ([]*model.Booking, int64, error)
	ProcessPayment(ctx context.Context, bookingID string, paymentMethod string) (string, error)
}

type bookingService struct {
	repo               repository.BookingRepository
	catalogClient      catalogv1.CatalogService
	notificationClient notificationv1.NotificationService
}

func NewBookingService(repo repository.BookingRepository, catalogClient catalogv1.CatalogService, notificationClient notificationv1.NotificationService) BookingService {
	return &bookingService{
		repo:               repo,
		catalogClient:      catalogClient,
		notificationClient: notificationClient,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID, showID, sessionID, seatAreaID string, quantity int32) (*model.Booking, error) {
	// 1. Check availability via Catalog Service
	checkResp, err := s.catalogClient.CheckAvailability(ctx, &catalogv1.CheckAvailabilityRequest{
		SessionId:  sessionID,
		SeatAreaId: seatAreaID,
		Quantity:   quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	if !checkResp.Available {
		return nil, ErrNotEnoughSeats
	}

	// 2. Calculate Total Price
	pricePerSeat, err := decimal.NewFromString(checkResp.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format from catalog: %w", err)
	}
	totalPrice := pricePerSeat.Mul(decimal.NewFromInt32(quantity))

	// 3. Create Booking Record (Pending)
	booking := &model.Booking{
		UserID:     userID,
		ShowID:     showID,
		SessionID:  sessionID,
		SeatAreaID: seatAreaID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
		Status:     model.BookingStatusPending,
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking record: %w", err)
	}

	// 4. Reserve Seats
	// Note: In a distributed system, we might want to do this before creating the booking record or use a saga.
	// For simplicity, we do it here. If reservation fails, we should technically fail the booking or mark it as failed.
	reserveResp, err := s.catalogClient.ReserveSeats(ctx, &catalogv1.ReserveSeatsRequest{
		SessionId:  sessionID,
		SeatAreaId: seatAreaID,
		Quantity:   quantity,
		OrderId:    booking.ID, // Using BookingID as OrderID
	})

	if err != nil {
		// Attempt to mark booking as failed
		_ = s.repo.UpdateStatus(ctx, booking.ID, model.BookingStatusFailed)
		return nil, fmt.Errorf("failed to reserve seats: %w", err)
	}

	if !reserveResp.Success {
		_ = s.repo.UpdateStatus(ctx, booking.ID, model.BookingStatusFailed)
		return nil, errors.New(reserveResp.Message)
	}

	return booking, nil
}

func (s *bookingService) GetBooking(ctx context.Context, bookingID string) (*model.Booking, error) {
	return s.repo.GetByID(ctx, bookingID)
}

func (s *bookingService) ListBookings(ctx context.Context, userID string, page, pageSize int, status *model.BookingStatus) ([]*model.Booking, int64, error) {
	return s.repo.List(ctx, page, pageSize, userID, status)
}

func (s *bookingService) ProcessPayment(ctx context.Context, bookingID string, paymentMethod string) (string, error) {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return "", err
	}
	if booking == nil {
		return "", ErrBookingNotFound
	}

	if booking.Status != model.BookingStatusPending {
		return "", ErrInvalidBookingState
	}

	// Simulate payment processing
	time.Sleep(500 * time.Millisecond)

	// Simulate success (we could add logic to fail based on paymentMethod == "fail")
	if paymentMethod == "fail" {
		_ = s.repo.UpdateStatus(ctx, bookingID, model.BookingStatusFailed)

		// Release seats
		_, _ = s.catalogClient.ReleaseSeats(ctx, &catalogv1.ReleaseSeatsRequest{
			SessionId:  booking.SessionID,
			SeatAreaId: booking.SeatAreaID,
			Quantity:   booking.Quantity,
			OrderId:    booking.ID,
		})

		return "", errors.New("payment failed")
	}

	// Payment Success
	if err := s.repo.UpdateStatus(ctx, bookingID, model.BookingStatusPaid); err != nil {
		return "", fmt.Errorf("failed to update booking status: %w", err)
	}

	// Send Notification
	// We do this asynchronously or synchronously. Based on plan, we just call it.
	// We don't block the response on email failure, just log it.
	_, _ = s.notificationClient.SendEmail(ctx, &notificationv1.SendEmailRequest{
		To:      "user@example.com", // In real app, fetch user email from Identity Service
		Subject: "Booking Confirmed",
		Body:    fmt.Sprintf("Your booking %s has been confirmed. Total paid: %s", booking.ID, booking.TotalPrice.String()),
	})

	return "txn_" + booking.ID, nil
}
