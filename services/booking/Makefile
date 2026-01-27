.PHONY: run run-mdns run-etcd build clean help

GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run

ETCD_ADDRESS?=localhost:2379

run-mdns:
	@echo "🚀 Starting booking service with mDNS registry..."
	$(GORUN) ./cmd/server/main.go --registry=mdns

run-etcd:
	@echo "🚀 Starting booking service with etcd registry ($(ETCD_ADDRESS))..."
	$(GORUN) ./cmd/server/main.go --registry=etcd --registry_address=$(ETCD_ADDRESS)

run: run-mdns

build:
	@echo "🔨 Building booking service..."
	$(GOBUILD) -o ../../bin/booking ./cmd/server/main.go
	@echo "✅ Build complete!"

clean:
	@echo "🧹 Cleaning..."
	rm -f ../../bin/booking

help:
	@echo "📋 Available targets:"
	@echo "  run-mdns  - 🏠 Run with mDNS registry (local development)"
	@echo "  run-etcd  - 🌐 Run with etcd registry (default: localhost:2379)"
	@echo "              Use ETCD_ADDRESS=host:port to specify custom address"
	@echo "  run       - 🏠 Alias for run-mdns"
	@echo "  build     - 🔨 Build booking binary"
	@echo "  clean     - 🧹 Remove built binary"
	@echo ""
	@echo "📝 Examples:"
	@echo "  make run-mdns"
	@echo "  make run-etcd"
	@echo "  make run-etcd ETCD_ADDRESS=192.168.1.100:2379"
