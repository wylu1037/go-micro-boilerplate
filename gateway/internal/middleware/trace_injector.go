package middleware

import (
	"net/http"

	"go-micro.dev/v4/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// TraceContextInjector is an HTTP middleware that injects the OpenTelemetry trace context
// from the standard context.Context into go-micro metadata.
// This bridges the gap between otelchi (which stores trace context in context.Context)
// and go-micro's opentelemetry wrapper (which reads from go-micro metadata).
func TraceContextInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get the global propagator
		propagator := otel.GetTextMapPropagator()

		// Create a carrier to hold the trace context headers
		carrier := make(propagation.MapCarrier)

		// Inject the trace context from context.Context into the carrier
		propagator.Inject(ctx, carrier)

		// Convert carrier to go-micro metadata
		md := make(metadata.Metadata)
		for k, v := range carrier {
			md.Set(k, v)
		}

		// Merge with existing metadata if any
		existingMd, ok := metadata.FromContext(ctx)
		if ok {
			for k, v := range existingMd {
				if _, exists := md[k]; !exists {
					md[k] = v
				}
			}
		}

		// Create new context with go-micro metadata
		ctx = metadata.NewContext(ctx, md)

		// Continue with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
