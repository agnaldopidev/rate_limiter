package infrastructure

import (
	"context"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
	"sync"
	"time"
)

type MemoryRateLimiter struct {
	counts map[string]int
	blocks map[string]time.Time
	mu     sync.Mutex
}

func NewMemoryRateLimiter() repositories.RateLimitRepository {
	return &MemoryRateLimiter{
		counts: make(map[string]int),
		blocks: make(map[string]time.Time),
	}
}

func (m *MemoryRateLimiter) CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Lógica de rate limiting aqui
	return true, 0, nil
}
