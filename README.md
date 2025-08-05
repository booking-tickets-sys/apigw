# API Gateway

A lightweight API gateway built with Go and Gin that provides HTTP endpoints for user authentication and ticket booking, communicating with microservices via gRPC. This project follows Go conventions and best practices for scalable microservice architecture.

## 🚀 Features

- **HTTP API Gateway**: RESTful endpoints for user authentication and ticket booking
- **gRPC Client**: Communicates with microservices (User Service, Order Service)
- **JWT Authentication**: Secure token-based authentication with middleware
- **CORS Support**: Cross-origin resource sharing enabled
- **Graceful Shutdown**: Proper server shutdown handling
- **Configuration Management**: YAML-based configuration with environment support
- **Health Check**: Built-in health check endpoint
- **Middleware Support**: Reusable middleware components
- **Docker Support**: Containerized deployment
- **Clean Architecture**: Well-organized code structure following Go conventions

## 📋 API Endpoints

### Authentication Endpoints

- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/login` - User login
- `POST /api/v1/users/refresh` - Refresh access token

### Ticket Management Endpoints

- `POST /api/v1/tickets/:event_id/purchase` - Purchase ticket (requires authentication)

### Health Check

- `GET /health` - Service health check

## 🏗️ Project Structure

```
apigw/
├── client/                 # Generated Protocol Buffer files
│   └── proto/             # Generated protobuf files
│       ├── user-svc.pb.go    # Generated protobuf messages
│       ├── user-svc_grpc.pb.go # Generated gRPC service definitions
│       ├── order-svc.pb.go  # Generated protobuf messages
│       └── order-svc_grpc.pb.go # Generated gRPC service definitions
├── cmd/                   # Application entry points
│   └── api/              # Main application binary
│       └── main.go       # Application entry point
├── internal/             # Private application code
│   ├── app/             # Application layer
│   │   ├── config/      # Configuration management
│   │   │   └── config.go # Configuration structs and loading
│   │   ├── domains/     # Domain layer
│   │   │   ├── dto/     # Data Transfer Objects
│   │   │   │   └── user.go # User DTOs
│   │   │   └── errs/    # Error definitions
│   │   │       └── errors.go # Custom error types
│   │   ├── handler/     # HTTP request handlers
│   │   │   ├── user.go  # User handler
│   │   │   └── order.go # Order handler
│   │   ├── middleware/  # HTTP middleware components
│   │   │   ├── error_handler.go # Error handling middleware
│   │   │   └── jwt.go   # JWT authentication middleware
│   │   └── router/      # HTTP routing
│   │       └── router.go # Route definitions
│   └── client/          # gRPC clients
│       ├── user.go      # User service client
│       └── order.go     # Order service client
├── pkg/                 # Public packages
│   └── utils/           # Utility functions
│       ├── crypt/       # Cryptographic utilities
│       │   └── token/   # JWT token utilities
│       │       ├── jwt_maker.go # JWT token maker
│       │       └── maker.go     # Token payload and validation
│       └── log/         # Logging utilities
│           └── logger.go # Logger configuration
├── proto/               # Protocol buffer definitions
│   ├── user-svc.proto   # User service protobuf definitions
│   ├── order-svc.proto  # Order service protobuf definitions
│   ├── Makefile         # Protobuf generation makefile
│   └── README.md        # Protobuf documentation
├── bin/                 # Build output directory
├── config.yaml          # Configuration file
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
├── Dockerfile          # Container definition
├── README.md           # This file
└── .gitignore          # Git ignore rules
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- Protocol Buffer compiler (`protoc`)
- User service running on port 50051
- Order service running on port 50052

### Installation

1. Clone the repository and navigate to the API gateway directory:
```bash
cd apigw
```

2. Initialize and update submodules:
```bash
git submodule update --init --recursive
```

3. Install dependencies and development tools:
```bash
make setup-dev
```

4. Build the application:
```bash
make build
```

