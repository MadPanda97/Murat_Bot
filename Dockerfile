# Этап сборки
FROM golang:1.22 AS builder

WORKDIR /app

# Установка прокси для стабильной загрузки зависимостей
ENV GOPROXY=https://proxy.golang.org,direct

# Копируем go.mod и go.sum отдельно для кеширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Сборка бинарника
RUN go build -o bot

# Финальный минимальный образ
FROM debian:bullseye-slim

WORKDIR /app

# Копируем бинарник из стадии builder
COPY --from=builder /app/bot .

# Railway передаёт переменные окружения как ENV,
# но если ты используешь .env, можешь скопировать его
COPY .env .env

# Открываем порт (если используешь вебхуки, иначе можно опустить)
# EXPOSE 8080

# Запуск бота
CMD ["./bot"]
