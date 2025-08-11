package redis

import (
	"context"
	"fmt"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(addr string, password string, db int) repositories.RateLimitRepository {
	return &RedisRateLimiter{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisRateLimiter) Allow(
	ctx context.Context,
	key string,
	limit int,
	window time.Duration,
) (bool, int, error) {
	// 1. Verifica se está bloqueado
	blockKey := fmt.Sprintf("block:%s", key)
	if blocked, _ := r.client.Get(ctx, blockKey).Bool(); blocked {
		return false, 0, nil
	}

	// 2. Incrementa o contador
	countKey := fmt.Sprintf("count:%s", key)
	current, err := r.client.Incr(ctx, countKey).Result()
	if err != nil {
		return false, 0, err
	}

	// 3. Seta expiração na primeira requisição
	if current == 1 {
		r.client.Expire(ctx, countKey, window)
	}

	// 4. Verifica o limite
	remaining := limit - int(current)
	if remaining < 0 {
		remaining = 0
	}

	if int(current) > limit {
		// Bloqueia a chave
		r.client.Set(ctx, blockKey, 1, window)
		r.client.Del(ctx, countKey) // Reseta o contador
		return false, remaining, nil
	}

	return true, remaining, nil
}
