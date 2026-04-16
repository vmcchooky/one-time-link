package main

import (
	"log"
	"net/http"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/httpapi"
	"one-time-link/backend/internal/secret"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// Initialize Redis client with connection pooling
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		PoolSize:     cfg.RedisPoolSize,
		MinIdleConns: cfg.RedisMinIdle,
		MaxRetries:   cfg.RedisMaxRetries,
	})

	// Create Redis-backed secret service
	secretService := secret.NewRedisService(redisClient)

	// Create server with rate limiting
	server := httpapi.NewServerWithRateLimiting(cfg, secretService, redisClient)

	log.Printf("starting %s on %s", cfg.ServiceName, cfg.ListenAddress())
	log.Printf("Redis: %s (pool: %d, min idle: %d, max retries: %d)",
		cfg.RedisAddr, cfg.RedisPoolSize, cfg.RedisMinIdle, cfg.RedisMaxRetries)
	log.Printf("Rate limiting enabled: create=10/hr, consume=20/hr, status=100/hr")

	if err := http.ListenAndServe(cfg.ListenAddress(), server.Handler()); err != nil {
		log.Fatal(err)
	}
}
