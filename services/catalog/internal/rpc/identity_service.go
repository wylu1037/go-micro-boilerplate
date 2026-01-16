package rpc

import (
	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
	"go-micro.dev/v4"
)

func NewIdentityService(
	microService micro.Service,
) identityv1.IdentityService {
	return identityv1.NewIdentityService(
		"ticketing.identity",
		microService.Client(),
	)
}
