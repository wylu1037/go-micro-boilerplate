[![Powered by Go Micro](https://img.shields.io/badge/Powered%20by-Go%20Micro-00ADD8?style=flat&logo=go&logoColor=white)](https://go-micro.dev)

# Go Micro Boilerplate

A microservices boilerplate for concert ticketing system built with go-micro v4 + buf + gRPC + Etcd.

## Tech Stack

| Component | Choice |
|-----------|--------|
| Framework | go-micro.dev/v4 |
| RPC | gRPC + Protobuf (buf) |
| API Gateway | go-micro API Handler + [chi](https://github.com/go-chi/chi) |
| Middleware | Custom go-micro HandlerWrapper (Validator, Logging, Recovery) |
| Database | PostgreSQL |
| Migration | golang-migrate |
| Cache | Redis |
| Message Queue | NATS |
| Logging | [zerolog](https://github.com/rs/zerolog) |
| Service Discovery | Etcd |

## Project Structure

```
.
├── buf.yaml                 # Buf workspace configuration
├── buf.gen.yaml             # Code generation rules
├── go.work                  # Go workspace
├── Makefile                 # Build scripts
├── docker-compose.yml       # Local dev infrastructure
│
├── proto/                   # Protobuf definitions
│   ├── common/v1/           # Shared types (pagination, errors)
│   ├── identity/v1/         # Identity service API
│   ├── catalog/v1/          # Catalog service API
│   ├── booking/v1/          # Booking service API
│   └── notification/v1/     # Notification service API
│
├── gen/go/                  # Generated Go code
│   ├── common/v1/
│   ├── identity/v1/
│   ├── catalog/v1/
│   ├── booking/v1/
│   └── notification/v1/
│
├── pkg/                     # Shared libraries
│   ├── config/              # Configuration loader (Viper)
│   ├── db/                  # PostgreSQL connection pool
│   ├── cache/               # Redis client wrapper
│   ├── auth/                # JWT utilities
│   ├── middleware/          # gRPC interceptors
│   ├── errors/              # Error handling
│   └── logger/              # Structured logging (zerolog)
│
├── services/                # Microservices
│   ├── identity/            # Identity service
│   ├── catalog/             # Catalog service
│   ├── booking/             # Booking service
│   └── notification/        # Notification service
│
├── gateway/                 # API Gateway (go-micro API Handler)
│
└── migrations/              # Database migrations (golang-migrate)
    ├── 000001_create_schemas.up.sql
    ├── 000001_create_schemas.down.sql
    ├── 000002_create_identity_tables.up.sql
    ├── 000002_create_identity_tables.down.sql
    ├── 000003_create_catalog_shows.up.sql
    ├── 000004_create_booking_orders.up.sql
    └── 000005_create_notification_templates.up.sql
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [buf](https://buf.build/docs/installation)
- [golang-migrate](https://github.com/golang-migrate/migrate)

### 1. Install Tools

```bash
# Install buf
brew install bufbuild/buf/buf

# Install migrate
brew install golang-migrate
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL, Redis, NATS
docker-compose up -d

# Optional: Start debug tools (pgAdmin, Redis Commander)
docker-compose --profile debug up -d
```

### 3. Generate Proto Code

```bash
make gen
```

### 4. Run Database Migrations

```bash
make migrate-up
```

### 5. Download Dependencies

```bash
make deps
```

### 6. Run Services

```bash
# Start etcd (required for service discovery)
docker-compose up -d etcd

# Run services with etcd registry
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=localhost:2379 make run-identity
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=localhost:2379 make run-catalog
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=localhost:2379 make run-booking
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=localhost:2379 make run-notification

# Run API Gateway
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=localhost:2379 make run-gateway
```

**Environment Variables for Service Discovery:**

| Variable | Default | Description |
|----------|---------|-------------|
| `MICRO_REGISTRY` | `etcd` | Registry type (required) |
| `MICRO_REGISTRY_ADDRESS` | `localhost:2379` | etcd server endpoint |

## Services

### API Gateway

The API Gateway is the single entry point for all client requests, built with **Chi router** + **go-micro API Handler**.

#### Architecture

```
┌────────────────────────────────────────────────────────────────────┐
│                        API Gateway (:8080)                         │
├────────────────────────────────────────────────────────────────────┤
│  HTTP Request → Chi Router → Middleware Stack → go-micro Handler  │
│                                    ↓                               │
│                      Service Discovery (Etcd)                      │
│                                    ↓                               │
│                    Dynamic Service Resolution                      │
│                                    ↓                               │
│                        Backend Microservices                       │
│                    (Auto-discovered via Etcd)                      │
└────────────────────────────────────────────────────────────────────┘
```

#### Middleware Stack

**Chi Router Middleware (Gateway Level):**

| Order | Middleware | Description |
|-------|------------|-------------|
| 1 | Recovery | Panic recovery with JSON error response |
| 2 | RequestID | Request tracing (chi built-in) |
| 3 | RealIP | Extract real client IP from proxy headers |
| 4 | Logging | Structured request logging (zerolog) |
| 5 | CORS | Cross-origin resource sharing (go-chi/cors) |
| 6 | Timeout | Request timeout (60s default) |
| 7 | RateLimiter | IP-based rate limiting with LRU cache (API routes only) |

**go-micro Handler Middleware (Service Level):**

| Order | Middleware | Description |
|-------|------------|-------------|
| 1 | Recovery | Panic recovery with stack trace logging |
| 2 | Logging | Request/response logging (service, endpoint, duration) |
| 3 | Validator | Protocol buffer validation (protovalidate) |

#### Key Features

- **Dynamic Service Discovery**: Auto-discover services via Etcd registry
- **Protocol Translation**: RESTful HTTP ↔ gRPC via go-micro API handler
- **Unified Entry Point**: Single endpoint for all microservices
- **No Static Configuration**: Services register/deregister dynamically
- **Structured Logging**: Request/response logging with request ID tracing
- **Rate Limiting**: Token bucket algorithm with LRU cache for 10K IPs
- **Graceful Shutdown**: Proper cleanup on SIGINT/SIGTERM

---

### Identity Service
- User registration & login
- JWT token management
- User profile management

### Catalog Service
- Show/concert management
- Session scheduling
- Seat areas & pricing
- Inventory initialization

### Booking Service
- Order creation & management
- Inventory reservation (Redis distributed lock)
- Payment integration
- Order state machine

### Notification Service
- Event subscription (NATS)
- SMS/Email delivery
- Message template management

---

## Business Architecture

For detailed business architecture, workflows, and usage scenarios, please refer to [Business Architecture Documentation](docs/architecture.md).

## Configuration

Override settings via environment variables:

```bash
export TICKETING_DATABASE_HOST=localhost
export TICKETING_DATABASE_PASSWORD=secret
export TICKETING_JWT_SECRET=your-secret-key
export MICRO_REGISTRY=etcd
export MICRO_REGISTRY_ADDRESS=localhost:2379
```

## Make Commands

```bash
make help           # Show all commands
make gen            # Generate proto code
make build          # Build all services
make test           # Run tests
make lint           # Run linters
make docker-build   # Build Docker images
```

## Development Guide

### Adding New APIs

1. Define `.proto` file in `proto/{service}/v1/`
2. Run `make gen` to generate code
3. Implement handler in `services/{service}/internal/handler`
4. Register handler with the service

### Adding New Service

1. Create proto directory under `proto/`
2. Create service directory under `services/`
3. Update `go.work` to add new module
4. Update `Makefile` to add build targets

## License

MIT
