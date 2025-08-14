import requests
import time
import unittest

BASE_URL = "http://localhost:8080"
DEFAULT_LIMIT = 5  # Limite padrão por IP
TOKEN_LIMITS = {
    "premium": 100,  # Token com limite alto
    "free": 2        # Token com limite baixo
}

class TestRateLimiter(unittest.TestCase):
    def setUp(self):
        self.session = requests.Session()
    
    def test_default_ip_limit(self):
        """Testa o limite padrão para requisições sem token"""
        print("\nTestando limite padrão por IP (5 reqs/segundo):")
        
        # Primeiras 5 requisições devem passar
        for i in range(DEFAULT_LIMIT):
            response = self.session.get(BASE_URL)
            self.assertEqual(response.status_code, 200, f"Req {i+1} deveria passar (status 200)")
            print(f"Req {i+1}: HTTP {response.status_code}")

        # A 6ª deve falhar
        response = self.session.get(BASE_URL)
        self.assertEqual(response.status_code, 429, "Deveria bloquear após o limite")
        print(f"Req 6: HTTP {response.status_code} (bloqueado como esperado)")

    def test_token_limits(self):
        """Testa limites personalizados por token"""
        print("\nTestando limites por token:")
        
        # Testa token 'free' com limite 2
        print(f"\nToken 'free' (limite: {TOKEN_LIMITS['free']}):")
        headers = {"API_KEY": "free"}
        for i in range(TOKEN_LIMITS["free"] + 1):
            response = self.session.get(BASE_URL, headers=headers)
            if i < TOKEN_LIMITS["free"]:
                self.assertEqual(response.status_code, 200, f"Req {i+1} deveria passar")
            else:
                self.assertEqual(response.status_code, 429, "Deveria bloquear após o limite")
            print(f"Req {i+1}: HTTP {response.status_code}")

        # Testa token 'premium' com limite 100
        print(f"\nToken 'premium' (limite: {TOKEN_LIMITS['premium']}):")
        headers = {"API_KEY": "premium"}
        for i in range(TOKEN_LIMITS["premium"] + 1):
            response = self.session.get(BASE_URL, headers=headers)
            if i < TOKEN_LIMITS["premium"]:
                self.assertEqual(response.status_code, 200, f"Req {i+1} deveria passar")
            else:
                self.assertEqual(response.status_code, 429, "Deveria bloquear após o limite")
            
            if (i+1) % 20 == 0 or i+1 == TOKEN_LIMITS["premium"] + 1:
                print(f"Req {i+1}: HTTP {response.status_code}")

    def test_block_expiration(self):
        """Verifica se o bloqueio expira após a janela de tempo"""
        print("\nTestando expiração do bloqueio:")
        
        # Excede o limite padrão
        for _ in range(DEFAULT_LIMIT + 1):
            response = self.session.get(BASE_URL)
        
        # Verifica o bloqueio
        self.assertEqual(response.status_code, 429)
        print(f"Status após exceder limite: HTTP {response.status_code}")

        # Espera a janela de tempo expirar (1 segundo)
        print(f"Aguardando {1} segundo para expirar o bloqueio...")
        time.sleep(1)

        # Nova requisição deve passar
        response = self.session.get(BASE_URL)
        self.assertEqual(response.status_code, 200)
        print(f"Status após expiração: HTTP {response.status_code}")

if __name__ == "__main__":
    unittest.main(verbosity=2)