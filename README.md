# Go Micro Boilerplate

A microservices boilerplate for concert ticketing system built with go-micro v5 + buf + gRPC.

## Tech Stack

| Component | Choice |
|-----------|--------|
| Framework | go-micro.dev/v5 |
| RPC | gRPC + Protobuf (buf) |
| Database | PostgreSQL |
| Migration | golang-migrate |
| Cache | Redis |
| Message Queue | NATS |
| Service Discovery | mDNS (dev) / Kubernetes (prod) |

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
│   └── logger/              # Structured logging (Zap)
│
├── services/                # Microservices
│   ├── identity/            # Identity service
│   ├── catalog/             # Catalog service
│   ├── booking/             # Booking service
│   └── notification/        # Notification service
│
├── api/                     # API Gateway
│
└── migrations/              # Database migrations
    ├── 001_create_schemas.sql
    ├── identity/
    ├── catalog/
    ├── booking/
    └── notification/
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
# Run each service in separate terminals
make run-identity
make run-catalog
make run-booking
make run-notification
```

## Services

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

## Configuration

Override settings via environment variables:

```bash
export TICKETING_DATABASE_HOST=localhost
export TICKETING_DATABASE_PASSWORD=secret
export TICKETING_JWT_SECRET=your-secret-key
export MICRO_REGISTRY=kubernetes  # Production
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
