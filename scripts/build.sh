#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Начинаем сборку telegram-reply-bot...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}Go не установлен${NC}"
    exit 1
fi

APP_NAME="telegram-reply-bot"
BINARY_NAME="bot"
BUILD_DIR="./bin"
MAIN_PATH="./cmd/bot"

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

echo -e "${YELLOW}Версия: ${VERSION}${NC}"
echo -e "${YELLOW}Время сборки: ${BUILD_TIME}${NC}"
echo -e "${YELLOW}Git commit: ${GIT_COMMIT}${NC}"

mkdir -p ${BUILD_DIR}

LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

platforms=("linux/amd64" "darwin/amd64" "windows/amd64")

for platform in "${platforms[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "${platform}"
    output_name=${BINARY_NAME}

    if [ ${GOOS} = "windows" ]; then
        output_name+='.exe'
    fi

    output_path=${BUILD_DIR}/${output_name}-${GOOS}-${GOARCH}
    if [ ${GOOS} = "windows" ]; then
        output_path+='.exe'
    fi

    echo -e "${YELLOW}Сборка для ${GOOS}/${GOARCH}...${NC}"
    
    env GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="${LDFLAGS}" \
        -o ${output_path} \
        ${MAIN_PATH}

    if [ $? -ne 0 ]; then
        echo -e "${RED}Ошибка сборки для ${GOOS}/${GOARCH}${NC}"
        exit 1
    fi
done

current_platform=$(go env GOOS)-$(go env GOARCH)
if [ -f "${BUILD_DIR}/${BINARY_NAME}-${current_platform}" ]; then
    ln -sf ${BINARY_NAME}-${current_platform} ${BUILD_DIR}/${BINARY_NAME}
fi

echo -e "${GREEN}Сборка завершена успешно!${NC}"
echo -e "${GREEN}Бинарные файлы находятся в директории: ${BUILD_DIR}${NC}"
ls -la ${BUILD_DIR}/