# API Gateway

A lightweight API gateway built with Go and Gin that provides HTTP endpoints for user authentication and ticket booking, communicating with microservices via gRPC. This project follows Go conventions and best practices for scalable microservice architecture.

## 🚀 Features

- **HTTP API Gateway**: RESTful endpoints for user authentication and ticket booking
- **gRPC Client**: Communicates with microservices (User Service, Order Service)
- **JWT Authentication**: Secure token-based authentication with middleware
- **Token Bucket Rate Limiting**: Advanced Redis-based rate limiting with token bucket algorithm
- **CORS Support**: Cross-origin resource sharing enabled
- **Graceful Shutdown**: Proper server shutdown handling
- **Configuration Management**: YAML-based configuration with environment support
- **Health Check**: Built-in health check endpoint
- **Middleware Support**: Reusable middleware components
- **Docker Support**: Multi-stage containerized deployment with security best practices
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
│   │   │   ├── cors.go  # CORS middleware
│   │   │   ├── error_handler.go # Error handling middleware
│   │   │   ├── jwt.go   # JWT authentication middleware
│   │   │   └── rate_limiter.go # Token bucket rate limiting middleware
│   │   └── router/      # HTTP routing
│   │       └── router.go # Route definitions
│   └── client/          # gRPC and Redis clients
│       ├── user.go      # User service client
│       ├── order.go     # Order service client
│       └── redis.go     # Redis client wrapper
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
├── Dockerfile          # Multi-stage container definition
├── README.md           # This file
└── .gitignore          # Git ignore rules
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- Protocol Buffer compiler (`protoc`)
- Redis server (for rate limiting)
- User service running on port 50051
- Order service running on port 50052

## 🔄 CI/CD Pipeline

This project includes a comprehensive CI/CD pipeline using GitHub Actions for automated testing, building, and deployment.

### Pipeline Features

- **Automated Testing**: Unit tests with race condition detection
- **Code Quality**: Automated linting and formatting checks
- **Security Scanning**: Static security analysis with gosec
- **Docker Building**: Automated Docker image creation
- **Deployment Ready**: Staging and production deployment stages

### Local Development

Run CI checks locally before pushing:

```bash
# Run all CI checks (format, lint, test, build)
make ci

# Individual checks
make fmt      # Format code
make lint     # Lint code
make test     # Run tests with race detection
make build    # Build application
```

### GitHub Actions Workflows

#### Simple CI Pipeline (`.github/workflows/ci-simple.yml`)
- Triggers on push to `main`/`develop` and pull requests
- Runs tests with race detection
- Performs code quality checks
- Builds the application
- Uploads build artifacts

#### Comprehensive CI/CD Pipeline (`.github/workflows/ci.yml`)
- Full CI/CD pipeline with security scanning
- Docker image building
- Staging deployment (develop branch)
- Production deployment (main branch)

### Pipeline Triggers

- **Push to main**: Full CI/CD pipeline with production deployment
- **Push to develop**: CI pipeline with staging deployment
- **Pull Requests**: CI pipeline for code review

For detailed CI/CD documentation, see [CI_README.md](CI_README.md).

### Installation

1. Clone the repository and navigate to the API gateway directory:
```bash
cd apigw
```

2. Initialize and update submodules:
```bash
git submodule update --init --recursive
```

3. Start Redis server:
```bash
# Using Docker
docker run -d --name redis -p 6379:6379 redis:alpine

# Or using Homebrew (macOS)
brew install redis
brew services start redis
```

4. Install dependencies and development tools:
```bash
make setup-dev
```

5. Build the application:
```bash
make build
```

6. Run the API gateway:
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

jwt:
  secret_key: "your-secret-key-change-in-production-super-secure-32-chars-minimum-2024"

redis:
  enabled: true
  host: "localhost"
  port: 6379
  db: 0
  token_bucket:
    capacity: 100           # Maximum number of tokens in the bucket
    refill_rate: 1.67       # Tokens per second (100 tokens per minute)
    refill_interval: "1m"   # How often to refill tokens

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
- `REDIS_ENABLED` - Enable/disable Redis
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port
- `REDIS_DB` - Redis database number

## 🚦 Token Bucket Rate Limiting

The API Gateway includes an advanced Redis-based token bucket rate limiter that provides:

### Features
- **Token Bucket Algorithm**: More sophisticated than simple counters, allowing for burst traffic
- **Distributed Rate Limiting**: Works across multiple server instances
- **Race Condition Safe**: Uses Redis pipelines for atomic operations
- **Burst Support**: Allows temporary burst requests beyond normal limits
- **User-based Limiting**: Limits by user ID when authenticated, falls back to IP
- **Configurable Parameters**: Customizable capacity, refill rate, and refill interval
- **Graceful Degradation**: Continues working if Redis is unavailable

### Token Bucket Configuration
```yaml
redis:
  token_bucket:
    capacity: 100           # Maximum tokens in bucket
    refill_rate: 1.67       # Tokens per second (100 per minute)
    refill_interval: "1m"   # Refill interval
```

### Rate Limit Headers
The API returns the following headers with each request:
- `X-RateLimit-Limit` - Maximum tokens allowed
- `X-RateLimit-Remaining` - Remaining tokens in bucket
- `X-RateLimit-Reset` - Unix timestamp when next refill occurs
- `X-RateLimit-RefillRate` - Tokens refilled per second

### Rate Limit Response
When rate limit is exceeded:
```json
{
  "error": "RATE_LIMIT_ERROR",
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Rate limit exceeded. Please try again later.",
  "details": {
    "remaining_tokens": 0,
    "next_refill": "2024-01-01T00:01:00Z",
    "capacity": 100,
    "refill_rate": 1.67
  }
}
```

### Custom Rate Limits
You can create custom token bucket rate limiter for specific endpoints:

```go
// Custom rate limiter: 10 tokens capacity, 0.17 tokens per second, 1-minute refill
customLimiter := middleware.CreateCustomTokenBucketMiddleware(
    redisClient,
    10,           // capacity
    0.17,         // refill rate (tokens per second)
    time.Minute,  // refill interval
    logger,
)
```

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
| `make docker-compose-up` | Start with Docker Compose |
| `make docker-compose-down` | Stop Docker Compose services |
| `make docker-compose-logs` | Show Docker Compose logs |
| `make docker-compose-clean` | Clean up Docker resources |

### Running Tests

```bash
# Run all tests
make test

# Run rate limiter tests (requires Redis)
go test ./internal/app/middleware -v
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
- `RATE_LIMIT_ERROR` - Rate limit exceeded (429)
- `INTERNAL_ERROR` - Internal server errors (500)

## 🐳 Docker Support

### Multi-Stage Build

The Dockerfile uses a multi-stage build for optimized production images:

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
# ... build process

# Final stage
FROM alpine:latest
# ... minimal runtime image
```

### Build Docker Image

```bash
make docker-build
```

### Run Docker Container

```bash
make docker-run
```

### Docker Compose Support

The project includes Docker Compose commands for easy development:

```bash
# Start all services
make docker-compose-up

# View logs
make docker-compose-logs

# Stop services
make docker-compose-down

# Clean up resources
make docker-compose-clean
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
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - user-service
      - order-service
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
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
- **Redis**: For distributed rate limiting

### Development Dependencies
- **testify**: Testing framework
- **golangci-lint**: Code linting

## 🏛️ Architecture

This API Gateway follows a clean architecture pattern:

1. **Entry Point** (`cmd/api/`): Application bootstrap and configuration
2. **Configuration** (`internal/app/config/`): Environment and service configuration management
3. **Handlers** (`internal/app/handler/`): HTTP request processing and response formatting
4. **Clients** (`internal/client/`): gRPC service communication layer and Redis client
5. **Routing** (`internal/app/router/`): HTTP route definitions and middleware setup
6. **Middleware** (`internal/app/middleware/`): HTTP middleware components (CORS, JWT, Error handling, Token Bucket Rate limiting)
7. **DTOs** (`internal/app/domains/dto/`): Data Transfer Objects for request/response
8. **Error Handling** (`internal/app/domains/errs/`): Custom error types and gRPC to HTTP error conversion
9. **Utilities** (`pkg/utils/`): JWT token utilities and logging
10. **API Definitions** (`client/proto/`): Generated Protocol Buffer contracts and gRPC stubs
11. **Shared Protos** (`proto/`): Protocol buffer definitions

### Data Flow

```
HTTP Request → Router → Middleware (CORS, Rate Limit, JWT) → Handler → gRPC Client → Microservice
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
export REDIS_HOST=your-redis-host
export REDIS_PORT=6379
```

## 🔒 Security Best Practices

### JWT Token Security
- Access tokens have short expiration times (15 minutes)
- Refresh tokens are stored securely and rotated
- Token validation on every protected endpoint
- Secure token generation and validation

### Rate Limiting Security
- Token bucket algorithm prevents abuse and DDoS attacks
- User-based limiting for authenticated requests
- IP-based fallback for anonymous requests
- Configurable burst limits for legitimate traffic spikes

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

### Rate Limiting Monitoring
- Rate limit headers provide visibility into usage
- Structured logging for rate limit events
- Redis metrics for rate limiting performance

### Logging
- Structured logging using Logrus
- Request/response logging for debugging
- Error logging with proper context
- Rate limiting event logging

## ⚡ Performance Considerations

### Optimization Tips
- Use connection pooling for gRPC clients
- Implement proper timeout configurations
- Use efficient JSON serialization
- Consider caching for frequently accessed data
- Redis connection pooling for rate limiting

### Scaling
- Stateless design allows horizontal scaling
- Load balancing across multiple instances
- gRPC connection pooling
- Efficient error handling
- Distributed rate limiting with Redis

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

#### 4. Redis Connection Issues
- Check if Redis server is running
- Verify Redis host and port configuration
- Check Redis authentication if configured
- Validate Redis database number

#### 5. Rate Limiting Issues
- Verify Redis connection
- Check rate limit configuration
- Monitor Redis memory usage
- Validate rate limit headers
- Check token bucket parameters

#### 6. Build Errors
```bash
# Clean and rebuild
make clean
make deps
make build
```

#### 7. Docker Issues
```bash
# Rebuild Docker image
make docker-build
# Check container logs
docker logs <container-id>
```

#### 8. Token Bucket Rate Limiting Issues
- Verify Redis is running and accessible
- Check token bucket configuration parameters
- Monitor token consumption patterns
- Validate refill rate calculations

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

# Check rate limiting headers
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/tickets/event-123/purchase

# Test rate limiting with token bucket
for i in {1..110}; do curl http://localhost:8080/health; done

# Monitor Redis for rate limiting data
redis-cli keys "*rate_limit*"
```

### Rate Limiting Debugging
```bash
# Check Redis rate limiting keys
redis-cli keys "*token_bucket*"

# Monitor token bucket state
redis-cli hgetall "token_bucket:client_ip"

# Check Redis memory usage
redis-cli info memory
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
- Test rate limiting functionality with Redis
- Validate token bucket rate limiting behavior
- **Run CI checks locally**: Always run `make ci` before pushing changes
- **Check CI pipeline**: Ensure GitHub Actions pass before merging PRs

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
- [ ] Rate limiting tested with Redis
- [ ] Token bucket algorithm validated
- [ ] CI checks pass locally (`make ci`)
- [ ] GitHub Actions pipeline passes
- [ ] No security vulnerabilities detected

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
- [x] ✅ Token bucket rate limiting with Redis
- [x] ✅ Multi-stage Docker build
- [x] ✅ Docker Compose support
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