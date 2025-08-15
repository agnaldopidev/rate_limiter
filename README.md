# Rate Limiter em Go

Um **Rate Limiter** desenvolvido em Go para limitar o número de requisições por endereço IP ou Token de acesso, utilizando **Redis** como mecanismo de persistência. Este projeto fornece um middleware personalizável para controlar o tráfego em servidores web, garantindo que o limite máximo de requisições seja respeitado.

---

## 📋 Recursos

- **Limitação por IP:** Controle baseado no número de requisições permitidas por IP por segundo.
- **Limitação por Token de Acesso:** Controle baseado em tokens customizados, com um limite configurável separado.
- **Configuração Flexível:** Configurações de limites e bloqueios via `.env`.
- **Mensagens de Bloqueio Personalizadas:** Respostas HTTP `429 - Too Many Requests`.
- **Persistência Redis:** Utilização do Redis com suporte a tempo de expiração para armazenar contadores e blocos.
- **Estrategia Abstrata:** Possibilidade de trocar o Redis por outro mecanismo de persistência.
- **Docker-Compose:** Facilita o deploy do Redis em ambiente local.

---

## 🚀 Como Executar o Projeto

### 1. **Pré-requisitos**
Certifique-se de ter os seguintes itens instalados:
- [Go 1.23+](https://golang.org/dl/)
- [Docker e Docker Compose](https://docs.docker.com/get-docker/)

---

### 2. **Clonar o Repositório**
Clone o projeto em sua máquina:
```bash
  git clone https://github.com/agnaldopidev/rate_limiter
  cd rate_limiter
```

---

### 3. **Configurar o Ambiente**

Crie um arquivo `.env` na raiz do projeto com as seguintes configurações:

```plaintext
# Configurações do Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Configurações do Rate Limiter
RATE_LIMIT_IP=5              # Limite de requisições por segundo por IP
RATE_LIMIT_TOKEN=10          # Limite de requisições por token por segundo
BLOCK_TIME=300               # Tempo de bloqueio (em segundos) se o limite for excedido
```

---

### 4. **Subir o Redis com Docker**

Utilize o Docker Compose para iniciar o servidor Redis:

```bash
  docker-compose up -d
```

Essa configuração utiliza a porta padrão `6379`. Certifique-se de que não há outros serviços ocupando essa porta.

---

### 5. **Executar o Projeto**

Compile e execute o servidor:

```bash
   go run cmd/server/main.go
```

O servidor será iniciado na porta `8080`.

---

## 📡 Testando o Rate Limiter

### 1. Requisições com Limite por IP

Envie requisições simples (exemplo para Linux ou Mac):

```bash
  for i in {1..10}; do curl -i http://localhost:8080; done
```

Resultado esperado:
- As primeiras 5 requisições são aceitas com código HTTP `200`.
- A sexta e demais retornam código HTTP `429` com a mensagem:
  ```plaintext
  you have reached the maximum number of requests or actions allowed within a certain time frame
  ```

---

### 2. Requisições com Token de Acesso

Caso seu limite por Token seja maior que o de IP, insira o header `API_KEY` no momento da requisição:

```bash
  for i in {1..15}; do curl -i -H "API_KEY: my_token" http://localhost:8080; done
```

Resultado esperado:
- As primeiras 10 requisições (de acordo com seu limite para tokens) são aceitas.
- A 11ª externa bloqueio com resposta HTTP `429`.

---

### 3 Teste com Python
```bash
  python tests/teste_python_rate_limiter/1_test_requisicao_token.py
  python tests/teste_python_rate_limiter/2_test_rate_limite.py
```

## 🛠 Estrutura do Projeto

```plaintext
.
├── docker-compose.yml      # Subir o Redis via Docker Compose
├── go.mod                  # Dependências do projeto Go
├── go.sum                  # Checksum das dependências
├── main.go                 # Inicialização do servidor e middleware
├── .env                    # Configurações do projeto
├── cmd
│    └── server
├── internal
│   ├── domain          
│   ├── infrastructure   
│   └── middleware.go      
│          
└── tests
    └── rete_limiter_test     # Testes automatizados para o Rate Limiter
```

---

## 🔍 Tecnologias Utilizadas

- **Linguagem:** [Go 1.23+](https://golang.org)
- **Banco de Dados:** [Redis](https://redis.io/)
- **Containerização:** [Docker/Docker Compose](https://www.docker.com/)
- **Gerenciamento de Dependências:** [Go Modules](https://github.com/golang/go/wiki/Modules)

---

## ✨ Futuras Melhorias

- **Limites Dinâmicos:** Endpoint HTTP para alterar limites de requisições sem reiniciar a aplicação.
- **Cache Distribuído:** Suporte a múltiplos nós Redis para alta disponibilidade.
- **Taxas Personalizadas:** Diferenças entre usuários gratuitos e premium.
- **Melhora nos Testes:** Adicionar casos de teste para cenários extremos.

---

## 🧪 Testes Unitários

Os testes automatizados verificam o comportamento do Redis e do middleware. Para rodar os testes:

```bash
  go test ./tests
```

Certifique-se de que o Redis está em funcionamento antes de executar os testes.

---

## 📄 Licença

Este projeto é disponibilizado sob a licença **MIT**. Consulte o arquivo `LICENSE` para mais detalhes.

---

## 🤝 Créditos & Contribuições

Desenvolvido por [Agnaldo Correia](https://github.com/seu-usuario).  
Contribuições são bem-vindas! Envie sugestões ou abra um pull request para melhorar o projeto.

---
