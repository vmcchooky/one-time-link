package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/secret"
)

func TestHealthEndpoint(t *testing.T) {
	server := NewServer(config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}, secret.NewInMemoryService())

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected CORS origin header to be set, got %q", got)
	}

	// Check response structure
	var response healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Service != "test-api" {
		t.Errorf("expected service name 'test-api', got %q", response.Service)
	}

	if response.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %q", response.Status)
	}
}

func TestOptionsRequestReturnsNoContent(t *testing.T) {
	server := NewServer(config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}, secret.NewInMemoryService())

	req := httptest.NewRequest(http.MethodOptions, "/api/secrets", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestRequestIDHeader(t *testing.T) {
	server := NewServer(config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}, secret.NewInMemoryService())

	t.Run("generates request ID when not provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		requestID := rec.Header().Get("X-Request-ID")
		if requestID == "" {
			t.Error("expected X-Request-ID header to be set")
		}
	})

	t.Run("echoes provided request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.Header.Set("X-Request-ID", "test-request-id-123")
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		requestID := rec.Header().Get("X-Request-ID")
		if requestID != "test-request-id-123" {
			t.Errorf("expected X-Request-ID to be 'test-request-id-123', got %q", requestID)
		}
	})
}

func TestRequestSizeLimit(t *testing.T) {
	server := NewServer(config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}, secret.NewInMemoryService())

	t.Run("accepts request within size limit", func(t *testing.T) {
		body := bytes.NewBufferString(`{"test": "data"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/secrets", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		// Should return 501 Not Implemented, not 413
		if rec.Code == http.StatusRequestEntityTooLarge {
			t.Error("small request should not be rejected for size")
		}
	})

	t.Run("rejects request exceeding size limit", func(t *testing.T) {
		// Create a payload larger than 15KB
		largePayload := strings.Repeat("a", 16*1024)
		body := bytes.NewBufferString(largePayload)
		req := httptest.NewRequest(http.MethodPost, "/api/secrets", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("expected status 413, got %d", rec.Code)
		}
	})
}

func TestCORSHeaders(t *testing.T) {
	server := NewServer(config.Config{
		ServiceName:   "test-api",
		Host:          "127.0.0.1",
		Port:          "8080",
		AllowedOrigin: "http://localhost:5173",
	}, secret.NewInMemoryService())

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	headers := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Expose-Headers",
	}

	for _, header := range headers {
		if rec.Header().Get(header) == "" {
			t.Errorf("expected %s header to be set", header)
		}
	}
}
