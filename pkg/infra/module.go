package infra

import (
	"go.uber.org/fx"
)

// Module provides common infrastructure dependencies (logger, database, cache, auth)
// Note: Config must be provided separately by each service using NewConfig(serviceName, schema)
var Module = fx.Options(
	fx.Provide(
		NewLogger,
		NewDatabase,
		NewMicroAuth,
		NewRedis,
		NewEtcd,
		NewDistributedLocker,
	),
)
