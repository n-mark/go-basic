# ДЗ 14: NoSQL хранилища (MongoDB + Redis)

## Описание
- **MongoDB** — основное хранилище доменных сущностей (`players`, `figures`, `moves`).
- **Redis** — журнал изменений (sorted set) с TTL:
  - ключ `history:<entity>` — глобальная лента по типу сущности;
  - ключ `history:<entity>:<id>` — история конкретной сущности (на этом ключе стоит TTL, который
    переустанавливается при каждом изменении).

Переключение между сохранением в json-файлы (предыдущий способ)
и MongoDB+Redis выполняется переменной окружения `STORAGE_TYPE`.

## Запуск инфраструктуры

```bash
docker compose up -d
docker compose ps   # дождаться статуса healthy
```

Будут подняты:
- `mongo` на порту `27017`;
- `redis` на порту `6379`.

## Переменные окружения

| Переменная          | Назначение                                 | По умолчанию                   |
| ------------------- | ------------------------------------------ | ------------------------------ |
| `STORAGE_TYPE`      | `local` (файлы) или `mongo` (Mongo+Redis)  | `mongo`                        |
| `MONGO_URI`         | URI подключения к MongoDB                  | `mongodb://localhost:27017`    |
| `MONGO_DB`          | Имя базы                                   | `go_basic`                     |
| `REDIS_ADDR`        | Адрес Redis                                | `localhost:6379`               |
| `REDIS_PASSWORD`    | Пароль Redis                               | (пусто)                        |
| `REDIS_ENTITY_TTL`  | TTL истории конкретной сущности            | `24h`                          |
| `REDIS_GLOBAL_TTL`  | TTL глобальной ленты                       | `168h` (7 дней)                |

## Как это работает

- `repository.MongoStorageProvider` реализует `StorageProvider`. Каждая сущность
  живёт в своей коллекции.
- `history.RedisLogger` пишет в sorted set с временной меткой в качестве score.
  TTL per-entity ключа обновляется при каждом изменении.
- `repository.Repo.recordHistory` — единая точка записи в журнал.

## Инструкция для тестового запуска

```bash
# 1. поднять инфраструктуру
docker compose up -d

# 2. запустить сервер
go run packages/launch/main.go

# 3. (опционально) запустить gRPC-клиент в ещё одном терминале
go run grpc_client/main.go

# 4. (опционально) посмотреть историю в Redis
docker exec -it go_basic_redis redis-cli ZRANGE history:player 0 -1
docker exec -it go_basic_redis redis-cli TTL history:player:1
```
