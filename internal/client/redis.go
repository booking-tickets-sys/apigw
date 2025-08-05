package client

import (
	"context"
	"fmt"
	"time"

	"apigw/internal/app/config"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisClient represents a Redis client wrapper
type RedisClient struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig, logger *logrus.Logger) (*RedisClient, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("Redis is not enabled")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"host": cfg.Host,
		"port": cfg.Port,
		"db":   cfg.DB,
	}).Info("Redis client connected successfully")

	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

// GetClient returns the underlying Redis client
func (rc *RedisClient) GetClient() *redis.Client {
	return rc.client
}

// Close closes the Redis connection
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

 