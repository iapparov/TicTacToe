# TicTacToe API

## Описание
API для игры "Крестики-нолики" с поддержкой регистрации, JWT-авторизации, создания и ведения игр, получения лидерборда и просмотра завершённых игр. Реализовано на Go с использованием PostgreSQL.


---

## Структура проекта

```
src/
├── cmd/
│   └── app/
│       └── main.go                # Точка входа в приложение
├── internal/
│   ├── app/
│   │   ├── game_service.go        # Бизнес-логика игры
│   │   ├── interface.go           # Интерфейсы сервисов
│   │   ├── jwt_models.go          # Модели для JWT
│   │   ├── jwt_provider.go        # Провайдер JWT-токенов
│   │   ├── model.go               # Модели приложения
│   │   └── user_service.go        # Бизнес-логика пользователей и JWT
│   ├── datasource/
│   │   ├── db_game.go             # Работа с таблицей игр в БД
│   │   ├── db_user.go             # Работа с таблицей пользователей в БД
│   │   ├── game_service.go        # Реализация сервисов для работы с играми
│   │   ├── interface.go           # Интерфейсы репозиториев
│   │   ├── mapper.go              # Маппинг моделей datasource <-> domain
│   │   ├── model.go               # Модели для БД
│   │   └── service.go             # Реализация сервисов для работы с БД
│   ├── di/
│   │   └── service.go             # DI-контейнер и запуск HTTP-сервера
│   └── web/
│       ├── game_handler.go        # HTTP-обработчики для игр
│       ├── mappers.go             # Маппинг моделей web <-> domain
│       ├── model.go               # Модели для web-слоя
│       ├── router.go              # Регистрация маршрутов
│       ├── user_handler.go        # HTTP-обработчики для пользователей и JWT
│       └── user_midlware.go       # JWT middleware для авторизации
```


---

## Основные возможности

- **Регистрация и авторизация пользователей (JWT)**
- **Создание новой игры (против игрока или компьютера)**
- **Подключение к игре**
- **Выполнение хода**
- **Получение списка всех игр пользователя**
- **Получение списка завершённых игр пользователя**
- **Получение лидерборда**
- **Валидация JWT для всех защищённых эндпоинтов**

---

## Основные компоненты

### 1. `cmd/app/main.go`
**Описание**: Точка входа в приложение. Настраивает зависимости и запускает HTTP-сервер.

**Ключевые функции**:
- `main()`: 
  - Настраивает DI-контейнер с помощью `fx.Provide`.
  - Регистрирует зависимости для работы с базой данных, бизнес-логики и веб-слоя.
  - Запускает HTTP-сервер через `fx.Invoke(di.StartHTTPServer)`.

---

### 2. `internal/app`

#### a. `game_service.go`
**Описание**: Реализация бизнес-логики для игры "Крестики-нолики".

**Ключевые функции**:
- `NewGame(Computer bool, Uuid string) (*CurrentGame)` — создание новой игры (против игрока или компьютера).
- `Connect(game *CurrentGame, Uuidgame string, Uuidplayero string) (*CurrentGame)` — подключение второго игрока.
- `FieldValidation(game *CurrentGame) (bool, error)` — проверка валидности игрового поля.
- `NextMove(game *CurrentGame) (*CurrentGame, error)` — выполнение следующего хода (с поддержкой minimax для компьютера).
- `GameIsOver(game *CurrentGame) bool` — проверка завершения игры.

#### b. `user_service.go`
**Описание**: Бизнес-логика пользователей и JWT.

**Ключевые функции**:
- Регистрация пользователя.
- Аутентификация и генерация JWT.
- Обновление access/refresh токенов.

#### c. `jwt_provider.go`, `jwt_models.go`
**Описание**: Генерация, валидация и работа с JWT-токенами.

---

### 3. `internal/datasource`

#### a. `db_game.go`
**Описание**: Работа с таблицей `games` в базе данных.

**Ключевые функции**:
- `SaveGame(currentgame *app.CurrentGame) error` — сохранение игры.
- `LoadGame(ID uuid.UUID) (*app.CurrentGame, error)` — загрузка игры по UUID.
- `GetGames() []string` — список всех открытых игр.
- `CurrentGame(Userid string) []string` — список игр пользователя.
- `GetEndedGames(uuid string) ([]string)` — список завершённых игр пользователя.
- `GetLeaderBoard(count int) ([]app.LeaderBoard, error)` — лидерборд по соотношению побед.

#### b. `db_user.go`
**Описание**: Работа с таблицей пользователей.

