package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/secret"
)

func TestGetSecretStatusEndpoint(t *testing.T) {
	cfg := config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}

	t.Run("returns pending status for existing secret", func(t *testing.T) {
		mockService := &mockSecretService{
			getSecretStatusFunc: func(ctx context.Context, secretID string) (*secret.SecretStatus, error) {
				return &secret.SecretStatus{
					SecretID:  secretID,
					Status:    "pending",
					CreatedAt: "2026-04-15T12:00:00Z",
					ExpiresAt: "2026-04-15T13:00:00Z",
				}, nil
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/secrets/test-secret-123/status", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var resp secret.SecretStatus
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "pending" {
			t.Errorf("expected status 'pending', got '%s'", resp.Status)
		}

		if resp.SecretID != "test-secret-123" {
			t.Errorf("expected secretId 'test-secret-123', got '%s'", resp.SecretID)
		}
	})

	t.Run("returns 404 for non-existent secret", func(t *testing.T) {
		mockService := &mockSecretService{
			getSecretStatusFunc: func(ctx context.Context, secretID string) (*secret.SecretStatus, error) {
				return &secret.SecretStatus{
					SecretID: secretID,
					Status:   "not_found",
					Message:  "Secret not found or has expired.",
				}, nil
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/secrets/non-existent/status", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}

		var resp secret.SecretStatus
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "not_found" {
			t.Errorf("expected status 'not_found', got '%s'", resp.Status)
		}
	})

	t.Run("returns 500 on service error", func(t *testing.T) {
		mockService := &mockSecretService{
			getSecretStatusFunc: func(ctx context.Context, secretID string) (*secret.SecretStatus, error) {
				return nil, errors.New("redis connection failed")
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/secrets/test-secret-123/status", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for empty secret ID", func(t *testing.T) {
		mockService := &mockSecretService{}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/secrets/status", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		// Empty secret ID results in route not found
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})
}

func TestConsumeSecretEndpoint(t *testing.T) {
	cfg := config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}

	t.Run("consumes secret successfully", func(t *testing.T) {
		mockService := &mockSecretService{
			consumeSecretFunc: func(ctx context.Context, secretID string) (*secret.ConsumeSecretResponse, error) {
				return &secret.ConsumeSecretResponse{
					SecretID:   secretID,
					Ciphertext: "dGVzdC1jaXBoZXJ0ZXh0",
					Nonce:      "MTIzNDU2Nzg5MDEy",
					Algorithm:  "AES-GCM",
					ConsumedAt: "2026-04-15T12:30:00Z",
				}, nil
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/secrets/test-secret-123/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var resp secret.ConsumeSecretResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.SecretID != "test-secret-123" {
			t.Errorf("expected secretId 'test-secret-123', got '%s'", resp.SecretID)
		}

		if resp.Ciphertext == "" {
			t.Error("ciphertext should not be empty")
		}

		if resp.Nonce == "" {
			t.Error("nonce should not be empty")
		}

		if resp.Algorithm != "AES-GCM" {
			t.Errorf("expected algorithm 'AES-GCM', got '%s'", resp.Algorithm)
		}
	})

	t.Run("returns 410 for already consumed secret", func(t *testing.T) {
		mockService := &mockSecretService{
			consumeSecretFunc: func(ctx context.Context, secretID string) (*secret.ConsumeSecretResponse, error) {
				return nil, errors.New("secret not found or already consumed")
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/secrets/test-secret-123/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusGone {
			t.Errorf("expected status 410, got %d", rec.Code)
		}

		var resp map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp["error"] != "already_consumed" {
			t.Errorf("expected error 'already_consumed', got '%s'", resp["error"])
		}
	})

	t.Run("returns 500 on service error", func(t *testing.T) {
		mockService := &mockSecretService{
			consumeSecretFunc: func(ctx context.Context, secretID string) (*secret.ConsumeSecretResponse, error) {
				return nil, errors.New("redis connection failed")
			},
		}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/secrets/test-secret-123/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for empty secret ID", func(t *testing.T) {
		mockService := &mockSecretService{}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/secrets/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		// Empty secret ID results in route not found
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("includes CORS headers", func(t *testing.T) {
		mockService := &mockSecretService{}
		server := NewServer(cfg, mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/secrets/test-secret-123/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		origin := rec.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:5173" {
			t.Errorf("expected CORS origin 'http://localhost:5173', got '%s'", origin)
		}
	})
}

func TestSecretRoutesMethodValidation(t *testing.T) {
	cfg := config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}

	mockService := &mockSecretService{}
	server := NewServer(cfg, mockService)

	t.Run("status endpoint rejects POST", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/secrets/test-secret-123/status", strings.NewReader("{}"))
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("consume endpoint rejects GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/secrets/test-secret-123/consume", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})
}