5. Run the API gateway:
```bash
make run
```

The API gateway will start on `http://localhost:8080`

### Alternative: Quick Development Setup

For a complete development environment in one command:

```bash
make setup-dev && make build && make run
```

## ⚙️ Configuration

The API gateway uses `config.yaml` for configuration:

```yaml
app:
  name: "booking-tickets-api-gateway"
  version: "1.0.0"
  environment: "development"

server:
  http:
    host: "0.0.0.0"
    port: 8080
    read_timeout: "30s"
    write_timeout: "30s"
    idle_timeout: "60s"
    graceful_shutdown_timeout: "30s"

services:
  user_service:
    name: "user-service"
    host: "localhost"
    port: 50051
    grpc:
      keepalive_time: "30s"
      keepalive_timeout: "5s"
      keepalive_permit_without_stream: true
  order_service:
    name: "order-service"
    host: "localhost"
    port: 50052
    grpc:
      keepalive_time: "30s"
      keepalive_timeout: "5s"
      keepalive_permit_without_stream: true

jwt:
  secret_key: "your-secret-key-change-in-production"
  access_token_duration: "15m"
  refresh_token_duration: "7d"
  issuer: "booking-tickets-api-gateway"
```

### Configuration Sources (in order of precedence):
1. **Environment Variables** (highest priority)
2. **Config File** (if found)
3. **Default Values** (lowest priority)

### Environment Variables

The application supports environment variable overrides with the following pattern:
- `APP_NAME` - Application name
- `APP_VERSION` - Application version
- `APP_ENVIRONMENT` - Application environment
- `SERVER_HTTP_HOST` - HTTP server host
- `SERVER_HTTP_PORT` - HTTP server port
- `SERVICES_USER_SERVICE_HOST` - User service host
- `SERVICES_USER_SERVICE_PORT` - User service port
- `SERVICES_ORDER_SERVICE_HOST` - Order service host
- `SERVICES_ORDER_SERVICE_PORT` - Order service port
- `JWT_SECRET_KEY` - JWT secret key

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

### Purchase Ticket (Requires Authentication)

```bash
curl -X POST http://localhost:8080/api/v1/tickets/event-123/purchase \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "ticket": {
    "id": "ticket-id",
    "eventId": "event-123",
    "userId": "user-id",
    "status": "confirmed",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/users/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "your-refresh-token-here"
  }'
```

**Response:**
```json
{
  "accessToken": "new-jwt-access-token"
}
```

## 🛠️ Development

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make run` | Run the application |
| `make dev` | Run in development mode |
| `make test` | Run tests |
| `make clean` | Clean build artifacts |
| `make deps` | Install dependencies |
| `make fmt` | Format code |
| `make lint` | Lint code |
| `make proto` | Generate protobuf files |
| `make install-tools` | Install development tools |
| `make setup-dev` | Complete development setup |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

### Running Tests

```bash
# Run all tests
make test
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint
```

### Protocol Buffer Generation

The project uses Protocol Buffers for gRPC communication. The protobuf files are generated from the `proto/` directory to the `client/proto/` directory.

```bash
# Install protobuf tools
make install-tools

