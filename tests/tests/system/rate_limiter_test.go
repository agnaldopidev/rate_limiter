package system

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

type TokenConfig struct {
	Token         string        `json:"token"`
	Limit         int           `json:"limit"`
	BlockDuration time.Duration `json:"block_duration_ms"`
}

func TestRateLimiterSystem(t *testing.T) {
	token := "test-system-token-" + time.Now().Format("150405") // Token único

	t.Run("deve respeitar tempo de bloqueio configurado", func(t *testing.T) {
		// 1. Configura o token via API
		err := setTokenLimit(token, 2, 3*time.Second)
		assert.NoError(t, err)

		// 2. Testa o limite
		for i := 0; i < 2; i++ {
			resp, err := doRequest(token)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Req %d deveria passar", i+1)
		}

		// 3. Verifica o bloqueio
		resp, err := doRequest(token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		assert.Equal(t, "3s", resp.Header.Get("Retry-After"))

		// 4. Verifica expiração do bloqueio
		time.Sleep(3*time.Second + 100*time.Millisecond)
		resp, err = doRequest(token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// Implementação real da função
func setTokenLimit(token string, limit int, blockDuration time.Duration) error {
	config := TokenConfig{
		Token:         token,
		Limit:         limit,
		BlockDuration: blockDuration,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		baseURL+"/config",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func doRequest(token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", baseURL+"/api", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("API_KEY", token)
	return http.DefaultClient.Do(req)
}
