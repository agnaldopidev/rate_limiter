package main

import (
	"net/http"
	"time"

	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/memory"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/middleware"
	_ "github.com/agnaldopidev/rate_limiter/internal/interfaces/repositories"
)

func main() {
	// 1. Inicializa o rate limiter em memória
	repo := memory.NewMemoryRateLimiter()

	// 2. Configura o middleware (5 reqs/segundo por IP)
	middleware := middleware.NewRateLimiterMiddleware(
		repo,
		5,           // Limite padrão (IP)
		time.Second, // Janela de tempo
	)
	middleware.SetTokenLimit("abc123", 100)
	middleware.SetTokenLimit("free-tier", 2)
	// 3. Roteador
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// 4. Aplica o middleware
	handler := middleware.Handler(mux)

	// 5. Inicia o servidor
	http.ListenAndServe(":8080", handler)
}
