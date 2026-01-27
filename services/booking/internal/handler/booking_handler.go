package handler

import (
	"context"

	bookingv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/booking/v1"
	commonv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/common/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/service"
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
	// TODO: Get UserID from context (extracted from auth token middleware)
	userID := "test-user-id" // Placeholder

	booking, err := h.svc.CreateBooking(ctx, userID, req.ShowId, req.SessionId, req.SeatAreaId, req.Quantity)
	if err != nil {
		return err
	}

	resp.Booking = toProtoBooking(booking)
	return nil
}

func (h *microBookingGrpcHandler) GetBooking(ctx context.Context, req *bookingv1.GetBookingRequest, resp *bookingv1.GetBookingResponse) error {
	booking, err := h.svc.GetBooking(ctx, req.BookingId)
	if err != nil {
		return err
	}

	resp.Booking = toProtoBooking(booking)
	return nil
}

func (h *microBookingGrpcHandler) ListBookings(ctx context.Context, req *bookingv1.ListBookingsRequest, resp *bookingv1.ListBookingsResponse) error {
	userID := "test-user-id" // Placeholder

	var status *model.BookingStatus
	if req.Status != nil && *req.Status != bookingv1.BookingStatus_BOOKING_STATUS_UNSPECIFIED {
		s := model.BookingStatus(req.Status.String()[15:]) // Strip prefix "BOOKING_STATUS_"
		status = &s
	}

	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 10
	}

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
	txnID, err := h.svc.ProcessPayment(ctx, req.BookingId, req.PaymentMethod)
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
	case model.BookingStatusPending:
		status = bookingv1.BookingStatus_BOOKING_STATUS_PENDING
	case model.BookingStatusPaid:
		status = bookingv1.BookingStatus_BOOKING_STATUS_PAID
	case model.BookingStatusCancelled:
		status = bookingv1.BookingStatus_BOOKING_STATUS_CANCELLED
	case model.BookingStatusFailed:
		status = bookingv1.BookingStatus_BOOKING_STATUS_FAILED
	}

	return &bookingv1.Booking{
		BookingId:  b.ID,
		UserId:     b.UserID,
		ShowId:     b.ShowID,
		SessionId:  b.SessionID,
		SeatAreaId: b.SeatAreaID,
		Quantity:   b.Quantity,
		TotalPrice: b.TotalPrice.String(),
		Status:     status,
		CreatedAt:  timestamppb.New(b.CreatedAt),
		UpdatedAt:  timestamppb.New(b.UpdatedAt),
	}
}
