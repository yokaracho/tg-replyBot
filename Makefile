.PHONY: build run test clean docker-build docker-run lint fmt deps generate help

# Переменные
APP_NAME=telegram-reply-bot
BINARY_NAME=bot
DOCKER_IMAGE=$(APP_NAME):latest
MAIN_PATH=./cmd/bot
BUILD_DIR=./bin

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

## Основные команды

help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(NC) %s\n", $1, $2 }' $(MAKEFILE_LIST)

build: ## Собрать приложение
	@echo "$(GREEN)Сборка приложения...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Сборка завершена: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: ## Запустить приложение
	@echo "$(GREEN)Запуск приложения...$(NC)"
	@go run $(MAIN_PATH)

run-local: build ## Запустить собранное приложение локально
	@echo "$(GREEN)Запуск собранного приложения...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

## Тестирование

test: ## Запустить все тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	@go test -v ./...

test-cover: ## Запустить тесты с покрытием
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии создан: coverage.html$(NC)"

test-race: ## Запустить тесты с проверкой гонок
	@echo "$(GREEN)Запуск тестов с проверкой гонок...$(NC)"
	@go test -race -v ./...
