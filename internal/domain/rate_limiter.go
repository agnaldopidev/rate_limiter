package domain

import "time"

// RateLimit representa a configuração básica do limitador
type RateLimit struct {
	Key           string        // Identificador (IP ou token)
	Limit         int           // Número máximo de requisições
	Window        time.Duration // Janela de tempo (ex: 1s, 1m)
	BlockDuration time.Duration // Tempo de bloqueio se exceder
}

// Result representa o resultado de uma verificação de limite
type Result struct {
	Allowed    bool          // Se a requisição é permitida
	Remaining  int           // Requisições restantes
	RetryAfter time.Duration // Tempo até poder fazer novas requisições
}
