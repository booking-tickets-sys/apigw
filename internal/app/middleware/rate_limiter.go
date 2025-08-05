package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// TokenBucketConfig holds token bucket rate limiter configuration
type TokenBucketConfig struct {
	RedisClient    *redis.Client
	Capacity       int           // Maximum number of tokens in the bucket
	RefillRate     float64       // Tokens per second
	RefillInterval time.Duration // How often to refill tokens
	Logger         *logrus.Logger
}

// TokenBucketInfo represents token bucket information
type TokenBucketInfo struct {
	RemainingTokens int           `json:"remaining_tokens"`
	NextRefill      time.Time     `json:"next_refill"`
	Capacity        int           `json:"capacity"`
	RefillRate      float64       `json:"refill_rate"`
	RefillInterval  time.Duration `json:"refill_interval"`
}

// TokenBucket represents a Redis-based token bucket rate limiter
type TokenBucket struct {
	config *TokenBucketConfig
}

// NewTokenBucket creates a new token bucket rate limiter instance
func NewTokenBucket(config *TokenBucketConfig) *TokenBucket {
	return &TokenBucket{
		config: config,
	}
}

// TokenBucketMiddleware creates a token bucket rate limiting middleware
func (tb *TokenBucket) TokenBucketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID)
		clientID := tb.getClientIdentifier(c)

		// Check rate limit using token bucket
		allowed, info, err := tb.checkTokenBucket(c.Request.Context(), clientID)
		if err != nil {
			tb.config.Logger.WithError(err).Error("Token bucket rate limit check failed")
			// On Redis error, allow the request but log the error
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(info.Capacity))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(info.RemainingTokens))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(info.NextRefill.Unix(), 10))
		c.Header("X-RateLimit-RefillRate", fmt.Sprintf("%.2f", info.RefillRate))

		if !allowed {
			tb.config.Logger.WithFields(logrus.Fields{
				"client_id":        clientID,
				"remaining_tokens": info.RemainingTokens,
				"capacity":         info.Capacity,
				"next_refill":      info.NextRefill,
			}).Warn("Token bucket rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "RATE_LIMIT_ERROR",
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Rate limit exceeded. Please try again later.",
				"details": gin.H{
					"remaining_tokens": info.RemainingTokens,
					"next_refill":      info.NextRefill,
					"capacity":         info.Capacity,
					"refill_rate":      info.RefillRate,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkTokenBucket checks if the request is within rate limits using token bucket algorithm
func (tb *TokenBucket) checkTokenBucket(ctx context.Context, clientID string) (bool, *TokenBucketInfo, error) {
	// If Redis client is nil, allow all requests
	if tb.config.RedisClient == nil {
		info := &TokenBucketInfo{
			RemainingTokens: tb.config.Capacity,
			NextRefill:      time.Now().Add(tb.config.RefillInterval),
			Capacity:        tb.config.Capacity,
			RefillRate:      tb.config.RefillRate,
			RefillInterval:  tb.config.RefillInterval,
		}
		return true, info, nil
	}

	now := time.Now()

	// Create keys for this client
	tokensKey := fmt.Sprintf("token_bucket:tokens:%s", clientID)
	lastRefillKey := fmt.Sprintf("token_bucket:last_refill:%s", clientID)

	// Use Redis pipeline for atomic operations
	pipe := tb.config.RedisClient.Pipeline()

	// Get current tokens and last refill time
	tokensCmd := pipe.Get(ctx, tokensKey)
	lastRefillCmd := pipe.Get(ctx, lastRefillKey)

	// Execute pipeline to get current state
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, nil, fmt.Errorf("redis pipeline execution failed: %w", err)
	}

	// Parse current tokens
	var currentTokens int
	if tokensCmd.Val() != "" {
		if val, err := strconv.Atoi(tokensCmd.Val()); err == nil {
			currentTokens = val
		}
	} else {
		currentTokens = tb.config.Capacity // Start with full bucket
	}

	// Parse last refill time
	var lastRefill time.Time
	if lastRefillCmd.Val() != "" {
		if timestamp, err := strconv.ParseInt(lastRefillCmd.Val(), 10, 64); err == nil {
			lastRefill = time.Unix(timestamp, 0)
		} else {
			lastRefill = now
		}
	} else {
		lastRefill = now
	}

	// Calculate time since last refill
	timeSinceLastRefill := now.Sub(lastRefill)

	// Calculate tokens to add based on refill rate and time elapsed
	tokensToAdd := int(tb.config.RefillRate * timeSinceLastRefill.Seconds())

	// Refill the bucket (but don't exceed capacity)
	newTokens := currentTokens + tokensToAdd
	if newTokens > tb.config.Capacity {
		newTokens = tb.config.Capacity
	}

	// Check if we have enough tokens
	if newTokens < 1 {
		// No tokens available, calculate next refill time
		nextRefill := lastRefill.Add(time.Duration(float64(time.Second) * (1.0 / tb.config.RefillRate)))

		info := &TokenBucketInfo{
			RemainingTokens: 0,
			NextRefill:      nextRefill,
			Capacity:        tb.config.Capacity,
			RefillRate:      tb.config.RefillRate,
			RefillInterval:  tb.config.RefillInterval,
		}
		return false, info, nil
	}

	// Consume one token
	newTokens--

	// Update Redis with new token count and refill time
	updatePipe := tb.config.RedisClient.Pipeline()
	updatePipe.Set(ctx, tokensKey, newTokens, 0)      // No expiration for tokens
	updatePipe.Set(ctx, lastRefillKey, now.Unix(), 0) // No expiration for last refill

	_, updateErr := updatePipe.Exec(ctx)
	if updateErr != nil {
		return false, nil, fmt.Errorf("redis update failed: %w", updateErr)
	}

	// Calculate next refill time
	nextRefill := now.Add(time.Duration(float64(time.Second) * (1.0 / tb.config.RefillRate)))

	info := &TokenBucketInfo{
		RemainingTokens: newTokens,
		NextRefill:      nextRefill,
		Capacity:        tb.config.Capacity,
		RefillRate:      tb.config.RefillRate,
		RefillInterval:  tb.config.RefillInterval,
	}

	return true, info, nil
}

// getClientIdentifier returns a unique identifier for the client
func (tb *TokenBucket) getClientIdentifier(c *gin.Context) string {
	// Try to get user ID from JWT context first
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%s", userID)
	}

	// Fall back to IP address
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = "unknown"
	}

	return fmt.Sprintf("ip:%s", clientIP)
}

// CreateCustomTokenBucketMiddleware creates a token bucket rate limiting middleware with custom configuration
func CreateCustomTokenBucketMiddleware(
	redisClient *redis.Client,
	capacity int,
	refillRate float64,
	refillInterval time.Duration,
	logger *logrus.Logger,
) gin.HandlerFunc {
	config := &TokenBucketConfig{
		RedisClient:    redisClient,
		Capacity:       capacity,
		RefillRate:     refillRate,
		RefillInterval: refillInterval,
		Logger:         logger,
	}

	limiter := NewTokenBucket(config)
	return limiter.TokenBucketMiddleware()
}
