package main

import (
	"github.com/agnaldopidev/rate_limiter/internal/infrastructure/redis"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Configuração do Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	// Configuração de limites por IP
	ipLimit, _ := strconv.Atoi(os.Getenv("IP_LIMIT"))
	ipWindow, _ := time.ParseDuration(os.Getenv("IP_WINDOW"))
	ipBlockDuration, _ := time.ParseDuration(os.Getenv("IP_BLOCK_DURATION"))

	redisLimiter := redis.NewRedisRateLimiter(
		redisAddr,
		redisPass,
		redisDB)

	// Inicializar o middleware de Rate Limiting
	limiter := middleware.NewRateLimiterMiddleware(redisLimiter, ipLimit, ipWindow, ipBlockDuration)

	tokenLimit, _ := strconv.Atoi(os.Getenv("TOKEN_LIMIT"))
	tokenBlockDuration, _ := time.ParseDuration(os.Getenv("TOKEN_BLOCK_DURATION"))
	limiter.SetTokenLimit("example_token", tokenLimit, tokenBlockDuration)

	mux := http.NewServeMux()
	mux.Handle("/", limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Requisição permitida!"))
	})))

	log.Println("Servidor iniciando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
