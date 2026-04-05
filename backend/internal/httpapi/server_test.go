package httpapi

import (
	"net/http"
	"net/http/httptest"
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
