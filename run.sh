#!/bin/bash

echo "Установка зависимостей..."
go mod download

echo "Запуск сервера..."
go run cmd/server/main.go
