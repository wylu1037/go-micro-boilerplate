# OpenTelemetry Integration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement distributed tracing across the Go Micro Boilerplate system using OpenTelemetry, enabling full visibility into request flows from Gateway to Services and Infrastructure.

**Architecture:**

- **Pattern:** Standardized SDK integration via `pkg/telemetry` module.
- **Protocol:** OTLP (gRPC) for all telemetry data export.
- **Collector:** OpenTelemetry Collector as the central processor.
- **Backend:** Jaeger for trace storage and visualization.
- **Scope:** Gateway, Identity, Catalog, Booking, Notification services, plus DB (PostgreSQL) and Cache (Redis) instrumentation.

**Tech Stack:**

- Go 1.25+
- OpenTelemetry Go SDK (`go.opentelemetry.io/otel`)
- OTel Collector
- Jaeger
- Libraries: `otelgrpc`, `otelpgx`, `redisotel`

---

## Phase 1: Tracing Infrastructure & Core Integration

### Task 1: Verify Infrastructure Connectivity

**Goal:** Ensure the Go services can reach the already deployed OpenTelemetry Collector.

**Configuration:**

- Collector gRPC Endpoint: `localhost:4317` (or service name if in same network)
- Collector HTTP Endpoint: `localhost:4318`

**Step 1: Verify Environment Variables**
Ensure the following environment variables can be set for services:

- `OTEL_EXPORTER_OTLP_ENDPOINT`
- `OTEL_SERVICE_NAME`

**Step 2: (Optional) Test Connectivity**
If running locally, ensure `telnet localhost 4317` works.
If running in Docker, ensure services can resolve the collector hostname.

### Task 2: Create Telemetry Package

**Files:**

- Create: `pkg/telemetry/tracer.go`
- Create: `pkg/telemetry/provider.go`
- Modify: `pkg/go.mod`

**Step 1: Create Tracer Provider Initialization Code**
Implement `NewTracerProvider` function in `pkg/telemetry/provider.go`.
It should:

1. Accept service name and version.
2. Configure OTLP gRPC exporter pointing to the collector.
3. Configure Resource attributes (service.name, service.version, environment).
4. Set TraceSampler (AlwaysSample for dev, configurable later).
5. Set global `otel.SetTracerProvider` and `otel.SetTextMapPropagator`.

**Step 2: Add Shutdown Hook**
Ensure the provider returns a shutdown function to be called on service exit.

### Task 3: Integrate with Infrastructure Module

**Files:**

- Modify: `pkg/infra/module.go`
- Modify: `pkg/config/config.go` (if config needed)

**Step 1: Add Telemetry to FX Module**
Update `pkg/infra/module.go` to provide the TracerProvider.
Use `fx.Lifecycle` to manage the shutdown hook.

```go
func NewTracerProvider(lc fx.Lifecycle, cfg *config.Config) (*sdktrace.TracerProvider, error) {
    // ... impl ...
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return tp.Shutdown(ctx)
        }
    })
    return tp, nil
}
```

### Task 4: Instrument API Gateway

**Files:**

- Modify: `gateway/internal/bootstrap/http_server.go`
- Modify: `gateway/go.mod`

**Step 1: Add OTel Middleware to Chi Router**
In `gateway/internal/bootstrap/http_server.go`, add `otelchi` middleware or custom middleware that starts a span.
Start the span _before_ the request is passed to go-micro handler.

### Task 5: Instrument Microservices (Identity as Pilot)

**Files:**

- Modify: `services/identity/internal/bootstrap/micro_server.go`
- Verify: `services/identity/go.mod`

**Step 1: Validate gRPC Instrumentation**
Go Micro v4 + OpenTelemetry wrapper usually handles propagation.
Ensure `opentelemetry.NewHandlerWrapper()` is correctly configured in the service construction (it seemed to be present but maybe commented out or misused in user snippets).
Confirm `otelgrpc` or go-micro plugin transmits the TraceContext.

### Task 6: Validate Database & Redis Instrumentation

**Files:**

- Validate: `pkg/db/postgres.go`
- Validate: `pkg/infra/redis.go`

**Step 1: Review Existing Instrumentation**
The user noted `otelpgx` and `redisotel` are already there.
Ensure they are using the global TracerProvider we set up in Task 2.
Since `NewPool` and `NewRedis` are in `pkg/infra`, and we are adding `NewTracerProvider` to `pkg/infra`, we need to ensure the order of initialization is correct or that they pick up the global `otel` implementation _after_ it's initialized.
_Refinement:_ It might be safer to inject `*sdktrace.TracerProvider` into `NewPool` and `NewRedis` to enforce dependency order in FX, even if they use the global one internally (or better, configure them to use the injected instance).

### Task 7: Rollout to All Services

**Files:**

- Modify: `services/catalog/internal/bootstrap/micro_server.go`
- Modify: `services/booking/internal/bootstrap/micro_server.go`
- Modify: `services/notification/internal/bootstrap/micro_server.go`

**Step 1: Apply Identity Service Pattern**
Replicate the setup from Task 5 to all other services.

## Phase 2: NATS & Manual Instrumentation (Follow-up)

_To be detailed after Phase 1 is verifying traces._
