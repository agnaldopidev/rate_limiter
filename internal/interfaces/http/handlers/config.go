package handlers

import (
	"encoding/json"
	"github.com/agnaldopidev/rate_limiter/internal/interfaces/http/middleware"
	"net/http"
	"time"
)

type ConfigHandler struct {
	rateLimiterMiddleware *middleware.RateLimiterMiddleware
}

type TokenConfig struct {
	Token         string        `json:"token"`
	Limit         int           `json:"limit"`
	BlockDuration time.Duration `json:"block_duration_ms"` // Novo campo
}

func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config TokenConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Agora passando todos os três parâmetros
	h.rateLimiterMiddleware.SetTokenLimit(
		config.Token,
		config.Limit,
		config.BlockDuration, // Novo parâmetro
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "config updated"}`))
}

func NewConfigHandler(rateLimiter *middleware.RateLimiterMiddleware) *ConfigHandler {
	return &ConfigHandler{
		rateLimiterMiddleware: rateLimiter,
	}
}