# Generate protobuf files
make proto
```

### Development Setup

For a complete development environment setup:

```bash
# This will install tools, dependencies, and generate proto files
make setup-dev
```

## 🔧 Error Handling

The API Gateway provides comprehensive error handling with proper HTTP status codes:

### Error Response Format
```json
{
  "error": "ERROR_TYPE",
  "code": "SPECIFIC_ERROR_CODE",
  "message": "Human-readable error message"
}
```

### Error Types
- `VALIDATION_ERROR` - Input validation errors (400)
- `AUTHENTICATION_ERROR` - Authentication failures (401)
- `AUTHORIZATION_ERROR` - Authorization failures (403)
- `NOT_FOUND_ERROR` - Resource not found (404)
- `INTERNAL_ERROR` - Internal server errors (500)

## 🐳 Docker Support

### Build Docker Image

```bash
make docker-build
```

### Run Docker Container

```bash
make docker-run
```

### Docker Compose Example

```yaml
version: '3.8'
services:
  apigw:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVICES_USER_SERVICE_HOST=user-service
      - SERVICES_USER_SERVICE_PORT=50051
      - SERVICES_ORDER_SERVICE_HOST=order-service
      - SERVICES_ORDER_SERVICE_PORT=50052
      - JWT_SECRET_KEY=your-secure-jwt-secret
    depends_on:
      - user-service
      - order-service
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  user-service:
    image: your-user-service:latest
    ports:
      - "50051:50051"
    restart: unless-stopped
  
  order-service:
    image: your-order-service:latest
    ports:
      - "50052:50052"
    restart: unless-stopped
```

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
  "service": "booking-tickets-api-gateway",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 📦 Dependencies

### Core Dependencies
- **Gin**: HTTP web framework
- **gRPC**: For communication with microservices
- **Viper**: Configuration management
- **Protobuf**: For gRPC message definitions
- **JWT**: For token-based authentication
- **Logrus**: Structured logging

### Development Dependencies
- **testify**: Testing framework
- **golangci-lint**: Code linting

## 🏛️ Architecture

This API Gateway follows a clean architecture pattern:

1. **Entry Point** (`cmd/api/`): Application bootstrap and configuration
2. **Configuration** (`internal/app/config/`): Environment and service configuration management
3. **Handlers** (`internal/app/handler/`): HTTP request processing and response formatting
4. **Clients** (`internal/client/`): gRPC service communication layer
5. **Routing** (`internal/app/router/`): HTTP route definitions and middleware setup
6. **Middleware** (`internal/app/middleware/`): HTTP middleware components (CORS, JWT, Error handling)
7. **DTOs** (`internal/app/domains/dto/`): Data Transfer Objects for request/response
8. **Error Handling** (`internal/app/domains/errs/`): Custom error types and gRPC to HTTP error conversion
9. **Utilities** (`pkg/utils/`): JWT token utilities and logging
10. **API Definitions** (`client/proto/`): Generated Protocol Buffer contracts and gRPC stubs
11. **Shared Protos** (`proto/`): Protocol buffer definitions

### Data Flow

```
HTTP Request → Router → Middleware (CORS, JWT) → Handler → gRPC Client → Microservice
                ↓
