// internal/infrastructure/memory/memory_rate_limiter.go
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
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

func (m *MemoryRateLimiter) Allow(
	ctx context.Context,
	key string,
	limit int,
	window time.Duration,
	blockDuration time.Duration, // Novo parâmetro
) (bool, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verifica bloqueio existente
	if blockTime, exists := m.blocks[key]; exists {
		if time.Since(blockTime) < blockDuration { // Usa blockDuration aqui
			return false, 0, nil
		}
		delete(m.blocks, key)
	}

	// Lógica de contagem
	m.counts[key]++
	if m.counts[key] == 1 {
		go func() {
			time.Sleep(window)
			m.mu.Lock()
			delete(m.counts, key)
			m.mu.Unlock()
		}()
	}

	remaining := limit - m.counts[key]
	if remaining < 0 {
		remaining = 0
	}

	if m.counts[key] > limit {
		m.blocks[key] = time.Now()
		delete(m.counts, key)
		return false, remaining, nil
	}

	return true, remaining, nil
}
