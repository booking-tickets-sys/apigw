package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"apigw/internal/app/config"
	"apigw/internal/app/router"
	"apigw/internal/client"
	"apigw/pkg/utils/crypt/token"
	logutils "apigw/pkg/utils/log"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	if err := logutils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger := logutils.GetLogger()

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Configuration validation failed: %v", err)
	}

	// Create clients
	userClient, err := client.NewUserServiceClient(&cfg.Services.UserService)
	if err != nil {
		logger.Fatalf("Failed to create user client: %v", err)
	}
	orderClient, err := client.NewOrderServiceClient(&cfg.Services.OrderService)
	if err != nil {
		logger.Fatalf("Failed to create order client: %v", err)
	}

	// Initialize Redis client for rate limiting
	var redisClient *client.RedisClient
	if cfg.Redis.Enabled {
		redisClient, err = client.NewRedisClient(&cfg.Redis, logger)
		if err != nil {
			logger.Fatalf("Failed to create Redis client: %v", err)
		}
		defer redisClient.Close()
		logger.Info("Redis client initialized for rate limiting")
	} else {
		logger.Info("Redis is disabled, rate limiting will not be available")
	}

	// Ensure clients are properly closed on exit
	defer func() {
		if userClient != nil {
			if err := userClient.Close(); err != nil {
				logger.WithError(err).Error("Failed to close user client")
			}
		}
		if orderClient != nil {
			if err := orderClient.Close(); err != nil {
				logger.WithError(err).Error("Failed to close order client")
			}
		}
	}()

	// Initialize token maker
	tokenMaker, err := token.NewJWTTokenMaker(cfg.JWT.SecretKey)
	if err != nil {
		logger.Fatalf("Failed to create token maker: %v", err)
	}

	// Setup router
	router := router.SetupRouter(cfg, userClient, orderClient, redisClient, tokenMaker, logger)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"address":     serverAddr,
			"port":        cfg.Server.HTTP.Port,
			"host":        cfg.Server.HTTP.Host,
			"environment": cfg.App.Environment,
			"version":     cfg.App.Version,
		}).Info("API Gateway server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.HTTP.GracefulShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("API Gateway server exited")
}
