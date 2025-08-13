git package redis2_test

import (
	"context"
	"testing"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/redis"
	"github.com/stretchr/testify/assert" // Auxilia nos asserts
)

// Configuração inicial para testes
var (
	ctx           = context.Background()
	redisAddr     = "localhost:6379" // Certifique-se de que o Redis esteja rodando
	redisPassword = ""
	redisDB       = 0
	limiter       = redis.NewRedisRateLimiter(redisAddr, redisPassword, redisDB)
)

func TestRedisRateLimiter_Allowed(t *testing.T) {
	// Configuração do teste
	key := "test-key"
	limit := 5
	window := time.Second // Janela de 1 segundo
	blockDuration := 2 * time.Second

	// Cenário 1: Permitir requisições abaixo do limite
	for i := 1; i <= limit; i++ {
		allowed, remaining, err := limiter.Allow(ctx, key, limit, window, blockDuration)
		assert.NoError(t, err, "Erro inesperado ao verificar limite")
		assert.True(t, allowed, "Requisição deveria ser permitida")
		assert.Equal(t, limit-i, remaining, "Requisições restantes estão incorretas")
	}

	// Cenário 2: Rejeitar após atingir limite
	allowed, remaining, err := limiter.Allow(ctx, key, limit, window, blockDuration)
	assert.NoError(t, err)
	assert.False(t, allowed, "Requisição deveria ser bloqueada")
	assert.Equal(t, 0, remaining, "Não deveria haver requisições restantes")
}

func TestRedisRateLimiter_BlockAfterLimit(t *testing.T) {
	// Configuração do teste
	key := "block-test-key"
	limit := 3
	window := time.Second            // Janela de 1 segundo
	blockDuration := 3 * time.Second // Bloqueio de 3 segundos

	// Atingir o limite
	for i := 1; i <= limit; i++ {
		allowed, _, err := limiter.Allow(ctx, key, limit, window, blockDuration)
		assert.NoError(t, err)
		assert.True(t, allowed, "Requisição deveria ser permitida")
	}

	// Ultrapassar o limite
	allowed, _, err := limiter.Allow(ctx, key, limit, window, blockDuration)
	assert.NoError(t, err)
	assert.False(t, allowed, "Requisição deveria ser bloqueada")

	// Garantir que ainda está bloqueado
	time.Sleep(1 * time.Second)
	allowed, _, err = limiter.Allow(ctx, key, limit, window, blockDuration)
	assert.NoError(t, err)
	assert.False(t, allowed, "Requisição deveria continuar bloqueada")

	// Após o bloqueio expirar
	time.Sleep(3 * time.Second)
	allowed, _, err = limiter.Allow(ctx, key, limit, window, blockDuration)
	assert.NoError(t, err)
	assert.True(t, allowed, "Requisição deveria ser permitida após o bloqueio expirar")
}

func TestRedisRateLimiter_ResetAfterTTL(t *testing.T) {
	// Configuração do teste
	key := "ttl-reset-test-key"
	limit := 3
	window := 2 * time.Second // Janela de 2 segundos
	blockDuration := 3 * time.Second

	// Atingir o limite
	for i := 1; i <= limit; i++ {
		allowed, _, err := limiter.Allow(ctx, key, limit, window, blockDuration)
		assert.NoError(t, err)
		assert.True(t, allowed, "Requisição deveria ser permitida")
	}

	// Esperar o TTL expirar
	time.Sleep(window)

	// Requisições devem ser permitidas novamente após o reset
	allowed, _, err := limiter.Allow(ctx, key, limit, window, blockDuration)
	assert.NoError(t, err)
	assert.True(t, allowed, "Requisição deveria ser permitida após o reset do TTL")
}

func TestRedisRateLimiter_IndependentKeys(t *testing.T) {
	// Configuração do teste
	key1 := "independent-key-1"
	key2 := "independent-key-2"
	limit := 3
	window := time.Second // Janela de 1 segundo
	blockDuration := 2 * time.Second

	// Atingir o limite para a chave 1
	for i := 1; i <= limit; i++ {
		allowed, _, err := limiter.Allow(ctx, key1, limit, window, blockDuration)
		assert.NoError(t, err)
		assert.True(t, allowed, "Requisição deveria ser permitida")
	}

	// Garantir que a chave 2 ainda está dentro do limite
	for i := 1; i <= limit; i++ {
		allowed, _, err := limiter.Allow(ctx, key2, limit, window, blockDuration)
		assert.NoError(t, err)
		assert.True(t, allowed, "Requisição deveria ser permitida para chave 2")
	}
}
