package module

import (
	catalogv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/catalog/v1"
	notificationv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/notification/v1"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/handler"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/repository"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/service"
	"go-micro.dev/v4"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"booking",
	fx.Provide(
		repository.NewBookingRepository,
		service.NewBookingService,
		handler.NewBookingGrpcHandler,
		// Provide clients for other services
		func(service micro.Service) catalogv1.CatalogService {
			return catalogv1.NewCatalogService("ticketing.catalog", service.Client())
		},
		func(service micro.Service) notificationv1.NotificationService {
			return notificationv1.NewNotificationService("ticketing.notification", service.Client())
		},
	),
)
