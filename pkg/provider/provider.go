package provider

import (
	"go.uber.org/fx"
)

// InfraModule provides common infrastructure dependencies (logger, database, cache, auth)
// Note: Config must be provided separately by each service using NewConfig(serviceName, schema)
var InfraModule = fx.Options(
	fx.Provide(
		NewLogger,
		NewDatabase,
		NewMicroAuth,
	),
)
