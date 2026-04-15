package secret

import (
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidAlgorithm   = errors.New("algorithm must be AES-GCM")
	ErrInvalidTTL         = errors.New("ttl must be 3600, 86400, or 604800 seconds")
	ErrInvalidNonceLength = errors.New("nonce must be exactly 12 bytes (16 base64url characters)")
	ErrCiphertextTooLarge = errors.New("ciphertext exceeds maximum size")
	ErrEmptyCiphertext    = errors.New("ciphertext cannot be empty")
	ErrEmptyNonce         = errors.New("nonce cannot be empty")
)

const (
	MaxCiphertextSize = 15 * 1024 // 15KB (includes base64 overhead)
	NonceLength       = 12        // 12 bytes for AES-GCM
)

var allowedTTLs = map[int]bool{
	3600:   true, // 1 hour
	86400:  true, // 24 hours
	604800: true, // 7 days
}

// ValidateCreateSecretRequest validates the incoming request
func ValidateCreateSecretRequest(req CreateSecretRequest) error {
	// Validate algorithm
	if req.Algorithm != "AES-GCM" {
		return ErrInvalidAlgorithm
	}

	// Validate TTL
	if !allowedTTLs[req.TTLSeconds] {
		return ErrInvalidTTL
	}

	// Validate ciphertext
	if req.Ciphertext == "" {
		return ErrEmptyCiphertext
	}

	if len(req.Ciphertext) > MaxCiphertextSize {
		return ErrCiphertextTooLarge
	}

	// Validate nonce
	if req.Nonce == "" {
		return ErrEmptyNonce
	}

	// Decode and validate nonce length
	nonceBytes, err := decodeBase64Url(req.Nonce)
	if err != nil {
		return fmt.Errorf("invalid nonce encoding: %w", err)
	}

	if len(nonceBytes) != NonceLength {
		return ErrInvalidNonceLength
	}

	return nil
}

// decodeBase64Url decodes base64url string to bytes
func decodeBase64Url(s string) ([]byte, error) {
	// Add padding if needed
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	// Use URL encoding (base64url uses - and _ instead of + and /)
	return base64.URLEncoding.DecodeString(s)
}
