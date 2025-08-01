# API Gateway

A lightweight API gateway built with Go and Gin that provides HTTP endpoints for user authentication and management, communicating with microservices via gRPC. This project follows Go conventions and best practices for scalable microservice architecture.

## 🚀 Features

- **HTTP API Gateway**: RESTful endpoints for user authentication
- **gRPC Client**: Communicates with microservices
- **CORS Support**: Cross-origin resource sharing enabled
- **Graceful Shutdown**: Proper server shutdown handling
- **Configuration Management**: YAML-based configuration with environment support
- **Health Check**: Built-in health check endpoint
- **Middleware Support**: Reusable middleware components
- **Docker Support**: Containerized deployment


## 📋 API Endpoints

### Authentication Endpoints

- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/login` - User login
- `POST /api/v1/users/refresh` - Refresh access token

### Health Check

- `GET /health` - Service health check

## 🏗️ Project Structure

```
apigw/
├── api/                    # API definitions and contracts
│   └── proto/             # Protocol Buffer generated files
│       ├── user-svc.pb.go
│       └── user-svc_grpc.pb.go
├── cmd/                   # Application entry points
│   └── apigw/            # Main application binary
│       └── main.go       # Application entry point
├── config/               # Configuration management
│   ├── config.go         # Configuration structs and loading
│   └── config.yaml       # Configuration file
├── internal/             # Private application code
│   ├── client/          # gRPC clients
│   │   └── grpc_client.go
│   ├── handler/         # HTTP request handlers
│   │   └── user_handler.go
│   ├── middleware/      # HTTP middleware components
│   │   ├── cors.go
│   │   ├── health.go
│   │   └── logging.go
│   └── router/          # HTTP routing
│       └── router.go

├── proto/              # Protocol Buffer definitions
│   └── user-svc.proto
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
├── Dockerfile          # Container definition
└── README.md           # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- User service running on port 9090

### Installation

1. Clone the repository and navigate to the API gateway directory:
```bash
cd apigw
```

2. Install dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

4. Run the API gateway:
```bash
make run
```

The API gateway will start on `http://localhost:8080`

## ⚙️ Configuration

The API gateway uses `config/config.yaml` for configuration:

```yaml
app:
  name: "api-gateway"
  version: "1.0.0"
  environment: "development"

server:
  http:
    port: "8080"
    host: "localhost"
    graceful_shutdown_timeout: 30s

services:
  user_service:
    host: "localhost"
    port: 9090


```

### Environment Variables

The application supports environment variable overrides with the `APIGW` prefix:

- `APIGW_SERVER_HTTP_PORT` - HTTP server port
- `APIGW_SERVER_HTTP_HOST` - HTTP server host
- `APIGW_SERVICES_USER_SERVICE_HOST` - User service host
- `APIGW_SERVICES_USER_SERVICE_PORT` - User service port

## 📖 API Usage Examples

### User Registration

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "user": {
    "id": "user-id",
    "email": "user@example.com",
    "username": "testuser",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "accessToken": "jwt-access-token",
  "refreshToken": "jwt-refresh-token"
}
```

### User Login

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/users/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
  }'
```

## 🛠️ Development

### Available Make Commands

- `make build` - Build the application
- `make run` - Run the application
- `make dev` - Run in development mode
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make fmt` - Format code
- `make lint` - Lint code
- `make proto` - Generate protobuf files
- `make build-prod` - Build for production
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container
- `make install-tools` - Install development tools

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint
```

### Protocol Buffer Generation

```bash
# Install protobuf tools
make install-tools

# Generate protobuf files
make proto
```

## 🐳 Docker Support

### Build Docker Image

```bash
make docker-build
```

### Run Docker Container

```bash
make docker-run
```

## 🔧 Error Handling

The API gateway provides consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created (registration)
- `400` - Bad Request
- `401` - Unauthorized
- `500` - Internal Server Error

## 🌐 CORS Configuration

The API gateway includes CORS middleware that allows:
- All origins (`*`)
- Common HTTP methods (GET, POST, PUT, DELETE, OPTIONS)
- Standard headers including Authorization

## 💚 Health Check

The health check endpoint returns:

```json
{
  "status": "ok",
  "service": "api-gateway"
}
```

## 📦 Dependencies

### Core Dependencies
- **Gin**: HTTP web framework
- **gRPC**: For communication with microservices
- **Viper**: Configuration management
- **Protobuf**: For gRPC message definitions

### Development Dependencies
- **testify**: Testing framework
- **golangci-lint**: Code linting

## 🏛️ Architecture

This API Gateway follows a clean architecture pattern:

1. **Entry Point** (`cmd/apigw/`): Application bootstrap
2. **Configuration** (`config/`): Environment and service configuration
3. **Handlers** (`internal/handler/`): HTTP request processing
4. **Clients** (`internal/client/`): gRPC service communication
5. **Routing** (`internal/router/`): HTTP route definitions
6. **Middleware** (`internal/middleware/`): HTTP middleware components
7. **API Definitions** (`api/proto/`): Protocol Buffer contracts

## 🔄 Go Conventions

This project follows Go community conventions:

- **Standard Layout**: Uses `cmd/`, `internal/`, `api/` directories
- **Package Naming**: Packages match directory names
- **Import Paths**: Consistent module-based imports
- **Error Handling**: Proper error propagation
- **Testing**: Comprehensive test coverage
- **Documentation**: Clear API and code documentation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following Go conventions
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Format and lint your code (`make fmt && make lint`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- Follow Go coding standards
- Write tests for new features
- Update documentation as needed
- Use meaningful commit messages
- Keep functions small and focused

## 📄 License

This project is part of the booking tickets system.

## 📞 Support

For questions and support:
- Check the API documentation above
- Review the project structure guide
- Open an issue on GitHub

---

**Built with ❤️ using Go and Gin** 