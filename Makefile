# Makefile for Go-Web project
# This Makefile helps to run both backend and frontend services

# Variables
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BINARY_NAME := go-web-app

# Default target
.PHONY: help
help: ## Show this help
	@echo "Go-Web Project - Makefile Commands"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: dev
dev: ## Run both backend and frontend in development mode
	@echo "Starting both backend and frontend development servers..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Frontend will be available at http://localhost:3000"
	@echo ""
	@echo "Starting backend..."
	@cd $(BACKEND_DIR) && go run main.go &
	@echo "Starting frontend..."
	@cd $(FRONTEND_DIR) && npm run dev

.PHONY: backend
backend: ## Run backend service
	@echo "Starting backend service..."
	@cd $(BACKEND_DIR) && go run main.go

.PHONY: frontend
frontend: ## Run frontend development server
	@echo "Starting frontend development server..."
	@cd $(FRONTEND_DIR) && npm run dev

##@ Build

.PHONY: build
build: build-backend build-frontend ## Build both backend and frontend

.PHONY: build-backend
build-backend: ## Build backend binary
	@echo "Building backend..."
	@cd $(BACKEND_DIR) && go build -o ../$(BINARY_NAME) .

.PHONY: build-frontend
build-frontend: ## Build frontend
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm run build

##@ Testing

.PHONY: test
test: test-backend test-frontend ## Run tests for both backend and frontend

.PHONY: test-backend
test-backend: ## Run backend tests
	@echo "Running backend tests..."
	@cd $(BACKEND_DIR) && go test ./...

.PHONY: test-frontend
test-frontend: ## Run frontend tests
	@echo "Running frontend tests..."
	@cd $(FRONTEND_DIR) && npm run test

##@ Setup and Dependencies

.PHONY: setup
setup: setup-backend setup-frontend ## Setup both backend and frontend dependencies

.PHONY: setup-backend
setup-backend: ## Setup backend dependencies
	@echo "Setting up backend dependencies..."
	@cd $(BACKEND_DIR) && go mod tidy

.PHONY: setup-frontend
setup-frontend: ## Setup frontend dependencies
	@echo "Setting up frontend dependencies..."
	@cd $(FRONTEND_DIR) && npm install

##@ Database

.PHONY: migrate
migrate: ## Run database migrations (if any)
	@echo "Running database migrations..."
	@echo "Note: GORM AutoMigrate runs automatically when starting the backend"
	@cd $(BACKEND_DIR) && go run main.go --migrate

##@ Clean

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@cd $(FRONTEND_DIR) && rm -rf dist/

.PHONY: clean-deps
clean-deps: ## Clean dependencies
	@echo "Cleaning dependencies..."
	@cd $(BACKEND_DIR) && go clean -modcache
	@cd $(FRONTEND_DIR) && rm -rf node_modules/

##@ Docker (if you add Docker support later)

.PHONY: docker-build
docker-build: ## Build Docker images for both services
	@echo "Building Docker images..."
	@echo "Note: Dockerfiles need to be created first"
	#@docker build -t go-web-backend -f backend/Dockerfile backend/
	#@docker build -t go-web-frontend -f frontend/Dockerfile frontend/

.PHONY: docker-run
docker-run: ## Run services with Docker
	@echo "Running services with Docker..."
	@echo "Note: docker-compose.yml needs to be created first"
	#@docker-compose up

##@ Utility

.PHONY: lint
lint: lint-backend lint-frontend ## Run linters on both backend and frontend

.PHONY: lint-backend
lint-backend: ## Run linters on backend code
	@echo "Linting backend code..."
	@cd $(BACKEND_DIR) && go vet ./...

.PHONY: lint-frontend
lint-frontend: ## Run linters on frontend code
	@echo "Linting frontend code..."
	@cd $(FRONTEND_DIR) && npm run lint

.PHONY: format
format: format-backend format-frontend ## Format code in both backend and frontend

.PHONY: format-backend
format-backend: ## Format backend code
	@echo "Formatting backend code..."
	@cd $(BACKEND_DIR) && go fmt ./...

.PHONY: format-frontend
format-frontend: ## Format frontend code
	@echo "Formatting frontend code..."
	@cd $(FRONTEND_DIR) && npm run format
