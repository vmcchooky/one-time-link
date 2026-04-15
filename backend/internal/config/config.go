package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName   string
	Host          string
	Port          string
	AllowedOrigin string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() Config {
	return Config{
		ServiceName:   getEnv("APP_SERVICE_NAME", "one-time-link-api"),
		Host:          getEnv("APP_HOST", "0.0.0.0"),
		Port:          getEnv("APP_PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
	}
}

func (c Config) ListenAddress() string {
	return c.Host + ":" + c.Port
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}

	return fallback
}
