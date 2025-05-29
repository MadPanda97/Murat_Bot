FROM golang:1.24 AS builder

WORKDIR /app

# Отключаем proxy для приватных/нестандартных путей
ENV GOPRIVATE=freelans/*
ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot
