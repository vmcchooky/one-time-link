package config

import "os"

type Config struct {
	ServiceName   string
	Host          string
	Port          string
	AllowedOrigin string
}

func Load() Config {
	return Config{
		ServiceName:   getEnv("APP_SERVICE_NAME", "one-time-link-api"),
		Host:          getEnv("APP_HOST", "0.0.0.0"),
		Port:          getEnv("APP_PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
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
