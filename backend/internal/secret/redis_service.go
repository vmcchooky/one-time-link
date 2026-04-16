package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	// RedisOperationTimeout is the default timeout for Redis operations
	RedisOperationTimeout = 5 * time.Second
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
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

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
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

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

// GetSecretStatus checks if a secret exists and returns its status
func (s *RedisService) GetSecretStatus(ctx context.Context, secretID string) (*SecretStatus, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	// Get secret from Redis
	key := fmt.Sprintf("secret:%s", secretID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Secret not found or expired
		return &SecretStatus{
			SecretID: secretID,
			Status:   "not_found",
			Message:  "Secret not found or has expired.",
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from Redis: %w", err)
	}

	// Deserialize secret
	var secret Secret
	err = json.Unmarshal(data, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	// Return pending status
	return &SecretStatus{
		SecretID:  secretID,
		Status:    "pending",
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
		ExpiresAt: secret.ExpiresAt.Format(time.RFC3339),
	}, nil
}

// ConsumeSecret atomically retrieves and deletes a secret
func (s *RedisService) ConsumeSecret(ctx context.Context, secretID string) (*ConsumeSecretResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	// Use GETDEL for atomic read-and-delete
	key := fmt.Sprintf("secret:%s", secretID)
	data, err := s.client.GetDel(ctx, key).Bytes()
	if err == redis.Nil {
		// Secret not found, expired, or already consumed
		return nil, fmt.Errorf("secret not found or already consumed")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to consume secret from Redis: %w", err)
	}

	// Deserialize secret
	var secret Secret
	err = json.Unmarshal(data, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	// Return consumed secret data
	now := time.Now().UTC()
	return &ConsumeSecretResponse{
		SecretID:   secretID,
		Ciphertext: secret.Ciphertext,
		Nonce:      secret.Nonce,
		Algorithm:  secret.Algorithm,
		ConsumedAt: now.Format(time.RFC3339),
	}, nil
}
