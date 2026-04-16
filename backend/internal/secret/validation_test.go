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
		if err == nil {
			t.Error("expected validation error for invalid algorithm")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "algorithm" && valErr.Code == "invalid_algorithm" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected algorithm validation error")
			}
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
		if err == nil {
			t.Error("expected validation error for invalid TTL")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "ttlSeconds" && valErr.Code == "invalid_ttl" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected TTL validation error")
			}
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
		if err == nil {
			t.Error("expected validation error for empty ciphertext")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "ciphertext" && valErr.Code == "empty_ciphertext" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected ciphertext validation error")
			}
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
		if err == nil {
			t.Error("expected validation error for ciphertext too large")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "ciphertext" && valErr.Code == "ciphertext_too_large" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected ciphertext size validation error")
			}
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
		if err == nil {
			t.Error("expected validation error for empty nonce")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "nonce" && valErr.Code == "empty_nonce" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected nonce validation error")
			}
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
		if err == nil {
			t.Error("expected validation error for invalid nonce length")
		}

		// Check if it's a MultiValidationError
		if multiErr, ok := err.(*MultiValidationError); ok {
			found := false
			for _, valErr := range multiErr.Errors {
				if valErr.Field == "nonce" && valErr.Code == "invalid_nonce_length" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected nonce length validation error")
			}
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

	t.Run("returns multiple validation errors", func(t *testing.T) {
		req := CreateSecretRequest{
			Ciphertext: "",        // empty
			Nonce:      "",        // empty
			Algorithm:  "INVALID", // invalid
			TTLSeconds: 999,       // invalid
		}

		err := ValidateCreateSecretRequest(req)
		if err == nil {
			t.Fatal("expected validation errors")
		}

		// Check if it's a MultiValidationError with multiple errors
		if multiErr, ok := err.(*MultiValidationError); ok {
			if len(multiErr.Errors) != 4 {
				t.Errorf("expected 4 validation errors, got %d", len(multiErr.Errors))
			}
		} else {
			t.Error("expected MultiValidationError")
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
