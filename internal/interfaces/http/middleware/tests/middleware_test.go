package middleware_test

import (
	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/redis"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/middleware"
)

func TestRateLimiterMiddleware(t *testing.T) {
	// Mock do repositório

	// Configuração do Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	repo := redis.NewRedisRateLimiter(
		redisAddr,
		redisPass,
		redisDB)

	// Configura middleware
	ipLimit := 2
	blockDuration := time.Second * 5
	limiter := middleware.NewRateLimiterMiddleware(repo, ipLimit, time.Second, blockDuration)

	// Configura handler de teste
	handler := limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Requisição permitida!"))
	}))

	// Cria requisição de teste
	req := httptest.NewRequest("GET", "/", nil)

	// Verifica o limite de IP
	for i := 1; i <= ipLimit; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("esperava HTTP 200, recebeu: %d", rr.Code)
		}
	}

	// Testa se o terceiro acesso resulta em bloqueio
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("esperava HTTP 429, recebeu: %d", rr.Code)
	}
}
