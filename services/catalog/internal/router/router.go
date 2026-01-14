package router

import (
	catalogv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/catalog/v1"
	"go-micro.dev/v4"
)

func NewRouter(
	service micro.Service,
	handler catalogv1.CatalogServiceHandler,
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
	handler      catalogv1.CatalogServiceHandler
}

func (r *router) Register() {
	catalogv1.RegisterCatalogServiceHandler(r.microService.Server(), r.handler)
}
