// internal/infrastructure/memory/memory_rate_limiter.go
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
)

// MemoryRateLimiter implementa RateLimitRepository usando armazenamento em memória
type MemoryRateLimiter struct {
	counts map[string]int       // Contador de requisições por chave
	blocks map[string]time.Time // Registro de bloqueios por chave
	mu     sync.Mutex           // Mutex para evitar concorrência
}

// NewMemoryRateLimiter cria uma nova instância do limitador em memória
func NewMemoryRateLimiter() repositories.RateLimitRepository {
	return &MemoryRateLimiter{
		counts: make(map[string]int),
		blocks: make(map[string]time.Time),
	}
}

// Allow verifica se uma requisição é permitida
func (m *MemoryRateLimiter) Allow(
	ctx context.Context,
	key string,
	limit int,
	window time.Duration,
) (bool, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. Verifica se está bloqueado
	if blockTime, exists := m.blocks[key]; exists {
		if time.Since(blockTime) < window {
			return false, 0, nil // Ainda bloqueado
		}
		delete(m.blocks, key) // Remove bloqueio expirado
	}

	// 2. Incrementa o contador ou reinicia a janela
	m.counts[key]++
	if m.counts[key] == 1 {
		// Simula expiração da janela (em um sistema real, usaríamos goroutine ou timer)
		go func() {
			time.Sleep(window)
			m.mu.Lock()
			delete(m.counts, key)
			m.mu.Unlock()
		}()
	}

	// 3. Verifica o limite
	remaining := limit - m.counts[key]
	if remaining < 0 {
		remaining = 0
	}

	if m.counts[key] > limit {
		// Bloqueia a chave
		m.blocks[key] = time.Now()
		delete(m.counts, key) // Reseta o contador
		return false, remaining, nil
	}

	return true, remaining, nil
}
