package repositories

import (
	"context"
	"time"
)

type RateLimitRepository interface {
	CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (allowed bool, remaining int, err error)
}
