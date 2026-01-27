package handler

import (
	"context"

	"github.com/samber/lo"
	bookingv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/booking/v1"
	commonv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/common/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/service"
	"go-micro.dev/v4/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewBookingGrpcHandler(
	svc service.BookingService,
) bookingv1.BookingServiceHandler {
	return &microBookingGrpcHandler{svc: svc}
}

type microBookingGrpcHandler struct {
	svc service.BookingService
}

func (h *microBookingGrpcHandler) CreateBooking(ctx context.Context, req *bookingv1.CreateBookingRequest, resp *bookingv1.CreateBookingResponse) error {
	userID, ok := ctx.Value("userId").(string)
	if !ok || userID == "" {
		return errors.Unauthorized("ticketing.booking", "user unauthorized")
	}

	booking, err := h.svc.CreateBooking(ctx, userID, req.SessionId, req.SeatAreaId, req.Quantity)
	if err != nil {
		return err
	}

	resp.Booking = toProtoBooking(booking)
	return nil
}

func (h *microBookingGrpcHandler) GetBooking(ctx context.Context, req *bookingv1.GetBookingRequest, resp *bookingv1.GetBookingResponse) error {
	userID, ok := ctx.Value("userId").(string)
	if !ok || userID == "" {
		return errors.Unauthorized("ticketing.booking", "user unauthorized")
	}

	booking, err := h.svc.GetBooking(ctx, req.BookingId, userID)
	if err != nil {
		return err
	}

	resp.Booking = toProtoBooking(booking)
	return nil
}

func (h *microBookingGrpcHandler) ListBookings(ctx context.Context, req *bookingv1.ListBookingsRequest, resp *bookingv1.ListBookingsResponse) error {
	userID, ok := ctx.Value("userId").(string)
	if !ok || userID == "" {
		return errors.Unauthorized("ticketing.booking", "user unauthorized")
	}

	var status *model.BookingStatus
	if req.Status != nil && *req.Status != bookingv1.BookingStatus_BOOKING_STATUS_UNSPECIFIED {
		s := model.BookingStatus(req.Status.String()[15:]) // Strip prefix "BOOKING_STATUS_"
		status = &s
	}

	page := lo.Ternary(req.Page < 1, 1, int(req.Page))
	pageSize := lo.Ternary(req.PageSize < 1, 10, int(req.PageSize))

	bookings, total, err := h.svc.ListBookings(ctx, userID, page, pageSize, status)
	if err != nil {
		return err
	}

	pbBookings := make([]*bookingv1.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = toProtoBooking(b)
	}

	resp.Bookings = pbBookings
	resp.Pagination = &commonv1.PaginationResponse{
		TotalCount: total,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		TotalPages: int32((total + int64(pageSize) - 1) / int64(pageSize)),
	}
	return nil
}

func (h *microBookingGrpcHandler) ProcessPayment(ctx context.Context, req *bookingv1.ProcessPaymentRequest, resp *bookingv1.ProcessPaymentResponse) error {
	userID, ok := ctx.Value("userId").(string)
	if !ok || userID == "" {
		return errors.Unauthorized("ticketing.booking", "user unauthorized")
	}

	txnID, err := h.svc.ProcessPayment(ctx, req.BookingId, userID, req.PaymentMethod)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return nil // Return nil error so client receives the response with Success=false
	}

	resp.Success = true
	resp.Message = "Payment processed successfully"
	resp.TransactionId = txnID
	return nil
}

func toProtoBooking(b *model.Booking) *bookingv1.Booking {
	if b == nil {
		return nil
	}

	status := bookingv1.BookingStatus_BOOKING_STATUS_UNSPECIFIED
	switch b.Status {
	case model.BookingStatusPendingPayment:
		status = bookingv1.BookingStatus_BOOKING_STATUS_PENDING
	case model.BookingStatusPaid:
		status = bookingv1.BookingStatus_BOOKING_STATUS_PAID
	case model.BookingStatusCancelled:
		status = bookingv1.BookingStatus_BOOKING_STATUS_CANCELLED
	case model.BookingStatusRefunded:
		// Map refunded to cancelled for now, or extend proto enum
		status = bookingv1.BookingStatus_BOOKING_STATUS_CANCELLED
	case model.BookingStatusCompleted:
		// Map completed to paid for now, or extend proto enum
		status = bookingv1.BookingStatus_BOOKING_STATUS_PAID
	}

	return &bookingv1.Booking{
		BookingId:  b.ID,
		UserId:     b.UserID,
		ShowId:     "", // ShowID removed from model, set empty for backward compatibility
		SessionId:  b.SessionID,
		SeatAreaId: b.SeatAreaID,
		Quantity:   b.Quantity,
		TotalPrice: b.TotalAmount.String(),
		Status:     status,
		CreatedAt:  timestamppb.New(b.CreatedAt),
		UpdatedAt:  timestamppb.New(b.UpdatedAt),
	}
}
