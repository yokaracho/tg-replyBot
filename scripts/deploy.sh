#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Начинаем деплой telegram-reply-bot...${NC}"

DOCKER_IMAGE="telegram-reply-bot:latest"
COMPOSE_FILE="deployments/docker/docker-compose.yml"

if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker не установлен${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Docker Compose не установлен${NC}"
    exit 1
fi

if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo -e "${RED}TELEGRAM_BOT_TOKEN не установлен${NC}"
    exit 1
fi

echo -e "${YELLOW}Сборка Docker образа...${NC}"
docker build -t ${DOCKER_IMAGE} -f deployments/docker/Dockerfile .

echo -e "${YELLOW}Остановка старых контейнеров...${NC}"
docker-compose -f ${COMPOSE_FILE} down

echo -e "${YELLOW}Запуск новых контейнеров...${NC}"
docker-compose -f ${COMPOSE_FILE} up -d

echo -e "${YELLOW}Проверка статуса...${NC}"
sleep 5
docker-compose -f ${COMPOSE_FILE} ps

echo -e "${GREEN}Деплой завершен успешно!${NC}"
echo -e "${YELLOW}Для просмотра логов используйте: docker-compose -f ${COMPOSE_FILE} logs -f${NC}"