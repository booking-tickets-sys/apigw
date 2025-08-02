package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type ServerConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

type HTTPConfig struct {
	Port                    string        `mapstructure:"port"`
	Host                    string        `mapstructure:"host"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}

type ServicesConfig struct {
	UserService UserServiceConfig `mapstructure:"user_service"`
}

type UserServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (c *UserServiceConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func LoadConfig() (*Config, error) {
	// Set default values
	setDefaults()

	// Configure Viper for multiple config sources
	configureViper()

	// Read config from file
	if err := viper.ReadInConfig(); err != nil {
		// If config file is not found, continue with defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func configureViper() {
	// Set config name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add multiple config file paths to search
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/apigw")
	viper.AddConfigPath("$HOME/.apigw")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APIGW")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Allow config file path override via environment variable
	if configFile := viper.GetString("CONFIG_FILE"); configFile != "" {
		viper.SetConfigFile(configFile)
	}
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "api-gateway")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")

	// Server defaults
	viper.SetDefault("server.http.port", "8080")
	viper.SetDefault("server.http.host", "localhost")
	viper.SetDefault("server.http.graceful_shutdown_timeout", "30s")

	// Services defaults
	viper.SetDefault("services.user_service.host", "localhost")
	viper.SetDefault("services.user_service.port", 9090)
}
