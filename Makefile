.PHONY: help build up down test test-backend test-frontend clean logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker containers
	docker compose build

up: ## Start the application
	docker compose up -d
	@echo "Application started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080/api"

down: ## Stop the application
	docker compose down

logs: ## View application logs
	docker compose logs -f

test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests
	cd backend && go test ./... -v

test-frontend: ## Run frontend E2E tests
	cd frontend && npm install && npm run test:e2e

clean: ## Clean up containers, volumes, and build artifacts
	docker compose down -v
	rm -rf backend/data/*.db
	rm -rf frontend/node_modules
	rm -rf frontend/dist

restart: down up ## Restart the application

dev-backend: ## Run backend in development mode
	cd backend && go run main.go

dev-frontend: ## Run frontend in development mode
	cd frontend && npm install && npm run dev

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

install-backend: ## Install backend dependencies
	cd backend && go mod download

format-backend: ## Format backend code
	cd backend && go fmt ./...

format-frontend: ## Format frontend code
	cd frontend && npm run format || echo "Add prettier to package.json for formatting"
