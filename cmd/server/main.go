package main

import (
	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/redis"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/handlers"
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
		5,
		time.Second,
		30*time.Second, // Janela de tempo
	)
	//	middleware.SetTokenLimit("premium", 100)
	//	middleware.SetTokenLimit("free", 2)

	configHandler := handlers.NewConfigHandler(middleware)

	// 3. Roteador
	mux := http.NewServeMux()
	mux.HandleFunc("/config", configHandler.UpdateConfig) // Novo endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Redis Rate Limiter!"))
	})

	// 4. Aplica o middleware
	handler := middleware.Handler(mux)

	// 5. Inicia o servidor
	http.ListenAndServe(":8080", handler)
}
