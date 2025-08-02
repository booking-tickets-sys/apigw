package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"apigw/config"
	"apigw/internal/client"
	"apigw/internal/handler"
	"apigw/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize gRPC client for user service
	userServiceAddress := cfg.Services.UserService.GetAddress()
	userClient, err := client.NewUserServiceClient(userServiceAddress)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
	defer userClient.Close()

	// Initialize handlers
	userHandler := handler.NewUserHandler(userClient)

	// Setup router
	router := router.SetupRouter(userHandler)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API Gateway server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.HTTP.GracefulShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway server exited")
}
