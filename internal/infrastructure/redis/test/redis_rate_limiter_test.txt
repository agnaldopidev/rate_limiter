package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRateLimiter(t *testing.T) {
	ctx := context.Background()
	rl := memory.NewMemoryRateLimiter()
	key := "test-key"

	t.Run("deve respeitar o tempo de bloqueio", func(t *testing.T) {
		// Primeira requisição (deve passar)
		allowed, remaining, err := rl.Allow(ctx, key, 1, time.Second, 2*time.Second)
		assert.True(t, allowed)
		assert.Equal(t, 0, remaining)
		assert.Nil(t, err)

		// Segunda requisição (deve bloquear)
		allowed, remaining, err = rl.Allow(ctx, key, 1, time.Second, 2*time.Second)
		assert.False(t, allowed)
		assert.Equal(t, 0, remaining)
		assert.Nil(t, err)

		// Espera 1s (ainda bloqueado)
		time.Sleep(1 * time.Second)
		allowed, _, _ = rl.Allow(ctx, key, 1, time.Second, 2*time.Second)
		assert.False(t, allowed)

		// Espera mais 1.1s (total 2.1s > blockDuration)
		time.Sleep(1100 * time.Millisecond)
		allowed, remaining, err = rl.Allow(ctx, key, 1, time.Second, 2*time.Second)
		assert.True(t, allowed)
		assert.Equal(t, 0, remaining)
		assert.Nil(t, err)
	})
}
