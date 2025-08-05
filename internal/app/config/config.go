package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// AppConfig represents application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

// HTTPConfig represents HTTP server configuration
type HTTPConfig struct {
	Host                    string        `mapstructure:"host"`
	Port                    int           `mapstructure:"port"`
	ReadTimeout             time.Duration `mapstructure:"read_timeout"`
	WriteTimeout            time.Duration `mapstructure:"write_timeout"`
	IdleTimeout             time.Duration `mapstructure:"idle_timeout"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}

// ServicesConfig represents microservices configuration
type ServicesConfig struct {
	UserService  ServiceConfig `mapstructure:"user_service"`
	OrderService ServiceConfig `mapstructure:"order_service"`
}

// UserServiceConfig is an alias for ServiceConfig for user service
type UserServiceConfig = ServiceConfig

// OrderServiceConfig is an alias for ServiceConfig for order service
type OrderServiceConfig = ServiceConfig

// ServiceConfig represents individual service configuration
type ServiceConfig struct {
	Name string     `mapstructure:"name"`
	Host string     `mapstructure:"host"`
	Port int        `mapstructure:"port"`
	GRPC GRPCConfig `mapstructure:"grpc"`
}

// GRPCConfig represents gRPC client configuration
type GRPCConfig struct {
	KeepaliveTime                time.Duration `mapstructure:"keepalive_time"`
	KeepaliveTimeout             time.Duration `mapstructure:"keepalive_timeout"`
	KeepalivePermitWithoutStream bool          `mapstructure:"keepalive_permit_without_stream"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	Issuer               string        `mapstructure:"issuer"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Enable environment variable support
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "booking-tickets-api-gateway")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")

	// Server defaults
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.http.read_timeout", "30s")
	v.SetDefault("server.http.write_timeout", "30s")
	v.SetDefault("server.http.idle_timeout", "60s")
	v.SetDefault("server.http.graceful_shutdown_timeout", "30s")

	// JWT defaults
	v.SetDefault("jwt.secret_key", "booking-tickets-api-gateway-secret-key-2024-development")
	v.SetDefault("jwt.access_token_duration", "15m")
	v.SetDefault("jwt.refresh_token_duration", "168h")
	v.SetDefault("jwt.issuer", "booking-tickets-api-gateway")

	// Service defaults
	v.SetDefault("services.user_service.name", "user-service")
	v.SetDefault("services.user_service.host", "localhost")
	v.SetDefault("services.user_service.port", 50051)
	v.SetDefault("services.user_service.grpc.keepalive_time", "30s")
	v.SetDefault("services.user_service.grpc.keepalive_timeout", "5s")
	v.SetDefault("services.user_service.grpc.keepalive_permit_without_stream", true)

	v.SetDefault("services.order_service.name", "order-service")
	v.SetDefault("services.order_service.host", "localhost")
	v.SetDefault("services.order_service.port", 50052)
	v.SetDefault("services.order_service.grpc.keepalive_time", "30s")
	v.SetDefault("services.order_service.grpc.keepalive_timeout", "5s")
	v.SetDefault("services.order_service.grpc.keepalive_permit_without_stream", true)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if c.Server.HTTP.Port <= 0 || c.Server.HTTP.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.HTTP.Port)
	}

	if c.Server.HTTP.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.Server.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT secret key must be set")
	}

	if c.Services.UserService.Host == "" {
		return fmt.Errorf("user service host is required")
	}

	if c.Services.OrderService.Host == "" {
		return fmt.Errorf("order service host is required")
	}

	return nil
}
