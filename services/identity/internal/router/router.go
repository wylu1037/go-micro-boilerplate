package router

import (
	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
	"go-micro.dev/v4"
)

func NewRouter(
	service micro.Service,
	handler identityv1.IdentityServiceHandler,
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
	handler      identityv1.IdentityServiceHandler
}

func (r *router) Register() {
	identityv1.RegisterIdentityServiceHandler(r.microService.Server(), r.handler)
}
