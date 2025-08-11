package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/domain"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
)

type RateLimiterMiddleware struct {
	repo        repositories.RateLimitRepository
	ipLimit     domain.RateLimit // Config padrão para IP
	tokenLimits map[string]int   // Limites personalizados por token
	mu          sync.Mutex
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
		tokenLimits: make(map[string]int),
	}
}

// Adiciona ou atualiza um limite personalizado para um token
func (m *RateLimiterMiddleware) SetTokenLimit(token string, limit int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokenLimits[token] = limit
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Identifica a chave (IP ou token)
		key := r.RemoteAddr
		limit := m.ipLimit.Limit
		window := m.ipLimit.Window

		if token := r.Header.Get("API_KEY"); token != "" {
			m.mu.Lock()
			if customLimit, exists := m.tokenLimits[token]; exists {
				limit = customLimit // Usa o limite personalizado do token
			}
			m.mu.Unlock()
			key = token
		}

		// 2. Verifica o rate limit
		allowed, remaining, err := m.repo.Allow(r.Context(), key, limit, window)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 3. Responde com 429 se exceder o limite
		if !allowed {
			w.Header().Set("Retry-After", window.String())
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limit exceeded", "limit": ` + string(limit) + `}`))
			return
		}

		// 4. Headers opcionais para debug
		w.Header().Set("X-RateLimit-Limit", string(limit))
		w.Header().Set("X-RateLimit-Remaining", string(remaining))

		next.ServeHTTP(w, r)
	})
}
