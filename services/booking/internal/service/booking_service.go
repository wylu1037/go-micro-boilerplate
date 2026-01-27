package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	CreateBooking(ctx context.Context, userID, sessionID, seatAreaID string, quantity int32) (*model.Booking, error)
	GetBooking(ctx context.Context, bookingID string, userID string) (*model.Booking, error)
	ListBookings(ctx context.Context, userID string, page, pageSize int, status *model.BookingStatus) ([]*model.Booking, int64, error)
	ProcessPayment(ctx context.Context, bookingID string, userID string, paymentMethod string) (string, error)
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

func (s *bookingService) CreateBooking(ctx context.Context, userID, sessionID, seatAreaID string, quantity int32) (*model.Booking, error) {
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

	// 2. Calculate Unit Price and Total Amount
	unitPrice, err := decimal.NewFromString(checkResp.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format from catalog: %w", err)
	}
	totalAmount := unitPrice.Mul(decimal.NewFromInt32(quantity))

	// 3. Generate Order Number
	orderNo, err := generateOrderNo()
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}

	// 4. Set expiration time (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	// 5. Create Booking Record (Pending Payment)
	booking := &model.Booking{
		OrderNo:     orderNo,
		UserID:      userID,
		SessionID:   sessionID,
		SeatAreaID:  seatAreaID,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalAmount: totalAmount,
		Status:      model.BookingStatusPendingPayment,
		ExpiresAt:   &expiresAt,
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking record: %w", err)
	}

	// 6. Reserve Seats
	// Note: In a distributed system, we might want to do this before creating the booking record or use a saga.
	// For simplicity, we do it here. If reservation fails, we should technically fail the booking or mark it as failed.
	reserveResp, err := s.catalogClient.ReserveSeats(ctx, &catalogv1.ReserveSeatsRequest{
		SessionId:  sessionID,
		SeatAreaId: seatAreaID,
		Quantity:   quantity,
		OrderId:    booking.ID, // Using BookingID as OrderID
	})

	if err != nil {
		// Attempt to mark booking as cancelled
		_ = s.repo.UpdateStatus(ctx, booking.ID, model.BookingStatusCancelled)
		return nil, fmt.Errorf("failed to reserve seats: %w", err)
	}

	if !reserveResp.Success {
		_ = s.repo.UpdateStatus(ctx, booking.ID, model.BookingStatusCancelled)
		return nil, errors.New(reserveResp.Message)
	}

	return booking, nil
}

func (s *bookingService) GetBooking(ctx context.Context, bookingID string, userID string) (*model.Booking, error) {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil || booking.UserID != userID {
		return nil, ErrBookingNotFound
	}
	return booking, nil
}

func (s *bookingService) ListBookings(ctx context.Context, userID string, page, pageSize int, status *model.BookingStatus) ([]*model.Booking, int64, error) {
	return s.repo.List(ctx, page, pageSize, userID, status)
}

func (s *bookingService) ProcessPayment(ctx context.Context, bookingID string, userID string, paymentMethod string) (string, error) {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return "", err
	}
	if booking == nil || booking.UserID != userID {
		return "", ErrBookingNotFound
	}

	if booking.Status != model.BookingStatusPendingPayment {
		return "", ErrInvalidBookingState
	}

	// Simulate payment processing
	time.Sleep(500 * time.Millisecond)

	// Simulate success (we could add logic to fail based on paymentMethod == "fail")
	if paymentMethod == "fail" {
		_ = s.repo.UpdateStatus(ctx, bookingID, model.BookingStatusCancelled)

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
		Body:    fmt.Sprintf("Your booking %s has been confirmed. Total paid: %s", booking.ID, booking.TotalAmount.String()),
	})

	return "txn_" + booking.ID, nil
}

// generateOrderNo generates a unique order number
func generateOrderNo() (string, error) {
	// Generate order number in format: ORD + timestamp + random hex
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("ORD%d%s", time.Now().Unix(), hex.EncodeToString(b)), nil
}
