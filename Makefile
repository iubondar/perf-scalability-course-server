.PHONY: build test test-all run clean deploy check-platform help lint vet fmt

# Переменные
BINARY_NAME=perf-server
MAIN_PATH=./cmd/server
BUILD_DIR=./bin

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать бинарник для текущей платформы
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: ## Собрать бинарники для Linux (amd64 и arm64)
	@echo "Building Linux binaries..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "Binaries built:"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-*

test: vet ## Запустить тесты (включает go vet)
	@echo "Running tests..."
	go test -v ./...

test-all: lint test ## Запустить все проверки (lint + vet + test)

test-coverage: vet ## Запустить тесты с покрытием (включает go vet)
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

run: ## Запустить сервер локально
	@echo "Running server..."
	go run $(MAIN_PATH)

clean: ## Очистить артефакты сборки
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f perf-server-linux-*
	@echo "Clean complete"

deploy: ## Деплой на удалённую VM (требует DEPLOY_HOST, DEPLOY_USER)
	@if [ -z "$(DEPLOY_HOST)" ] || [ -z "$(DEPLOY_USER)" ]; then \
		echo "Error: DEPLOY_HOST and DEPLOY_USER must be set"; \
		echo "Usage: make deploy DEPLOY_HOST=example.com DEPLOY_USER=ubuntu"; \
		exit 1; \
	fi
	@./scripts/deploy.sh $(DEPLOY_HOST) $(DEPLOY_USER) $(DEPLOY_PATH)

check-platform: ## Проверить платформу удалённой VM
	@if [ -z "$(DEPLOY_HOST)" ] || [ -z "$(DEPLOY_USER)" ]; then \
		echo "Error: DEPLOY_HOST and DEPLOY_USER must be set"; \
		echo "Usage: make check-platform DEPLOY_HOST=example.com DEPLOY_USER=ubuntu"; \
		exit 1; \
	fi
	@./scripts/check-platform.sh $(DEPLOY_HOST) $(DEPLOY_USER)

lint: ## Запустить линтер (если установлен golangci-lint)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
	fi

fmt: ## Форматировать код
	@echo "Formatting code..."
	go fmt ./...

vet: ## Запустить go vet (автоматически запускается с test)
	@echo "Running go vet..."
	@go vet ./...
