package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter implements rate limiting using Redis
type Limiter struct {
	client *redis.Client
}

// Config holds rate limit configuration for an endpoint
type Config struct {
	Limit  int           // Maximum requests allowed
	Window time.Duration // Time window for the limit
}

// Result represents the result of a rate limit check
type Result struct {
	Allowed    bool      // Whether the request is allowed
	Limit      int       // Maximum requests allowed
	Remaining  int       // Remaining requests in current window
	ResetAt    time.Time // When the limit resets
	RetryAfter int       // Seconds to wait before retrying (if not allowed)
}

// NewLimiter creates a new rate limiter
func NewLimiter(client *redis.Client) *Limiter {
	return &Limiter{
		client: client,
	}
}

// Allow checks if a request is allowed under the rate limit
func (l *Limiter) Allow(ctx context.Context, key string, config Config) (*Result, error) {
	// Create Redis key
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Increment counter
	count, err := l.client.Incr(ctx, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// Set expiry on first request
	if count == 1 {
		err = l.client.Expire(ctx, redisKey, config.Window).Err()
		if err != nil {
			return nil, fmt.Errorf("failed to set rate limit expiry: %w", err)
		}
	}

	// Get TTL for reset time
	ttl, err := l.client.TTL(ctx, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit TTL: %w", err)
	}

	// Calculate reset time
	resetAt := time.Now().Add(ttl)

	// Check if limit exceeded
	allowed := count <= int64(config.Limit)
	remaining := config.Limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	result := &Result{
		Allowed:    allowed,
		Limit:      config.Limit,
		Remaining:  remaining,
		ResetAt:    resetAt,
		RetryAfter: int(ttl.Seconds()),
	}

	return result, nil
}

// Reset clears the rate limit for a key (useful for testing)
func (l *Limiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	return l.client.Del(ctx, redisKey).Err()
}
