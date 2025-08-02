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
- **Shared Proto Definitions**: Uses git submodules for shared protocol buffer definitions
- **Automated Build System**: Comprehensive Makefile for development workflow
- **Security**: JWT token validation and secure communication
- **Monitoring**: Built-in metrics and health monitoring
- **Scalability**: Designed for horizontal scaling

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
├── pb/                    # Generated Protocol Buffer files
│   ├── user-svc.pb.go    # Generated protobuf messages
│   └── user-svc_grpc.pb.go # Generated gRPC service definitions
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
│   │   └── health.go
│   └── router/          # HTTP routing
│       └── router.go
├── submodules/          # Git submodules for shared protos
│   └── user-svc.proto   # Shared protocol buffer definitions
├── bin/                 # Build output directory
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
├── Dockerfile          # Container definition
└── README.md           # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or later
- Protocol Buffer compiler (`protoc`)
- User service running on port 9090

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

### Configuration Sources (in order of precedence):
1. **Environment Variables** (highest priority)
2. **Config File** (if found)
3. **Default Values** (lowest priority)

### Environment Variables

The application supports environment variable overrides with the `APIGW` prefix:

- `APIGW_SERVER_HTTP_PORT` - HTTP server port
- `APIGW_SERVER_HTTP_HOST` - HTTP server host
- `APIGW_SERVICES_USER_SERVICE_HOST` - User service host
- `APIGW_SERVICES_USER_SERVICE_PORT` - User service port

### Config File Locations
The service will search for `config.yaml` in the following locations:
- Current directory (`.`)
- `./config/` directory
- `/etc/apigw/` (for system-wide config)
- `$HOME/.apigw/` (for user-specific config)

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

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/users/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
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
| `make test-coverage` | Run tests with coverage |
| `make clean` | Clean build artifacts |
| `make deps` | Install dependencies |
| `make fmt` | Format code |
| `make lint` | Lint code |
| `make proto` | Generate protobuf files |
| `make proto-clean` | Clean generated protobuf files |
| `make proto-regenerate` | Clean and regenerate protobuf files |
| `make build-prod` | Build for production |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |
| `make install-tools` | Install development tools |
| `make setup-dev` | Complete development setup |

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

The project uses Protocol Buffers for gRPC communication. The protobuf files are generated from the `submodules/` directory to the `pb/` directory.

```bash
# Install protobuf tools
make install-tools

# Generate protobuf files from submodules
make proto

# Clean generated protobuf files
make proto-clean

# Regenerate protobuf files (clean + generate)
make proto-regenerate
```

### Development Setup

For a complete development environment setup:

```bash
# This will install tools, dependencies, and generate proto files
make setup-dev
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

### Docker Compose (Optional)

If you have multiple services, you can create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  apigw:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APIGW_SERVICES_USER_SERVICE_HOST=user-service
      - APIGW_SERVICES_USER_SERVICE_PORT=9090
    depends_on:
      - user-service
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  user-service:
    image: your-user-service:latest
    ports:
      - "9090:9090"
    restart: unless-stopped
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

1. **Entry Point** (`cmd/apigw/`): Application bootstrap and configuration
2. **Configuration** (`config/`): Environment and service configuration management
3. **Handlers** (`internal/handler/`): HTTP request processing and response formatting
4. **Clients** (`internal/client/`): gRPC service communication layer
5. **Routing** (`internal/router/`): HTTP route definitions and middleware setup
6. **Middleware** (`internal/middleware/`): HTTP middleware components (CORS, health, etc.)
7. **API Definitions** (`pb/`): Generated Protocol Buffer contracts and gRPC stubs
8. **Shared Protos** (`submodules/`): Shared protocol buffer definitions across services

### Data Flow

```
HTTP Request → Router → Middleware → Handler → gRPC Client → Microservice
                ↓
HTTP Response ← Handler ← gRPC Response ← Microservice
```

## 🔄 Go Conventions

This project follows Go community conventions:

- **Standard Layout**: Uses `cmd/`, `internal/`, `pb/` directories
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
make build-prod
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
export APIGW_SERVER_HTTP_PORT=8080
export APIGW_SERVICES_USER_SERVICE_HOST=your-user-service-host
export APIGW_SERVICES_USER_SERVICE_PORT=9090
```

### Kubernetes Deployment (Optional)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apigw
spec:
  replicas: 3
  selector:
    matchLabels:
      app: apigw
  template:
    metadata:
      labels:
        app: apigw
    spec:
      containers:
      - name: apigw
        image: apigw:latest
        ports:
        - containerPort: 8080
        env:
        - name: APIGW_SERVICES_USER_SERVICE_HOST
          value: "user-service"
        - name: APIGW_SERVICES_USER_SERVICE_PORT
          value: "9090"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: apigw-service
spec:
  selector:
    app: apigw
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 🔒 Security Best Practices

### JWT Token Security
- Access tokens have short expiration times (15-30 minutes)
- Refresh tokens are stored securely and rotated
- Token validation on every protected endpoint

### Communication Security
- Use TLS/SSL for all production communications
- Implement proper CORS policies
- Validate all input data

### Environment Security
- Never commit sensitive configuration to version control
- Use environment variables for secrets
- Implement proper logging without sensitive data

## 📊 Monitoring and Observability

### Health Monitoring
- Built-in health check endpoint at `/health`
- Docker health checks configured
- Kubernetes liveness and readiness probes

### Logging
- Structured logging for better observability
- Request/response logging for debugging
- Error logging with proper context

### Metrics (Future Enhancement)
- Request rate monitoring
- Response time tracking
- Error rate monitoring
- gRPC connection status

## ⚡ Performance Considerations

### Optimization Tips
- Use connection pooling for gRPC clients
- Implement proper timeout configurations
- Use efficient JSON serialization
- Consider caching for frequently accessed data

### Scaling
- Stateless design allows horizontal scaling
- Load balancing across multiple instances
- Database connection pooling
- CDN for static assets (if applicable)

## 🐛 Troubleshooting

### Common Issues

#### 1. Protobuf Generation Errors
```bash
# Clean and regenerate protobuf files
make proto-regenerate
```

#### 2. gRPC Connection Issues
- Check if the user service is running
- Verify network connectivity
- Check firewall settings
- Validate service configuration

#### 3. Build Errors
```bash
# Clean and rebuild
make clean
make deps
make build
```

#### 4. Docker Issues
```bash
# Rebuild Docker image
make docker-build
# Check container logs
docker logs <container-id>
```

### Debug Mode
Enable debug logging by setting the environment variable:
```bash
export APIGW_APP_ENVIRONMENT=development
```

### Log Analysis
```bash
# View application logs
docker logs -f apigw-container

# Check health endpoint
curl http://localhost:8080/health
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
- Update submodules when proto definitions change
- Ensure all make commands work correctly
- Test both local and Docker deployments

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No linter errors
- [ ] Protobuf files are regenerated if needed
- [ ] Build and deployment work correctly
- [ ] Security considerations addressed
- [ ] Performance impact evaluated

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
- **Booking Service**: Ticket booking and management microservice
- **Payment Service**: Payment processing microservice

## 📈 Roadmap

### Planned Features
- [ ] Metrics collection and monitoring
- [ ] Rate limiting middleware
- [ ] Circuit breaker pattern implementation
- [ ] API versioning support
- [ ] GraphQL endpoint support
- [ ] WebSocket support for real-time features
- [ ] Advanced caching strategies
- [ ] Distributed tracing integration

### Performance Improvements
- [ ] Connection pooling optimization
- [ ] Response compression
- [ ] Request batching
- [ ] Background job processing

---

**Built with ❤️ using Go and Gin** 