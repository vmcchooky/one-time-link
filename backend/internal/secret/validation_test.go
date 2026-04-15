package secret

import (
	"strings"
	"testing"
)

func TestValidateCreateSecretRequest(t *testing.T) {
	t.Run("valid request passes", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdC1jaXBoZXJ0ZXh0",
			Nonce:      "MTIzNDU2Nzg5MDEy", // 12 bytes base64url
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != nil {
			t.Errorf("valid request should pass, got error: %v", err)
		}
	})

	t.Run("rejects invalid algorithm", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdA",
			Nonce:      "MTIzNDU2Nzg5MDEy",
			Algorithm:  "AES-CBC",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrInvalidAlgorithm {
			t.Errorf("expected ErrInvalidAlgorithm, got %v", err)
		}
	})

	t.Run("rejects invalid TTL", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdA",
			Nonce:      "MTIzNDU2Nzg5MDEy",
			Algorithm:  "AES-GCM",
			TTLSeconds: 7200, // not in allowed list
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrInvalidTTL {
			t.Errorf("expected ErrInvalidTTL, got %v", err)
		}
	})

	t.Run("accepts all valid TTL values", func(t *testing.T) {
		validTTLs := []int{3600, 86400, 604800}

		for _, ttl := range validTTLs {
			req := CreateSecretRequest{
				Ciphertext: "dGVzdA",
				Nonce:      "MTIzNDU2Nzg5MDEy",
				Algorithm:  "AES-GCM",
				TTLSeconds: ttl,
			}

			err := ValidateCreateSecretRequest(req)
			if err != nil {
				t.Errorf("TTL %d should be valid, got error: %v", ttl, err)
			}
		}
	})

	t.Run("rejects empty ciphertext", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "",
			Nonce:      "MTIzNDU2Nzg5MDEy",
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrEmptyCiphertext {
			t.Errorf("expected ErrEmptyCiphertext, got %v", err)
		}
	})

	t.Run("rejects ciphertext exceeding size limit", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: strings.Repeat("a", MaxCiphertextSize+1),
			Nonce:      "MTIzNDU2Nzg5MDEy",
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrCiphertextTooLarge {
			t.Errorf("expected ErrCiphertextTooLarge, got %v", err)
		}
	})

	t.Run("rejects empty nonce", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdA",
			Nonce:      "",
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrEmptyNonce {
			t.Errorf("expected ErrEmptyNonce, got %v", err)
		}
	})

	t.Run("rejects nonce with invalid length", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdA",
			Nonce:      "c2hvcnQ", // too short (5 bytes)
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err != ErrInvalidNonceLength {
			t.Errorf("expected ErrInvalidNonceLength, got %v", err)
		}
	})

	t.Run("rejects nonce with invalid base64url encoding", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "dGVzdA",
			Nonce:      "invalid!!!base64",
			Algorithm:  "AES-GCM",
			TTLSeconds: 3600,
		}

		err := ValidateCreateSecretRequest(req)
		if err == nil {
			t.Error("expected error for invalid base64url encoding")
		}
	})
}

func TestDecodeBase64Url(t *testing.T) {
	t.Run("decodes valid base64url", func(t *testing.T) {
		// "hello" in base64url
		input := "aGVsbG8"
		expected := []byte("hello")

		result, err := decodeBase64Url(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(result) != string(expected) {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("handles padding correctly", func(t *testing.T) {
		// Test with different padding scenarios
		tests := []struct {
			input    string
			expected string
		}{
			{"YQ", "a"},        // needs ==
			{"YWI", "ab"},      // needs =
			{"YWJj", "abc"},    // no padding needed
			{"YWJjZA", "abcd"}, // needs ==
		}

		for _, tt := range tests {
			result, err := decodeBase64Url(tt.input)
			if err != nil {
				t.Errorf("input %s: unexpected error: %v", tt.input, err)
				continue
			}

			if string(result) != tt.expected {
				t.Errorf("input %s: expected %s, got %s", tt.input, tt.expected, result)
			}
		}
	})
}
