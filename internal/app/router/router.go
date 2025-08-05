package router

import (
	"apigw/internal/app/config"
	"apigw/internal/app/handler"
	"apigw/internal/app/middleware"
	"apigw/internal/client"
	"apigw/pkg/utils/crypt/token"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(
	cfg *config.Config,
	userClient *client.UserServiceClient,
	orderClient *client.OrderServiceClient,
	redisClient *client.RedisClient,
	jwtMaker *token.JWTMaker,
	logger *logrus.Logger,
) *gin.Engine {
	// Set Gin mode
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware(logger))

	// Add token bucket rate limiter middleware if Redis is available
	if redisClient != nil {
		tokenBucketMiddleware := middleware.CreateCustomTokenBucketMiddleware(
			redisClient.GetClient(),
			cfg.Redis.TokenBucket.Capacity,
			cfg.Redis.TokenBucket.RefillRate,
			cfg.Redis.TokenBucket.RefillInterval,
			logger,
		)
		router.Use(tokenBucketMiddleware)
		logger.WithFields(logrus.Fields{
			"capacity":        cfg.Redis.TokenBucket.Capacity,
			"refill_rate":     cfg.Redis.TokenBucket.RefillRate,
			"refill_interval": cfg.Redis.TokenBucket.RefillInterval,
		}).Info("Token bucket rate limiter middleware enabled")
	} else {
		logger.Info("Token bucket rate limiter middleware disabled (Redis not available)")
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Create handlers
	userHandler := handler.NewUserHandler(userClient, logger)
	orderHandler := handler.NewOrderHandler(orderClient, logger)

	// Create JWT middleware
	jwtMiddleware := middleware.JWTMiddleware(jwtMaker, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes (no authentication required)
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/refresh", userHandler.RefreshToken)
		}

		// Order routes (authentication required)
		orders := api.Group("/orders")
		orders.Use(jwtMiddleware)
		{
			orders.POST("/:event_id/purchase", orderHandler.PurchaseTicket)
		}
	}

	return router
}
