package handler

import (
	"context"

	notificationv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/notification/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/service"
)

type microNotificationGrpcHandler struct {
	svc service.NotificationService
}

func NewNotificationGrpcHandler(
	svc service.NotificationService,
) notificationv1.NotificationServiceHandler {
	return &microNotificationGrpcHandler{svc: svc}
}

func (h *microNotificationGrpcHandler) SendEmail(ctx context.Context, req *notificationv1.SendEmailRequest, resp *notificationv1.SendEmailResponse) error {
	msgID, err := h.svc.SendEmail(ctx, req.To, req.Subject, req.Body)
	if err != nil {
		resp.Success = false
		return err
	}

	resp.Success = true
	resp.MessageId = msgID
	return nil
}

func (h *microNotificationGrpcHandler) SendSMS(ctx context.Context, req *notificationv1.SendSMSRequest, resp *notificationv1.SendSMSResponse) error {
	msgID, err := h.svc.SendSMS(ctx, req.Phone, req.Message)
	if err != nil {
		resp.Success = false
		return err
	}

	resp.Success = true
	resp.MessageId = msgID
	return nil
}