HTTP Response ← Handler ← gRPC Response ← Microservice
```

## 🔄 Go Conventions

This project follows Go community conventions:

- **Standard Layout**: Uses `cmd/`, `internal/`, `pkg/` directories
- **Package Naming**: Packages match directory names
- **Import Paths**: Consistent module-based imports
- **Error Handling**: Proper error propagation and logging
- **Testing**: Comprehensive test coverage
- **Documentation**: Clear API and code documentation
- **Build Automation**: Comprehensive Makefile for all development tasks

## 🚀 Deployment

### Production Build

```bash
# Build optimized binary for production
make build
```

### Docker Deployment

```bash
# Build and run with Docker
make docker-build
make docker-run
```

### Environment Configuration

For production deployment, ensure proper environment variables are set:

```bash
export APP_ENVIRONMENT=production
export SERVER_HTTP_PORT=8080
export SERVICES_USER_SERVICE_HOST=your-user-service-host
export SERVICES_USER_SERVICE_PORT=50051
export SERVICES_ORDER_SERVICE_HOST=your-order-service-host
export SERVICES_ORDER_SERVICE_PORT=50052
export JWT_SECRET_KEY=your-secure-jwt-secret
```

## 🔒 Security Best Practices

### JWT Token Security
- Access tokens have short expiration times (15 minutes)
- Refresh tokens are stored securely and rotated
- Token validation on every protected endpoint
- Secure token generation and validation

### Communication Security
- Use TLS/SSL for all production communications
- Implement proper CORS policies
- Validate all input data
- Secure gRPC communication

### Environment Security
- Never commit sensitive configuration to version control
- Use environment variables for secrets
- Implement proper logging without sensitive data

## 📊 Monitoring and Observability

### Health Monitoring
- Built-in health check endpoint at `/health`
- Docker health checks configured
- Structured logging for better observability

### Logging
- Structured logging using Logrus
- Request/response logging for debugging
- Error logging with proper context

## ⚡ Performance Considerations

### Optimization Tips
- Use connection pooling for gRPC clients
- Implement proper timeout configurations
- Use efficient JSON serialization
- Consider caching for frequently accessed data

### Scaling
- Stateless design allows horizontal scaling
- Load balancing across multiple instances
- gRPC connection pooling
- Efficient error handling

## 🐛 Troubleshooting

### Common Issues

#### 1. Protobuf Generation Errors
```bash
# Clean and regenerate protobuf files
make proto
```

#### 2. gRPC Connection Issues
- Check if the services are running
- Verify network connectivity
- Check firewall settings
- Validate service configuration

#### 3. JWT Authentication Issues
- Verify JWT secret configuration
- Check token expiration
- Validate token format
- Check authorization headers

#### 4. Build Errors
```bash
# Clean and rebuild
make clean
make deps
make build
```

#### 5. Docker Issues
```bash
# Rebuild Docker image
make docker-build
# Check container logs
docker logs <container-id>
```

### Debug Mode
Enable debug logging by setting the environment variable:
```bash
export APP_ENVIRONMENT=development
```

### Log Analysis
```bash
# View application logs
docker logs -f apigw-container

# Check health endpoint
curl http://localhost:8080/health

# Check authentication
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/tickets/event-123/purchase
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Initialize submodules: `git submodule update --init --recursive`
4. Make your changes following Go conventions
5. Add tests for new functionality
6. Ensure all tests pass (`make test`)
7. Format and lint your code (`make fmt && make lint`)
8. Commit your changes (`git commit -m 'Add amazing feature'`)
9. Push to the branch (`git push origin feature/amazing-feature`)
10. Open a Pull Request

### Development Guidelines

- Follow Go coding standards and conventions
- Write tests for new features and maintain high test coverage
- Update documentation as needed
- Use meaningful commit messages following conventional commits
- Keep functions small and focused on single responsibilities
- Update protobuf files when service definitions change
- Ensure all make commands work correctly
- Test both local and Docker deployments
- Verify error handling and mapping

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No linter errors
- [ ] Protobuf files are regenerated if needed
- [ ] Build and deployment work correctly
- [ ] Security considerations addressed
- [ ] Performance impact evaluated
- [ ] Error handling properly implemented

## 📄 License

This project is part of the booking tickets system.

## 📞 Support

For questions and support:
- Check the API documentation above
- Review the project structure guide
- Open an issue on GitHub
- Check the Makefile for available commands
- Review the troubleshooting section

## 🔗 Related Projects

This API Gateway is part of a larger microservices architecture:
- **User Service**: Authentication and user management microservice
- **Order Service**: Ticket booking and order management microservice

## 📈 Roadmap

### Planned Features
- [x] ✅ JWT authentication middleware
- [x] ✅ gRPC client communication
- [x] ✅ Configuration management
- [x] ✅ Health check endpoint
- [x] ✅ CORS support
- [x] ✅ Graceful shutdown
- [x] ✅ Error handling and gRPC to HTTP error conversion
- [ ] Rate limiting middleware
- [ ] Metrics collection and monitoring
- [ ] Circuit breaker pattern implementation
- [ ] API versioning support
- [ ] Advanced caching strategies
- [ ] Distributed tracing integration

### Performance Improvements
- [ ] Connection pooling optimization
- [ ] Response compression
- [ ] Request batching
- [ ] Background job processing

---

**Built with ❤️ using Go and Gin** 