package secret

import (
	"context"
	"sync"
)

type Service interface {
	Health(ctx context.Context) HealthStatus
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
