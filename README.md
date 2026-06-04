# ДЗ 10 Реализовать web сервер для управления данными через запросы

## Запуск
```bash
go run packages/launch/main.go
```
Сервер запускается на `localhost:8080`, консольный интерфейс слушает порт `8081`.

## REST API
Все CRUD эндпойнты доступны под `/api/v1/...`.

### Игра
| Метод | Путь | Описание |
| --- | --- | --- |
| `POST /api/v1/game` | Создать игру | `{"board_size":8,"player_one_name":"Alice","player_two_name":"Bob"}`
| `POST /api/v1/game/:id/move` | Сделать ход | `{"player":"Alice","position_from":"E2","position_to":"E4"}`
| `POST /api/v1/game/:id/automove` | Запланировать автоходы | `{"player":"Alice","moves_amount":2}`
| `GET /api/v1/game/:id` | Показать доску по id игры (через браузер)|
| `POST /api/v1/game/:id/surrender`, `POST /api/v1/game/:id/stop` | Сдаться / остановить игру


### CRUD для сущностей
Поддерживаются `player`, `figure`, `move`.
Каждый `POST` требует `Content-Type: application/json`.
| Сущность | POST | PUT | GET list | GET by id | DELETE |
| --- | --- | --- | --- | --- | --- |
| `player` | `/api/v1/player` — `name`, `figures_color`, optional `moves`, `figures_took` | `/api/v1/player/:id` | `/api/v1/players` | `/api/v1/player/:id` | `/api/v1/player/:id` |
| `figure` | `/api/v1/figure` — `symbol`, `piece_color`, `game_color` | `/api/v1/figure/:id` | `/api/v1/figures` | `/api/v1/figure/:id` | `/api/v1/figure/:id` |
| `move` | `/api/v1/move` — `position_from`, `position_to`, `figure`, `time_took` (ms) | `/api/v1/move/:id` | `/api/v1/moves` | `/api/v1/move/:id` | `/api/v1/move/:id` |

#### Примеры
```bash
curl -X POST localhost:8080/api/v1/player -H 'Content-Type: application/json' \
  -d '{"name":"Alice","figures_color":"white"}'

curl localhost:8080/api/v1/players
```

## Консольный клиент
Консольный клиент находится в директории console_client.
```bash
cd console_client
go run main.go
```