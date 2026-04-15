package secret

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestRedisServiceIntegration tests the Redis service with a real Redis instance
// This test requires Redis to be running on localhost:6379
// Skip this test if Redis is not available
func TestRedisServiceIntegration(t *testing.T) {
	// Try to connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// Check if Redis is available
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}
	defer client.Close()

	service := NewRedisService(client)

	t.Run("creates secret successfully", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "test-ciphertext-base64url",
			Nonce:      "test-nonce-12b",
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		resp, err := service.CreateSecret(ctx, req)
		if err != nil {
			t.Fatalf("CreateSecret failed: %v", err)
		}

		if resp.SecretID == "" {
			t.Error("SecretID should not be empty")
		}

		if resp.ExpiresAt == "" {
			t.Error("ExpiresAt should not be empty")
		}

		// Verify secret exists in Redis
		key := "secret:" + resp.SecretID
		exists, err := client.Exists(ctx, key).Result()
		if err != nil {
			t.Fatalf("Failed to check Redis key: %v", err)
		}

		if exists != 1 {
			t.Error("Secret should exist in Redis")
		}

		// Verify TTL is set
		ttl, err := client.TTL(ctx, key).Result()
		if err != nil {
			t.Fatalf("Failed to get TTL: %v", err)
		}

		if ttl <= 0 || ttl > time.Hour {
			t.Errorf("TTL should be between 0 and 1 hour, got %v", ttl)
		}

		// Cleanup
		client.Del(ctx, key)
	})

	t.Run("health check returns healthy when Redis is available", func(t *testing.T) {
		status := service.Health(ctx)

		if status.Store != "redis" {
			t.Errorf("expected store 'redis', got '%s'", status.Store)
		}

		if status.Mode != "healthy" {
			t.Errorf("expected mode 'healthy', got '%s'", status.Mode)
		}
	})

	t.Run("secret expires after TTL", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "test-expiring-secret",
			Nonce:      "test-nonce-exp",
			Algorithm:  "AES-GCM",
			TTLSeconds: 1, // 1 second for quick test
		}

		resp, err := service.CreateSecret(ctx, req)
		if err != nil {
			t.Fatalf("CreateSecret failed: %v", err)
		}

		key := "secret:" + resp.SecretID

		// Verify secret exists
		exists, _ := client.Exists(ctx, key).Result()
		if exists != 1 {
			t.Error("Secret should exist immediately after creation")
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Verify secret is gone
		exists, _ = client.Exists(ctx, key).Result()
		if exists != 0 {
			t.Error("Secret should be expired and removed from Redis")
		}
	})
}
