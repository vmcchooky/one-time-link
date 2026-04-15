package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

const apiBaseURL = "http://localhost:8080"

type CreateSecretRequest struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Algorithm  string `json:"algorithm"`
	TTLSeconds int    `json:"ttlSeconds"`
}

type CreateSecretResponse struct {
	SecretID  string `json:"secretId"`
	ExpiresAt string `json:"expiresAt"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

func TestMilestone2Comprehensive(t *testing.T) {
	// Check if backend is running
	resp, err := http.Get(apiBaseURL + "/healthz")
	if err != nil {
		t.Skip("Backend not running, skipping integration tests")
		return
	}
	resp.Body.Close()

	t.Run("Size Limits", func(t *testing.T) {
		t.Run("small plaintext should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1zbWFsbC1wbGFpbnRleHQ",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("medium plaintext ~5KB should pass", func(t *testing.T) {
			mediumText := strings.Repeat("a", 5000)
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: mediumText,
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("large plaintext ~10KB should pass", func(t *testing.T) {
			largeText := strings.Repeat("a", 10000)
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: largeText,
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("too large plaintext >15KB should fail", func(t *testing.T) {
			tooLargeText := strings.Repeat("a", 16000)
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: tooLargeText,
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 413)
		})
	})

	t.Run("TTL Validation", func(t *testing.T) {
		t.Run("TTL 3600 should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC10dGwtMzYwMA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("TTL 86400 should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC10dGwtODY0MDA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 86400,
			}, 201)
		})

		t.Run("TTL 604800 should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC10dGwtNjA0ODAw",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 604800,
			}, 201)
		})

		t.Run("invalid TTL 7200 should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1pbnZhbGlkLXR0bA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 7200,
			}, 400)
		})

		t.Run("zero TTL should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC16ZXJvLXR0bA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 0,
			}, 400)
		})

		t.Run("negative TTL should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1uZWdhdGl2ZS10dGw",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: -1,
			}, 400)
		})
	})

	t.Run("Algorithm Validation", func(t *testing.T) {
		t.Run("AES-GCM should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1hZXMtZ2Nt",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("AES-CBC should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1hZXMtY2Jj",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-CBC",
				TTLSeconds: 3600,
			}, 400)
		})

		t.Run("ChaCha20 should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1jaGFjaGEyMA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "ChaCha20",
				TTLSeconds: 3600,
			}, 400)
		})

		t.Run("empty algorithm should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1lbXB0eS1hbGdv",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "",
				TTLSeconds: 3600,
			}, 400)
		})
	})

	t.Run("Nonce Validation", func(t *testing.T) {
		t.Run("12-byte nonce should pass", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC12YWxpZC1ub25jZQ",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 201)
		})

		t.Run("8-byte nonce should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1zaG9ydC1ub25jZQ",
				Nonce:      "MTIzNDU2Nzg",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 400)
		})

		t.Run("16-byte nonce should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1sb25nLW5vbmNl",
				Nonce:      "MTIzNDU2Nzg5MDEyMzQ1Ng",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 400)
		})

		t.Run("empty nonce should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "dGVzdC1lbXB0eS1ub25jZQ",
				Nonce:      "",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 400)
		})
	})

	t.Run("Ciphertext Validation", func(t *testing.T) {
		t.Run("empty ciphertext should fail", func(t *testing.T) {
			testCreateSecret(t, CreateSecretRequest{
				Ciphertext: "",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: 3600,
			}, 400)
		})
	})
}

func testCreateSecret(t *testing.T, req CreateSecretRequest, expectedStatus int) {
	t.Helper()

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", apiBaseURL+"/api/secrets", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", fmt.Sprintf("test-%d", time.Now().Unix()))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	if resp.StatusCode == 201 {
		var result CreateSecretResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result.SecretID == "" {
			t.Error("secretId should not be empty")
		}

		if result.ExpiresAt == "" {
			t.Error("expiresAt should not be empty")
		}

		t.Logf("Created secret: %s, expires: %s", result.SecretID, result.ExpiresAt)
	}
}
