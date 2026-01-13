package provider

import (
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/handler"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/repository"
	"github.com/wylu1037/go-micro-boilerplate/services/identity/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"identity",
	fx.Provide(
		repository.NewUserRepository,
		repository.NewTokenRepository,
		service.NewIdentityService,
		handler.NewIdentityHandler,
	),
)
