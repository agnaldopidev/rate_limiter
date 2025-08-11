package main

import (
	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/redis"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/middleware"
	"net/http"
	"time"
)

func main() {
	redisLimiter := redis.NewRedisRateLimiter(
		"localhost:6379",
		"",
		0)

	// 2. Configura o middleware (5 reqs/segundo por IP)
	middleware := middleware.NewRateLimiterMiddleware(
		redisLimiter,
		5,           // Limite padrão (IP)
		time.Second, // Janela de tempo
	)
	middleware.SetTokenLimit("premium", 100)
	middleware.SetTokenLimit("free", 2)
	// 3. Roteador
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Redis Rate Limiter!"))
	})

	// 4. Aplica o middleware
	handler := middleware.Handler(mux)

	// 5. Inicia o servidor
	http.ListenAndServe(":8080", handler)
}
