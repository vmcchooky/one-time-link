package secret

import (
	"context"
	"fmt"
	"sync"
)

type Service interface {
	Health(ctx context.Context) HealthStatus
	CreateSecret(ctx context.Context, req CreateSecretRequest) (*CreateSecretResponse, error)
}

type HealthStatus struct {
	Store string `json:"store"`
	Mode  string `json:"mode"`
}

type InMemoryService struct {
	mu sync.RWMutex
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{}
}

func (s *InMemoryService) Health(ctx context.Context) HealthStatus {
	_ = ctx

	s.mu.RLock()
	defer s.mu.RUnlock()

	return HealthStatus{
		Store: "in-memory placeholder",
		Mode:  "scaffold",
	}
}

func (s *InMemoryService) CreateSecret(ctx context.Context, req CreateSecretRequest) (*CreateSecretResponse, error) {
	_ = ctx
	_ = req
	return nil, fmt.Errorf("in-memory service does not support CreateSecret")
}
