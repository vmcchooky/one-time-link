package secret

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// Validation errors with detailed messages
var (
	ErrInvalidAlgorithm   = errors.New("algorithm must be 'AES-GCM'")
	ErrInvalidTTL         = errors.New("ttl_seconds must be one of: 3600 (1 hour), 86400 (24 hours), or 604800 (7 days)")
	ErrInvalidNonceLength = errors.New("nonce must be exactly 12 bytes when base64-decoded (16 base64url characters)")
	ErrCiphertextTooLarge = errors.New("ciphertext exceeds maximum size of 15KB")
	ErrEmptyCiphertext    = errors.New("ciphertext is required and cannot be empty")
	ErrEmptyNonce         = errors.New("nonce is required and cannot be empty")
	ErrInvalidNonceFormat = errors.New("nonce must be valid base64url-encoded string")
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

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

// ValidateCreateSecretRequest validates the incoming request
// Returns a slice of validation errors for better error reporting
func ValidateCreateSecretRequest(req CreateSecretRequest) error {
	var validationErrors []ValidationError

	// Validate algorithm
	if req.Algorithm != "AES-GCM" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "algorithm",
			Message: "Algorithm must be 'AES-GCM'",
			Code:    "invalid_algorithm",
		})
	}

	// Validate TTL
	if !allowedTTLs[req.TTLSeconds] {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "ttlSeconds",
			Message: "TTL must be one of: 3600 (1 hour), 86400 (24 hours), or 604800 (7 days)",
			Code:    "invalid_ttl",
		})
	}

	// Validate ciphertext
	if req.Ciphertext == "" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "ciphertext",
			Message: "Ciphertext is required and cannot be empty",
			Code:    "empty_ciphertext",
		})
	} else if len(req.Ciphertext) > MaxCiphertextSize {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "ciphertext",
			Message: fmt.Sprintf("Ciphertext exceeds maximum size of %d bytes", MaxCiphertextSize),
			Code:    "ciphertext_too_large",
		})
	}

	// Validate nonce
	if req.Nonce == "" {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "nonce",
			Message: "Nonce is required and cannot be empty",
			Code:    "empty_nonce",
		})
	} else {
		// Decode and validate nonce length
		nonceBytes, err := decodeBase64Url(req.Nonce)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "nonce",
				Message: "Nonce must be a valid base64url-encoded string",
				Code:    "invalid_nonce_format",
			})
		} else if len(nonceBytes) != NonceLength {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "nonce",
				Message: fmt.Sprintf("Nonce must be exactly %d bytes when decoded (16 base64url characters)", NonceLength),
				Code:    "invalid_nonce_length",
			})
		}
	}

	// Return errors if any
	if len(validationErrors) > 0 {
		return &MultiValidationError{Errors: validationErrors}
	}

	return nil
}

// MultiValidationError represents multiple validation errors
type MultiValidationError struct {
	Errors []ValidationError
}

func (e *MultiValidationError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Message
	}
	return fmt.Sprintf("%d validation errors occurred", len(e.Errors))
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
