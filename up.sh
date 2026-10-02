#!/usr/bin/env bash
# Поднимает локальный мини-стенд. Запускать из каталога этого репозитория.
set -euo pipefail

if [ ! -f .env ]; then
  echo "Нет .env — копирую из .env.example"
  cp .env.example .env
fi

docker compose -f docker-compose.dev.yml --env-file .env up -d --build

echo
echo "Стенд поднят:"
echo "  auth-service  http://localhost:3000/health"
echo "  app-service   http://localhost:8080/api/dev/health"
echo
echo "Токен для локального сервиса: ./bootstrap-local-token.sh"
