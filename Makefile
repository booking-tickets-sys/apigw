# Makefile for API Gateway

.PHONY: all build test clean run proto help docker-compose

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
	go test -v -race ./...

# Run CI checks
ci: fmt lint test build
	@echo "CI checks completed successfully!"

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

# Docker Compose commands
docker-compose-up:
	@echo "Starting API Gateway and Redis with Docker Compose..."
	docker compose up -d

docker-compose-down:
	@echo "Stopping all services..."
	docker compose down

docker-compose-logs:
	@echo "Showing logs for all services..."
	docker compose logs -f

docker-compose-monitoring:
	@echo "Starting services with monitoring..."
	docker compose --profile monitoring up -d

docker-compose-clean:
	@echo "Cleaning up Docker Compose resources..."
	docker compose down -v --remove-orphans
	docker system prune -f

# Run API Gateway only (for development)
docker-apigw-only:
	@echo "Building and running API Gateway only..."
	docker build -t apigw .
	docker run -p 8080:8080 --name apigw-dev apigw

# Show help
help:
	@echo "Available targets:"
	@echo "  all                    - Build the application (default)"
	@echo "  build                  - Build the application"
	@echo "  test                   - Run tests"
	@echo "  ci                     - Run all CI checks (fmt, lint, test, build)"
	@echo "  clean                  - Clean build artifacts"
	@echo "  run                    - Build and run the application"
	@echo "  server                 - Run server (alias for run)"
	@echo "  dev                    - Run in development mode"
	@echo "  proto                  - Update submodule and generate proto files"
	@echo "  deps                   - Install dependencies"
	@echo "  fmt                    - Format code"
	@echo "  lint                   - Lint code"
	@echo "  install-tools          - Install development tools"
	@echo "  setup-dev              - Complete development setup"
	@echo "  docker-build           - Build Docker image"
	@echo "  docker-run             - Run Docker container"
	@echo ""
	@echo "Docker Compose commands:"
	@echo "  docker-compose-up      - Start API Gateway and Redis"
	@echo "  docker-compose-down    - Stop all services"
	@echo "  docker-compose-logs    - Show logs for all services"
	@echo "  docker-compose-monitoring - Start with monitoring (Redis Commander)"
	@echo "  docker-compose-clean   - Clean up all Docker resources"
	@echo "  docker-apigw-only      - Run API Gateway only (for development)"
	@echo ""
	@echo "Environment variables:"
	@echo "  APP_ENVIRONMENT        - Set to 'production' for production mode"
	@echo "  JWT_SECRET_KEY         - JWT secret key"
	@echo "  DB_PASSWORD            - Database password"
	@echo "  DB_NAME                - Database name"
	@echo "  LOG_LEVEL              - Log level (debug, info, warn, error)"
	@echo ""
	@echo "  help                   - Show this help message"