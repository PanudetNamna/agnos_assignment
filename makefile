# Makefile

.PHONY: up down restart build logs ps \
        test lint \
        db-reset db-connect \
        mock-his-up \
        help

# ==================== DOCKER ====================

up:
	docker compose up --build -d

down:
	docker compose down

restart:
	docker compose down && docker compose up --build -d

build:
	docker compose build

logs:
	docker compose logs -f

logs-app:
	docker compose logs -f app

logs-db:
	docker compose logs -f postgres

ps:
	docker compose ps

# ==================== DEV ====================

run:
	go run ./cmd/main.go

test:
	go test ./... -v

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

# ==================== DATABASE ====================

db-reset:
	docker compose down -v
	docker compose up --build -d
	@echo "database reset complete"

db-connect:
	docker exec -it agnos_postgres psql -U postgres -d agnos

# ==================== MOCK HIS ====================

mock-his-up:
	go run ./cmd/mock_his/main.go

# ==================== CURL ====================

curl-health:
	curl -s http://localhost/health | jq

curl-create-staff:
	curl -s -X POST http://localhost/staff/create \
		-H "Content-Type: application/json" \
		-d '{"username": "staff_a", "password": "secret", "hospital": "Hospital A"}' | jq

curl-login:
	curl -s -X POST http://localhost/staff/login \
		-H "Content-Type: application/json" \
		-d '{"username": "staff_a", "password": "secret", "hospital": "Hospital A"}' | jq

curl-search: ## TOKEN=<jwt> make curl-search
	curl -s -X POST http://localhost/patient/search \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $(TOKEN)" \
		-d '{"national_id": "1234567890123"}' | jq

# ==================== HELP ====================

help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Docker:"
	@echo "  up              build and start all services"
	@echo "  down            stop all services"
	@echo "  restart         down + up"
	@echo "  build           build images"
	@echo "  logs            follow all logs"
	@echo "  logs-app        follow app logs"
	@echo "  logs-db         follow postgres logs"
	@echo "  ps              show running containers"
	@echo ""
	@echo "Dev:"
	@echo "  run             run app locally"
	@echo "  test            run all tests"
	@echo "  lint            run linter"
	@echo "  tidy            go mod tidy"
	@echo ""
	@echo "Database:"
	@echo "  db-reset        drop volume and recreate"
	@echo "  db-connect      connect to postgres shell"
	@echo ""
	@echo "Mock HIS:"
	@echo "  mock-his-up     run mock HIS locally on :9090"
	@echo ""
	@echo "Curl:"
	@echo "  curl-health     GET /health"
	@echo "  curl-login      POST /staff/login"
	@echo "  curl-search     POST /patient/search  (TOKEN=<jwt> make curl-search)"
	@echo ""
