# sbe-devstand-kit — локальный стенд разработки SBE-плагина

Комплект для разработки плагинов под систему SBE **вне реестра ЦУП**. Даёт:

- **чтение боевых** LLM, Базы знаний и почты — через центральный шлюз
  `dev-gateway` (только чтение);
- **локальную запись** — свой сервис и БД, поднятые в Docker.

Это L1 «мини-стенд»: локально живут только сервис разработчика, его БД и
локальный `auth-service`. Боевой контур не дублируется.

## Что нужно

- Docker и Docker Compose.
- `bash`, `openssl`, `curl` (для скриптов).
- Плагин `sbe-devstand`, установленный в Obsidian (см. его репозиторий).
- Доступ к стенду: роль в приложении «Стенд», выдаёт администратор.

## Раскладка

Склонируйте рядом два репозитория (комплект и ядро системы):

```
~/sbe-devstand-kit/     ← этот репозиторий
~/sbe-core/             ← https://github.com/Epyur/sbe-core
```

Другой путь к ядру — переменная `SBE_CORE_DIR` в `.env`.

## Быстрый старт

```bash
./up.sh                      # поднять локальный стенд
./bootstrap-local-token.sh   # создать ключ и напечатать JWT для сервиса
```

Затем в настройках плагина `sbe-devstand`:

- **Адрес локального сервиса** — `http://localhost:8080`;
- **Токен локального сервиса** — JWT из `bootstrap-local-token.sh`.

Проверка: `curl http://localhost:8080/api/dev/health` → `{"status":"ok"}`.

Боевое чтение (LLM/KB/почта) настраивать не нужно — оно идёт через
`https://epyur.fvds.ru/api/dev/*` под вашим токеном ЦУП.

## Из чего состоит

| Путь | Что |
|---|---|
| `docker-compose.dev.yml` | auth-db + auth-service + app-db + заготовка app-service |
| `app-service/` | шаблон Go-сервиса разработчика (JWT, своя таблица, запись) |
| `plugin-template/` | минимальный каркас плагина Obsidian |
| `bootstrap-local-token.sh` | создаёт устройство/ключ в локальной БД auth и печатает JWT |
| `up.sh` / `down.sh` | поднять / остановить |

## Почему так

Плагин вне реестра не может писать в боевые сервисы. Поэтому чтение идёт через
шлюз (там белый список только read-ручек), а запись — в локальный контейнер.
Локально используется приложение `dev`: `auth-service` уже умеет его
регистрировать, а `app-service` проверяет токен тем же секретом.

## Ключевые репозитории

Клонируйте по необходимости (для своего агента/чтения исходников). Доступ к
закрытым репозиториям выдаёт администратор.

| Репозиторий | Что это |
|---|---|
| `https://github.com/Epyur/sbe-core` | ядро клиента, `auth-service`, docker, `dev-gateway` |
| `https://github.com/Epyur/sbe-devstand` | плагин стенда (этот инструмент) |
| `https://github.com/Epyur/sbe-devstand-kit` | локальный комплект (этот репозиторий) |
| `https://github.com/Epyur/sbe-apstore` | ЦУП (магазин и авторизация) |
| `https://github.com/Epyur/sbe-apstore-registry` | записи реестра плагинов |
| `https://github.com/Epyur/sbe-llm` | центр LLM |
| `https://github.com/Epyur/sbe-kb` | база знаний (закрытый) |
| `https://github.com/Epyur/sbe-mailer` | письма |
| `https://github.com/Epyur/sbe-documents` | документы |
| `https://github.com/Epyur/sbe-photobank` | фотобанк |
| `https://github.com/Epyur/sbe-lims` | ЛИМС (там же lab-service) |
| `https://github.com/Epyur/sbe-requests` | заявки на испытания |
| `https://github.com/Epyur/sbe-ekn` | справочник ЕКН |
| `https://github.com/Epyur/sbe-contacts` | контакты |
| `https://github.com/Epyur/sbe-agent` | агент LogicTEAM.007 |
| `https://github.com/Epyur/sbe-presentations` | мастер презентаций |
| `https://github.com/Epyur/sbe-dashboards` | дашборды (приватный) |
| `https://github.com/Epyur/sbe-calendar` | календарь |
| `https://github.com/Epyur/sbe-tasks` | задачи |
| `https://github.com/Epyur/sbe-yougile` | интеграция YouGile |
| `https://github.com/Epyur/sbe-access` | консоль доступов |

## Ограничения

- Шлюз пускает только чтение — писать в боевые сервисы из стенда нельзя.
- Боевые KB/LLM/почта видны в объёме прав вашей учётной записи.
- Локальная копия центральных сервисов (KB, почта, LLM) не поднимается.
