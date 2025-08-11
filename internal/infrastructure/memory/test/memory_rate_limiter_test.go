package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRateLimiter_Allow(t *testing.T) {
	ctx := context.Background()
	limiter := memory.NewMemoryRateLimiter()
	key := "test-key"
	limit := 5
	window := 1 * time.Second

	t.Run("Permite requisições dentro do limite", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			allowed, remaining, err := limiter.Allow(ctx, key, limit, window)
			assert.True(t, allowed, "Deveria permitir a requisição %d", i+1)
			assert.Equal(t, limit-(i+1), remaining, "Restantes incorretos na requisição %d", i+1)
			assert.NoError(t, err)
		}
	})

	t.Run("Bloqueia requisições acima do limite", func(t *testing.T) {
		allowed, remaining, err := limiter.Allow(ctx, key, limit, window)
		assert.False(t, allowed, "Deveria bloquear a requisição acima do limite")
		assert.Equal(t, 0, remaining, "Deveria mostrar 0 requisições restantes")
		assert.NoError(t, err)
	})

	t.Run("Libera após o bloqueio expirar", func(t *testing.T) {
		// Espera a janela de tempo expirar
		time.Sleep(window + 100*time.Millisecond)

		allowed, remaining, err := limiter.Allow(ctx, key, limit, window)
		assert.True(t, allowed, "Deveria permitir após o bloqueio expirar")
		assert.Equal(t, limit-1, remaining, "Deveria reiniciar o contador")
		assert.NoError(t, err)
	})
}
