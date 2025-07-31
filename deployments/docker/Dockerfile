FROM golang:1.21-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot ./cmd/bot

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/bot .

# Копируем конфигурации
COPY --from=builder /app/configs ./configs

# Создаем пользователя без привилегий
RUN adduser -D -s /bin/sh botuser
USER botuser

# Порт (если нужен для webhook)
EXPOSE 8080

CMD ["./bot"]