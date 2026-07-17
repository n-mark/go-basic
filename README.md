# ДЗ 15: Реляционные базы данных (PostgreSQL)

## Описание

- **PostgreSQL** — основное хранилище доменных сущностей (`players`, `figures`, `moves`).
- Миграции выполняются пакетом **goose** при старте сервиса из встроенного в результирующий артефакт через `embed.FS`.
- Запись идёт через `INSERT ... ON CONFLICT (id) DO UPDATE` (upsert) — идемпотентно, в единой транзакции.

Переключение между сохранением в json-файлы (предыдущий способ) и PostgreSQL выполняется
переменной окружения `STORAGE_TYPE`.

## Запуск инфраструктуры

```bash
docker compose up -d
docker compose ps    # дождаться статуса healthy
```

Будет поднят:

- `postgres` на порту `5432`.

## Переменные окружения

| Переменная        | Назначение                                          | По умолчанию              |
| ----------------- | --------------------------------------------------- | ------------------------- |
| `STORAGE_TYPE`    | `local` (файлы) или `pg` (PostgreSQL)               | `pg`                      |
| `PG_HOST`         | Хост PostgreSQL                                     | `localhost`               |
| `PG_PORT`         | Порт PostgreSQL                                     | `5432`                    |
| `PG_USER`         | Пользователь                                        | `gobasic`                 |
| `PG_PASSWORD`     | Пароль                                              | `gobasic`                 |
| `PG_DB`           | Имя базы                                            | `gobasic`                 |
| `PG_SSLMODE`      | Режим SSL (`disable` для локальной разработки)      | `disable`                 |

DSN собирается из этих переменных в `internal/config`.

## Как это работает

- `repository.PostgresStorageProvider` реализует `StorageProvider`. Каждая сущность
  хранится в своей таблице.
- `internal/migrations/00001_init.sql` описывает схему и поднимается
  через `goose.UpContext` при старте `CommonServer`.
- `INSERT … ON CONFLICT (id) DO UPDATE SET … = EXCLUDED.…` — upsert по `id`.
  Если запись уже есть — она обновляется, если нет — вставляется.

## Инструкция для тестового запуска

```bash
# 1. поднять инфраструктуру
docker compose up -d

# 2. запустить сервер (сам прогонит миграции и подключится)
go run packages/launch/main.go

# 3. (опционально) запустить gRPC-клиент в ещё одном терминале
go run grpc_client/main.go

# 4. (опционально) посмотреть данные в PostgreSQL
docker exec -it go_basic_postgres psql -U postgres -d go_basic -c "SELECT * FROM players;"
docker exec -it go_basic_postgres psql -U postgres -d go_basic -c "SELECT * FROM figures;"
docker exec -it go_basic_postgres psql -U postgres -d go_basic -c "SELECT * FROM moves;"
```


