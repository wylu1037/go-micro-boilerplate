package router

import (
	notificationv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/notification/v1"
	"go-micro.dev/v4"
)

func NewRouter(
	service micro.Service,
	handler notificationv1.NotificationServiceHandler,
) Router {
	return &router{
		microService: service,
		handler:      handler,
	}
}

type Router interface {
	Register()
}

type router struct {
	microService micro.Service
	handler      notificationv1.NotificationServiceHandler
}

func (r *router) Register() {
	notificationv1.RegisterNotificationServiceHandler(r.microService.Server(), r.handler)
}
