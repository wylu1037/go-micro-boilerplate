.PHONY: all gen build test clean run-% migrate-up migrate-down lint docker-build docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Services
SERVICES=gateway identity catalog booking notification

# Default target
all: gen build

# ==================== Proto Generation ====================

gen:
	@echo "Generating protobuf code..."
	buf generate
	@echo "Done!"

lint-proto:
	@echo "Linting protobuf files..."
	buf lint

breaking:
	@echo "Checking for breaking changes..."
	buf breaking --against '.git#branch=main'

# ==================== Build ====================

build: $(addprefix build-,$(SERVICES))

build-%:
	@echo "Building $*..."
	@if [ "$*" = "gateway" ]; then \
		cd gateway && $(GOBUILD) -o ../bin/$* ./cmd/$*; \
	else \
		cd services/$* && $(GOBUILD) -o ../../bin/$* ./...; \
	fi

# ==================== Test ====================

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

# ==================== Run Services ====================

run-%:
	@echo "Running $*..."
	@if [ "$*" = "gateway" ]; then \
		cd gateway && $(GOCMD) run ./cmd/$*; \
	else \
		cd services/$* && $(GOCMD) run .; \
	fi

# ==================== Database Migrations ====================

MIGRATE=migrate
DB_URL=postgres://postgres:postgres@localhost:5432/ticketing?sslmode=disable

migrate-up:
	@echo "Running migrations..."
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	$(MIGRATE) -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@echo "Creating new migration: $(NAME)"
	$(MIGRATE) create -ext sql -dir migrations -seq $(NAME)

# ==================== Dependencies ====================

deps:
	@echo "Downloading dependencies..."
	cd pkg && $(GOMOD) tidy
	cd gen && $(GOMOD) tidy
	cd gateway && $(GOMOD) tidy
	@for service in identity catalog booking notification; do \
		echo "Tidying $$service..."; \
		cd services/$$service && $(GOMOD) tidy; \
		cd ../..; \
	done

# ==================== Docker ====================

docker-build: $(addprefix docker-build-,$(SERVICES))

docker-build-%:
	@echo "Building Docker image for $*..."
	docker build -t ticketing-$*:latest -f deploy/docker/Dockerfile.$* .

docker-up:
	@echo "Starting Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

docker-logs:
	docker-compose logs -f

# ==================== Linting ====================

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

# ==================== Clean ====================

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf gen/go/

# ==================== Help ====================

help:
	@echo "Available targets:"
	@echo "  gen              - Generate protobuf code using buf"
	@echo "  lint-proto       - Lint protobuf files"
	@echo "  breaking         - Check for breaking proto changes"
	@echo "  build            - Build all services"
	@echo "  build-<service>  - Build specific service (gateway, identity, catalog, booking, notification)"
	@echo "  test             - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-integration - Run integration tests"
	@echo "  run-<service>    - Run specific service"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback last migration"
	@echo "  migrate-create   - Create new migration (use NAME=<name>)"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  docker-build     - Build Docker images"
	@echo "  docker-up        - Start Docker Compose stack"
	@echo "  docker-down      - Stop Docker Compose stack"
	@echo "  lint             - Run golangci-lint"
	@echo "  clean            - Clean build artifacts"
