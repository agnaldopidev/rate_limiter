// internal/interfaces/repositories/rate_limit_repository.go
package repositories

import (
	"context"
	"time"
)

// RateLimitRepository define o contrato para um repositório de rate limiting
type RateLimitRepository interface {
	// Allow verifica se uma requisição é permitida
	Allow(
		ctx context.Context,
		key string, // Identificador (IP/token)
		limit int, // Limite de requisições
		window time.Duration, // Janela de tempo (ex: 1s, 1m)
	) (allowed bool, remaining int, err error)
}
