package router

import (
	bookingv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/booking/v1"
	"go-micro.dev/v4"
)

func NewRouter(
	service micro.Service,
	handler bookingv1.BookingServiceHandler,
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
	handler      bookingv1.BookingServiceHandler
}

func (r *router) Register() {
	bookingv1.RegisterBookingServiceHandler(r.microService.Server(), r.handler)
}
