.PHONY: gen lint breaking format

gen:
	@echo "🔄 Generating protobuf code..."
	buf generate
	@echo "✅ Done!"

lint:
	@echo "🔍 Linting protobuf files..."
	buf lint

breaking:
	@echo "⚠️  Checking for breaking changes..."
	buf breaking --against '.git#branch=main'

format:
	@echo "✨ Formatting protobuf files..."
	buf format -w

help:
	@echo "📋 Available targets:"
	@echo "  gen       - 🔄 Generate protobuf code using buf"
	@echo "  lint      - 🔍 Lint protobuf files"
	@echo "  breaking  - ⚠️  Check for breaking proto changes"
	@echo "  format    - ✨ Format protobuf files"
