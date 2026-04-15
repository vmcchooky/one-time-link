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

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Create Redis-backed secret service
	secretService := secret.NewRedisService(redisClient)

	server := httpapi.NewServer(cfg, secretService)

	log.Printf("starting %s on %s", cfg.ServiceName, cfg.ListenAddress())
	log.Printf("connected to Redis at %s", cfg.RedisAddr)

	if err := http.ListenAndServe(cfg.ListenAddress(), server.Handler()); err != nil {
		log.Fatal(err)
	}
}