**Ключевые функции**:
- Сохранение и поиск пользователей по логину/UUID.

#### c. `game_service.go`
**Описание**: Адаптер между бизнес-логикой и репозиторием (реализация интерфейса GameService).

#### d. `mapper.go`, `model.go`
**Описание**: Маппинг моделей между слоями и описание сущностей для БД.

---

### 4. `internal/web`

#### a. `game_handler.go`
**Описание**: HTTP-обработчики для работы с играми.

**Ключевые функции**:
- `PlayGame(w http.ResponseWriter, r *http.Request)` — обработка хода игрока.
- `CreateGame(w http.ResponseWriter, r *http.Request)` — создание новой игры.
- `Connect(w http.ResponseWriter, r *http.Request)` — подключение к игре.
- `GetGames(w http.ResponseWriter, r *http.Request)` — список всех открытых игр.
- `CurrentGame(w http.ResponseWriter, r *http.Request)` — текущие игры пользователя.
- `GetEndedGames(w http.ResponseWriter, r *http.Request)` — завершённые игры пользователя.
- `GetLeaderBoard(w http.ResponseWriter, r *http.Request)` — лидерборд.
- `UserInfo(w http.ResponseWriter, r *http.Request)` — информация о пользователе.

#### b. `user_handler.go`
**Описание**: HTTP-обработчики для регистрации, логина и работы с JWT.

#### c. `user_midlware.go`
**Описание**: JWT middleware для авторизации всех защищённых эндпоинтов.

#### d. `router.go`
**Описание**: Регистрация всех маршрутов приложения.

#### e. `mappers.go`, `model.go`
**Описание**: Маппинг моделей web <-> domain, DTO для API.

---

### 5. `internal/di/service.go`
**Описание**: DI-контейнер и запуск HTTP-сервера через fx.

---

## Основные эндпоинты

### Аутентификация и пользователи

- `POST /register` — регистрация пользователя  
  **body:** `{ "login": "user", "password": "pass" }`

- `POST /login` — вход, получение JWT  
  **body:** `{ "login": "user", "password": "pass" }`

- `POST /refresh-access` — обновление access-токена  
  **body:** `{ "refresh_token": "<refresh_token>" }`

- `POST /refresh-refresh` — обновление refresh-токена  
  **body:** `{ "refresh_token": "<refresh_token>" }`

---

### Игры

- `POST /game` — создать новую игру  
  **body:** `{ "vs_computer": true|false }`  
  **header:** `Authorization: Bearer <access_token>`

- `POST /connect/{game_id}` — подключиться к игре  
  **header:** `Authorization: Bearer <access_token>`

- `POST /game/{game_id}` — сделать ход  
  **body:** `{ "id": "...", "field": [[...],[...],[...]] }`  
  **header:** `Authorization: Bearer <access_token>`

- `GET /games` — получить список всех игр пользователя  
  **header:** `Authorization: Bearer <access_token>`

- `GET /endedgames` — получить список завершённых игр пользователя  
  **header:** `Authorization: Bearer <access_token>`

- `GET /leaderboard?count=N` — получить топ-N игроков по соотношению побед  
  **header:** `Authorization: Bearer <access_token>`

---

## Пример работы API

### Регистрация и вход

```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"login": "user1", "password": "password123"}'

curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"login": "user1", "password": "password123"}'
```

### Создание новой игры

```bash
curl -X POST http://localhost:8080/game \
-H "Authorization: Bearer <access_token>" \
-d '{"vs_computer": true}'
```

### Получение завершённых игр

```bash
curl -X GET http://localhost:8080/endedgames \
-H "Authorization: Bearer <access_token>"
```

### Получение лидерборда

```bash
curl -X GET "http://localhost:8080/leaderboard?count=10" \
-H "Authorization: Bearer <access_token>"
```

---

## Настройка базы данных

Установите переменную окружения:

```bash
export DB_URL=postgres://postgres:password@localhost:5432/TicTacToe
```

---


## Настройка secretphrase для JWT

Установите переменную окружения:

```bash
export JWT_ACCESS_SECRET="YOUR_JWT_SECRET"
export JWT_REFRESH_SECRET="YOUR_JWT_SECRET"
```

---
## Запуск проекта

```bash
cd src
go mod tidy
go run cmd/app/main.go
```

API будет доступно по адресу http://localhost:8080

---

## Примечания

- Для всех защищённых эндпоинтов требуется JWT в заголовке `Authorization: Bearer <access_token>`.
- Все ответы и ошибки возвращаются в формате JSON.