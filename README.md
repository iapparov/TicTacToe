# TicTacToe API

## Описание
Это проект API для игры "Крестики-нолики". API предоставляет функционал для создания игр, подключения игроков, выполнения ходов и проверки состояния игры. Реализовано на языке Go с использованием PostgreSQL для хранения данных.

---

## Структура проекта
```
src/
├── cmd/
│   └── app/
│       └── main.go // Точка входа в приложение
├── internal/
│   ├── app/
│   │   ├── interface.go // Интерфейсы для бизнес-логики
│   │   ├── model.go // Модели данных
│   │   ├── service.go // Реализация бизнес-логики
│   │   └── user_service.go // Логика работы с пользователями
│   ├── datasource/
│   │   ├── db_game.go // Работа с таблицей игр в базе данных
│   │   ├── db_user.go // Работа с таблицей пользователей
│   │   ├── interface.go // Интерфейсы для работы с базой данных
│   │   ├── mapper.go // Маппинг между сущностями базы данных и моделями
│   │   ├── model.go // Модели базы данных
│   │   └── service.go // Реализация сервисов для работы с базой данных
│   ├── di/
│   │   └── service.go // DI-контейнер для управления зависимостями
│   └── web/
│       ├── handler.go // HTTP-обработчики для работы с играми
│       ├── mappers.go // Маппинг между веб-моделями и доменными моделями
│       ├── model.go // Модели для веб-слоя
│       ├── router.go // Регистрация маршрутов
│       ├── user_handler.go // HTTP-обработчики для работы с пользователями
│       └── user_midlware.go // Middleware для авторизации пользователей
```

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

#### a. `service.go`
**Описание**: Реализация бизнес-логики для игры "Крестики-нолики".

**Ключевые функции**:
- `NewGame(Computer bool, Uuid string) (*CurrentGame)`:
  Создаёт новую игру. Если `Computer == true`, игра будет против компьютера.
  
- `Connect(game *CurrentGame, Uuidgame string, Uuidplayero string) (*CurrentGame)`:
  Подключает второго игрока к игре.

- `FieldValidation(game *CurrentGame) (bool, error)`:
  Проверяет валидность игрового поля.

- `NextMove(game *CurrentGame) (*CurrentGame, error)`:
  Выполняет следующий ход. Если игра против компьютера, используется алгоритм `minimax`.

- `GameIsOver(game *CurrentGame) bool`:
  Проверяет, завершена ли игра.

---

### 3. `internal/datasource`

#### a. `db_game.go`
**Описание**: Работа с таблицей `games` в базе данных.

**Ключевые функции**:
- `SaveGame(currentgame *app.CurrentGame) error`:
  Сохраняет игру в базе данных. Использует `ON CONFLICT` для обновления существующих записей.

- `LoadGame(ID uuid.UUID) (*app.CurrentGame, error)`:
  Загружает игру по её UUID.

- `CurrentGame(Userid string) []string`:
  Возвращает список игр, в которых участвует пользователь.

---

### 4. `internal/web`

#### a. `handler.go`
**Описание**: HTTP-обработчики для работы с играми.

**Ключевые функции**:
- `PlayGame(w http.ResponseWriter, r *http.Request)`:
  Обрабатывает ход игрока. Проверяет завершение игры, выполняет ход и возвращает обновлённое состояние игры.

- `CreateGame(w http.ResponseWriter, r *http.Request)`:
  Создаёт новую игру.

- `Connect(w http.ResponseWriter, r *http.Request)`:
  Подключает второго игрока к игре.

- `GetGames(w http.ResponseWriter, r *http.Request)`:
  Возвращает список доступных игр.

- `CurrentGame(w http.ResponseWriter, r *http.Request)`:
  Возвращает текущую игру для пользователя.

---

## Пример работы API

### 1. Создание новой игры
```bash
curl -X POST http://localhost:8080/game \
-H "Cookie: user_id=9e74ffe9-c435-48c0-b922-bd80d9e368bb" \
-d '{"vs_computer": true}'
```

### 2. Подключение второго игрока
```bash
curl -X POST http://localhost:8080/connect/8248de5d-6b1c-4208-83c3-c2049557b7e5 \
-H "Cookie: user_id=4929e317-7e1e-4799-86d2-8848149ba7f5"
```

### 3. Выполнение хода
```bash
curl -X POST http://localhost:8080/game/d25c5788-547d-492c-94ad-d4c06b07ddc5 \
-H "Cookie: user_id=9e74ffe9-c435-48c0-b922-bd80d9e368bb" \
-H "Content-Type: application/json" \
-d '{"id":"d25c5788-547d-492c-94ad-d4c06b07ddc5", "field":[[2,1,0],[2,2,1],[1,0,0]]}'
```

## Настройка базы данных
### Для подключения к базе данных используйте переменную окружения DB_URL. Пример:

```bash
export DB_URL=postgres://postgres:password@localhost:5432/TicTacToe
```

## Запуск проекта

### 1. Установите зависимости:
```
go mod tidy
```

### 2. Запустите приложение:
```bash
go run cmd/app/main.go
```

### 3. API будет доступно по адресу http://localhost:8080.
