.PHONY: run run-mdns run-etcd build clean help

GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run

ETCD_ADDRESS?=localhost:2379

run-mdns:
	@echo "🚀 Starting identity service with mDNS registry..."
	$(GORUN) ./cmd/server/main.go --registry=mdns

run-etcd:
	@echo "🚀 Starting identity service with etcd registry ($(ETCD_ADDRESS))..."
	$(GORUN) ./cmd/server/main.go --registry=etcd --registry_address=$(ETCD_ADDRESS)

run: run-mdns

build:
	@echo "🔨 Building identity service..."
	$(GOBUILD) -o ../../bin/identity ./cmd/server/main.go
	@echo "✅ Build complete!"

clean:
	@echo "🧹 Cleaning..."
	rm -f ../../bin/identity

help:
	@echo "📋 Available targets:"
	@echo "  run-mdns  - 🏠 Run with mDNS registry (local development)"
	@echo "  run-etcd  - 🌐 Run with etcd registry (default: localhost:2379)"
	@echo "              Use ETCD_ADDRESS=host:port to specify custom address"
	@echo "  run       - 🏠 Alias for run-mdns"
	@echo "  build     - 🔨 Build identity binary"
	@echo "  clean     - 🧹 Remove built binary"
	@echo ""
	@echo "📝 Examples:"
	@echo "  make run-mdns"
	@echo "  make run-etcd"
	@echo "  make run-etcd ETCD_ADDRESS=192.168.1.100:2379"
