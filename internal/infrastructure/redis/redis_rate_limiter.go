package redis

import (
	"context"
	"fmt"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
	"time"

	_ "github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
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
	blockDuration time.Duration, // Novo parâmetro
) (bool, int, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	countKey := fmt.Sprintf("count:%s", key)

	// Verifica bloqueio
	if ttl, _ := r.client.TTL(ctx, blockKey).Result(); ttl > 0 {
		return false, 0, nil
	}

	// Incrementa contador
	current, err := r.client.Incr(ctx, countKey).Result()
	if err != nil {
		return false, 0, err
	}

	// Seta expiração
	if current == 1 {
		r.client.Expire(ctx, countKey, window)
	}

	remaining := limit - int(current)
	if remaining < 0 {
		remaining = 0
	}

	if int(current) > limit {
		r.client.Set(ctx, blockKey, 1, blockDuration) // Usa blockDuration
		return false, remaining, nil
	}

	return true, remaining, nil
}
