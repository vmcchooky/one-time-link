package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisService implements secret storage using Redis
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new Redis-backed secret service
func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
	}
}

// CreateSecret stores an encrypted secret in Redis with TTL
func (s *RedisService) CreateSecret(ctx context.Context, req CreateSecretRequest) (*CreateSecretResponse, error) {
	// Generate unique secret ID
	secretID := uuid.New().String()

	// Calculate expiration time
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(req.TTLSeconds) * time.Second)

	// Create secret record
	secret := Secret{
		SecretID:   secretID,
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
		Algorithm:  req.Algorithm,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		TTLSeconds: req.TTLSeconds,
		ConsumedAt: nil,
	}

	// Serialize to JSON
	data, err := json.Marshal(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal secret: %w", err)
	}

	// Store in Redis with TTL
	key := fmt.Sprintf("secret:%s", secretID)
	err = s.client.Set(ctx, key, data, time.Duration(req.TTLSeconds)*time.Second).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store secret in Redis: %w", err)
	}

	return &CreateSecretResponse{
		SecretID:  secretID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// Health checks Redis connectivity
func (s *RedisService) Health(ctx context.Context) HealthStatus {
	err := s.client.Ping(ctx).Err()
	if err != nil {
		return HealthStatus{
			Store: "redis",
			Mode:  "unhealthy",
		}
	}

	return HealthStatus{
		Store: "redis",
		Mode:  "healthy",
	}
}
