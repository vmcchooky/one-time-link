package secret

import "time"

// Secret represents the encrypted secret data stored in Redis
type Secret struct {
	SecretID   string     `json:"secretId"`
	Ciphertext string     `json:"ciphertext"`
	Nonce      string     `json:"nonce"`
	Algorithm  string     `json:"algorithm"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	TTLSeconds int        `json:"ttlSeconds"`
	ConsumedAt *time.Time `json:"consumedAt,omitempty"`
}

// CreateSecretRequest represents the incoming request to create a secret
type CreateSecretRequest struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Algorithm  string `json:"algorithm"`
	TTLSeconds int    `json:"ttlSeconds"`
}

// CreateSecretResponse represents the response after creating a secret
type CreateSecretResponse struct {
	SecretID  string `json:"secretId"`
	ExpiresAt string `json:"expiresAt"`
}

// SecretStatus represents the status of a secret
type SecretStatus struct {
	SecretID  string `json:"secretId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// ConsumeSecretResponse represents the response when consuming a secret
type ConsumeSecretResponse struct {
	SecretID   string `json:"secretId"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Algorithm  string `json:"algorithm"`
	ConsumedAt string `json:"consumedAt"`
}
