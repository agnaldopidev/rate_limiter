package repositories

import (
	"context"
	"time"
)

type RateLimitRepository interface {
	Allow(
		ctx context.Context,
		key string, // Identificador (IP/token)
		limit int, // Limite de requisições
		window time.Duration, // Janela de tempo
		blockDuration time.Duration, // Tempo de bloqueio
	) (allowed bool, remaining int, err error)
}
