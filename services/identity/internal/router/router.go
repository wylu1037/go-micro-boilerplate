package router

import (
	identityv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/identity/v1"
	"google.golang.org/grpc"
)

func NewRouter(
	grpcServer *grpc.Server,
	identityServiceServer identityv1.IdentityServiceServer,
) Router {
	return &router{
		grpcServer:            grpcServer,
		identityServiceServer: identityServiceServer,
	}
}

type Router interface {
	Register()
}

type router struct {
	grpcServer            *grpc.Server
	identityServiceServer identityv1.IdentityServiceServer
}

func (r *router) Register() {
	identityv1.RegisterIdentityServiceServer(r.grpcServer, r.identityServiceServer)
}
