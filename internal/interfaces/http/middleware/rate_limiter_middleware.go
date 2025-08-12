package middleware

import (
	_ "context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/domain"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
)

// RateLimiterMiddleware gerencia a lógica de rate limiting
type RateLimiterMiddleware struct {
	repo        repositories.RateLimitRepository
	defaultIP   domain.RateLimit            // Config padrão para IPs
	tokenConfig map[string]domain.RateLimit // Config por token
	mu          sync.RWMutex                // Protege o map tokenConfig
}

// NewRateLimiterMiddleware cria uma nova instância do middleware
func NewRateLimiterMiddleware(
	repo repositories.RateLimitRepository,
	ipLimit int,
	ipWindow time.Duration,
	ipBlockDuration time.Duration,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		repo: repo,
		defaultIP: domain.RateLimit{
			Limit:         ipLimit,
			Window:        ipWindow,
			BlockDuration: ipBlockDuration,
		},
		tokenConfig: make(map[string]domain.RateLimit),
	}
}

// SetTokenLimit define ou atualiza um limite personalizado para um token
func (m *RateLimiterMiddleware) SetTokenLimit(token string, limit int, blockDuration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokenConfig[token] = domain.RateLimit{
		Limit:         limit,
		Window:        m.defaultIP.Window, // Herda a janela padrão
		BlockDuration: blockDuration,
	}
}

// GetConfig retorna a configuração aplicável para uma chave
func (m *RateLimiterMiddleware) GetConfig(key string) domain.RateLimit {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if config, exists := m.tokenConfig[key]; exists {
		return config
	}
	return m.defaultIP
}

// Handler é o middleware HTTP principal
func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Identifica a chave (IP ou token)
		key := getRequestKey(r)
		config := m.GetConfig(key)

		// 2. Verifica o rate limit
		allowed, remaining, err := m.repo.Allow(
			r.Context(),
			key,
			config.Limit,
			config.Window,
			config.BlockDuration,
		)

		if err != nil {
			handleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 3. Responde se exceder o limite
		if !allowed {
			w.Header().Set("Retry-After", config.BlockDuration.String())
			handleError(w,
				"you have reached the maximum number of requests",
				http.StatusTooManyRequests,
			)
			return
		}

		// 4. Adiciona headers informativos
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		next.ServeHTTP(w, r)
	})
}

// --- Funções auxiliares ---

func getRequestKey(r *http.Request) string {
	if token := r.Header.Get("API_KEY"); token != "" {
		return token //prioriza token
	}
	return r.RemoteAddr //Fallback para IP
}

func handleError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + message + `"}`))
}
