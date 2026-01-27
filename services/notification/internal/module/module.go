package module

import (
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/handler"
	"github.com/wylu1037/go-micro-boilerplate/services/notification/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"notification",
	fx.Provide(
		service.NewNotificationService,
		handler.NewNotificationGrpcHandler,
	),
)
