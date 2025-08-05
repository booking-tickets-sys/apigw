# Makefile for API Gateway

.PHONY: all build test clean run proto help

# Default target
all: build

# Build the application
build:
	@echo "Building API Gateway..."
	@mkdir -p bin
	go build -o bin/apigw ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run the application
run: build
	@echo "Running API Gateway..."
	./bin/apigw

# Run server (alias for run)
server: run

# Development: run in development mode
dev:
	@echo "Starting API Gateway in development mode..."
	@echo "Note: Make sure user-service and order-service are running"
	$(MAKE) server

# Setup proto (update submodule and generate files)
proto:
	@echo "Cleaning up existing proto files..."
	rm -rf client/proto/*.pb.go
	@echo "Updating proto submodule..."
	git submodule update --remote proto
	@echo "Generating protobuf files from proto/ to client/proto/..."
	
	protoc --proto_path=proto \
		--go_out=client/proto --go_opt=paths=source_relative \
		--go-grpc_out=client/proto --go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "Proto setup completed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Setup development environment
setup-dev: install-tools deps proto
	@echo "Development environment setup completed!"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t apigw .
	@echo "Docker image built successfully!"

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name apigw-container apigw

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  server       - Run server (alias for run)"
	@echo "  dev          - Run in development mode"
	@echo "  proto        - Update submodule and generate proto files"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  install-tools - Install development tools"
	@echo "  setup-dev    - Complete development setup"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message"