package middleware

import (
	"net/http"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/domain"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
)

type RateLimiterMiddleware struct {
	repo    repositories.RateLimitRepository
	ipLimit domain.RateLimit // Config padrão para IP
}

func NewRateLimiterMiddleware(
	repo repositories.RateLimitRepository,
	ipLimit int,
	window time.Duration,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		repo: repo,
		ipLimit: domain.RateLimit{
			Limit:         ipLimit,
			Window:        window,
			BlockDuration: window, // Bloqueio pelo mesmo período da janela
		},
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extrai o IP ou token (prioridade para token)
		key := r.RemoteAddr // IP como fallback
		if token := r.Header.Get("API_KEY"); token != "" {
			key = token
		}

		// 2. Verifica o rate limit
		allowed, remaining, err := m.repo.Allow(r.Context(), key, m.ipLimit.Limit, m.ipLimit.Window)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 3. Resposta se exceder o limite
		if !allowed {
			w.Header().Set("Retry-After", m.ipLimit.BlockDuration.String())
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "you have reached the maximum number of requests"}`))
			return
		}

		// 4. Adiciona headers de rate limit (opcional)
		w.Header().Set("X-RateLimit-Limit", string(m.ipLimit.Limit))
		w.Header().Set("X-RateLimit-Remaining", string(remaining))

		next.ServeHTTP(w, r)
	})
}
