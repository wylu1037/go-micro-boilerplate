package provider

import (
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/handler"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/repository"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"catalog",
	fx.Provide(
		repository.NewShowRepository,
		repository.NewVenueRepository,
		repository.NewSessionRepository,
		repository.NewSeatAreaRepository,
		service.NewCatalogService,
		handler.NewCatalogHandler,
	),
)
