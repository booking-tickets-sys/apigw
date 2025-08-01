.PHONY: build run test clean deps proto proto-clean proto-regenerate install-tools setup-dev fmt lint build-prod docker-build docker-run

# Build the API gateway
build:
	go build -o bin/apigw ./cmd/apigw

# Run the API gateway
run:
	go run ./cmd/apigw

# Run the API gateway in development mode
dev:
	go run ./cmd/apigw

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out
	rm -f api/proto/*.pb.go

# Generate protobuf files from local proto directory
# Output: api/proto/user-svc.pb.go, api/proto/user-svc_grpc.pb.go
proto:
	@echo "Generating protobuf files..."
	@mkdir -p api/proto
	protoc --proto_path=proto \
		--go_out=api/proto --go_opt=module=apigw/api/proto \
		--go-grpc_out=api/proto --go-grpc_opt=module=apigw/api/proto \
		user-svc.proto
	@echo "Protobuf files generated successfully in api/proto/"

# Clean generated protobuf files
proto-clean:
	@echo "Cleaning generated protobuf files..."
	rm -f api/proto/*.pb.go
	@echo "Protobuf files cleaned"

# Regenerate protobuf files (clean + generate)
proto-regenerate: proto-clean proto

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Development tools installed successfully"

# Setup development environment (install tools, deps, and generate proto)
setup-dev: install-tools deps proto

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/apigw ./cmd/apigw

# Docker build
docker-build:
	docker build -t apigw .

# Docker run
docker-run:
	docker run -p 8080:8080 apigw 