#!/usr/bin/env bash
# Создаёт локальное устройство и ключ прямо в БД локального auth-service и
# печатает JWT для приложения `dev` — SMTP и правки боевого кода не нужны.
#
# Требуется: bash, openssl, curl, docker.
# Запускать из каталога этого репозитория после ./up.sh.
set -euo pipefail

EMAIL="${LOCAL_EMAIL:-dev@tn.ru}"
APP_ID="dev"
[ -f .env ] && . ./.env
EMAIL="${LOCAL_EMAIL:-$EMAIL}"
SERVICE_SECRET="${LOCAL_SERVICE_SECRET:-local-dev-secret}"

COMPOSE="docker compose -f docker-compose.dev.yml"
[ -f .env ] && COMPOSE="$COMPOSE --env-file .env"

KEY="$(openssl rand -hex 32)"
if command -v uuidgen >/dev/null 2>&1; then
  DEVICE="$(uuidgen | tr '[:upper:]' '[:lower:]')"
else
  DEVICE="$(cat /proc/sys/kernel/random/uuid)"
fi
KEY_HASH="$(printf '%s' "$KEY" | openssl dgst -sha256 -hex | awk '{print $NF}')"

$COMPOSE exec -T auth-db psql -U auth -d auth -v ON_ERROR_STOP=1 <<SQL
INSERT INTO users (email) VALUES ('$EMAIL') ON CONFLICT DO NOTHING;
INSERT INTO devices (device_id, user_id, label)
  VALUES ('$DEVICE', '$EMAIL', 'dev-stand') ON CONFLICT (device_id) DO NOTHING;
INSERT INTO keys (key_hash, device_id, status)
  VALUES ('$KEY_HASH', '$DEVICE', 'active') ON CONFLICT (key_hash) DO NOTHING;
INSERT INTO apps (app_id, name, owner_email, service_secret)
  VALUES ('$APP_ID', 'Стенд разработки', '$EMAIL', '$SERVICE_SECRET')
  ON CONFLICT (app_id) DO NOTHING;
SQL

echo "Ключ устройства создан. Токен для приложения $APP_ID:" >&2
curl -s -X POST http://localhost:3000/auth/token \
  -H 'Content-Type: application/json' \
  -d "{\"key\":\"$KEY\",\"app_id\":\"$APP_ID\"}"
echo
