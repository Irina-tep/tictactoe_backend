# Структура проекта "Крестики-Нолики" (Tic-Tac-Toe)

Проект реализует игру "Крестики-Нолики" с REST API, поддержкой двух режимов (PvP и PvC), алгоритмом минимакса для AI игрока и аутентификацией пользователей. Данные хранятся в PostgreSQL.

## 📁 Общая структура

```
src/
├── web/                    # Слой HTTP (handlers, DTO, routing, middleware)
├── domain/                 # Бизнес-логика (доменные модели, минимакс)
├── datasource/             # Слой данных (PostgreSQL, репозитории, сервисы)
├── di/                     # Внедрение зависимостей (uber/fx)
├── datasource/database/    # SQL-скрипты
├── go.work                 # Go workspace файл
└── go.work.sum             # Go workspace checksum
```

## 🏗️ Детальное описание модулей

### 1. **web/** — Слой HTTP (API)
Обработка HTTP запросов, валидация, аутентификация, преобразование данных.

| Файл | Назначение |
|------|------------|
| `router.go` | Маршрутизация HTTP запросов и middleware аутентификации |
| `handler.go` | Обработчики HTTP запросов (SignUp, Authenticate, CreateGame, GetGame, UpdateGame, JoinGame, GetGames, GetInfo) |
| `model.go` | Data Transfer Objects (DTO) для запросов/ответов |
| `mapper.go` | Преобразование между DTO и доменными моделями |
| `errors.go` | Ошибки уровня API |

**Эндпоинты:**

| Метод | Путь | Описание | Аутентификация |
|-------|------|----------|:--------------:|
| POST | `/auth/signup` | Регистрация пользователя | ❌ |
| POST | `/auth/login` | Авторизация (Basic Auth) | ❌ |
| POST | `/game` | Создать новую игру | ✅ |
| GET | `/game/{id}` | Получить состояние игры | ✅ |
| POST | `/game/{id}` | Сделать ход | ✅ |
| POST | `/game/{id}/join` | Присоединиться к PvP-игре | ✅ |
| GET | `/games` | Список активных игр | ✅ |
| GET | `/info/{id}` | Информация об игроке | ✅ |

### 2. **domain/** — Бизнес-логика (доменный слой)
Доменные модели и сервисы, реализующие игровую логику и алгоритм минимакса.

| Файл | Назначение |
|------|------------|
| `constants.go` | Константы игры (игроки, состояния, типы игр, ошибки) |
| `current_game.go` | Модель текущей игры с UUID, игроками и состоянием |
| `game_board.go` | Модель игрового поля 3×3 с логикой проверок и минимаксом |
| `user.go` | Модель пользователя с валидацией логина/пароля |
| `service_interface.go` | Интерфейсы сервисов (игра, пользователи) |
| `service_implementation.go` | Реализация сервиса с алгоритмом минимакса |

**Ключевые компоненты:**
- `CurrentGame` — доменная модель игры с UUID, полем, игроками, состоянием и типом
- `GameField` — игровое поле 3×3 с методами `MakeMove`, `CheckResult`, `IsFinished`, `Evaluate`
- `User` — модель пользователя с валидацией
- `GameServiceInterface` — реализация алгоритма минимакса и определения текущего игрока

### 3. **datasource/** — Слой данных (инфраструктура)
Реализация персистентности в PostgreSQL, репозитории и сервисы доступа к данным.

| Файл | Назначение |
|------|------------|
| `model.go` | Структуры данных для хранения (GameData, UserData) |
| `repository.go` | Реализация GameRepository (сохранение/загрузка игр через pgxpool) |
| `user_repository.go` | Реализация UserRepository (сохранение/загрузка пользователей) |
| `service.go` | Сервис координации (GameService) — бизнес-логика ходов, проверки |
| `user_service.go` | Сервис аутентификации (AuthorizationService) — регистрация, Basic Auth |
| `errors.go` | Ошибки уровня данных |
| `database/game.sql` | SQL-скрипты для создания таблиц и отладки |

**Ключевые компоненты:**
- `GameStorage` — пул соединений PostgreSQL (pgxpool), реализует `GameRepository`
- `UserStorage` — обёртка над GameStorage для работы с пользователями, реализует `UserRepository`
- `GameService` — сервис, координирующий репозиторий и доменную логику (ходы, проверки, минимакс)
- `AuthorizationService` — сервис регистрации и аутентификации через Basic Auth

### 4. **di/** — Внедрение зависимостей
Настройка и инициализация компонентов приложения с помощью uber/fx.

| Файл | Назначение |
|------|------------|
| `graph.go` | Граф зависимостей: создание Storage, репозиториев, сервисов, HTTP-сервера |
| `main.go` | Точка входа, запуск fx-приложения |

## 🔄 Поток данных (Data Flow)

```
HTTP Request → web (handler) → datasource (GameService) → domain (минимакс) → datasource (GameRepository) → PostgreSQL
```

1. **Клиент** отправляет HTTP запрос (с Basic Auth для защищённых эндпоинтов)
2. **web/handler** проверяет аутентификацию через middleware, валидирует запрос
3. **datasource/GameService** загружает игру из БД, вызывает доменную логику
4. **domain** выполняет бизнес-логику (минимакс, проверка поля, определение хода)
5. **datasource/GameRepository** сохраняет обновлённое состояние в PostgreSQL

## 🧩 Зависимости между модулями

```
di → web → datasource → domain
```

- **di** зависит от всех модулей (собирает граф зависимостей через uber/fx)
- **web** зависит от **datasource** (использует GameService и AuthorizationService)
- **datasource** зависит от **domain** (использует доменные модели и интерфейсы)
- **domain** не зависит от других модулей (чистая бизнес-логика)

## 🎯 Ключевые архитектурные принципы

### 1. **Инверсия зависимостей (DIP)**
- Интерфейсы (`GameRepository`, `UserRepository`) определены в **datasource**
- Реализации зависят от интерфейсов, а не наоборот

### 2. **Разделение ответственности**
- **web**: HTTP обработка, валидация, аутентификация, DTO
- **domain**: бизнес-правила, алгоритм минимакса, модели
- **datasource**: персистентность (PostgreSQL), репозитории, сервисы координации
- **di**: композиция компонентов, жизненный цикл приложения

### 3. **Два режима игры**
- **PvC** (Player vs Computer) — игрок против AI (минимакс)
- **PvP** (Player vs Player) — два игрока через join

### 4. **Аутентификация**
- Basic Auth (base64(login:password))
- Middleware проверяет авторизацию для защищённых эндпоинтов
- UserID передаётся через контекст запроса

## 🔧 Технологии и инструменты

- **Go 1.25.5** — основной язык разработки
- **uber/fx** — фреймворк для внедрения зависимостей
- **google/uuid** — генерация уникальных идентификаторов
- **jackc/pgx/v5** — драйвер PostgreSQL и пул соединений (pgxpool)
- **Go Workspace** — управление мультимодульным проектом
- **net/http** — HTTP сервер и роутинг (Go 1.22+ PathValue)

## 🚀 Запуск проекта

### Требования
- Go 1.25.5+
- PostgreSQL (например, через Docker)

### Настройка БД
```bash
# Создание пользователя и базы данных (пример)
psql -U postgres -c "CREATE USER tictacuser WITH PASSWORD 'password';"
psql -U postgres -c "CREATE DATABASE tictactoe OWNER tictacuser;"
```

### Запуск через DI слой (рекомендуется)
```bash
cd src/di
go run main.go
```

Сервер запускается на `http://localhost:8080`

Таблицы `games` и `player` создаются автоматически при запуске.

## 📊 Примеры использования API

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "player1", "password": "pass123"}'
```

### 2. Авторизация (получение токена)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
```
В ответ получите `user_id` и `token`.

### 3. Создание новой игры (PvC — с компьютером)
```bash
curl -X POST http://localhost:8080/game \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"game_type": "pvc"}'
```

### 4. Создание новой игры (PvP — с другим игроком)
```bash
curl -X POST http://localhost:8080/game \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"game_type": "pvp"}'
```

### 5. Выполнить ход (POST /game/{id})
Нужно отправить поле с одним изменением — поставить 1 (крестик, игрок X) в пустую клетку. Для PvC сервер сам сделает ответный ход (нолик, 2).

```bash
curl -X POST http://localhost:8080/game/ВАШ_UUID_ИГРЫ \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"field": [[1,0,0],[0,0,0],[0,0,0]]}'
```

### 6. Получить состояние игры
```bash
curl -X GET http://localhost:8080/game/ВАШ_UUID_ИГРЫ \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
```

### 7. Присоединиться к PvP-игре
```bash
curl -X POST http://localhost:8080/game/ВАШ_UUID_ИГРЫ/join \
  -H "Authorization: Basic $(echo -n 'player2:pass123' | base64)"
```

### 8. Список активных игр
```bash
curl -X GET http://localhost:8080/games \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
```

### 9. Информация об игроке
```bash
curl -X GET http://localhost:8080/info/ID_ИГРОКА \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
```

## 🎮 Логика игры

- Игровое поле: матрица 3×3
- Значения: `0` (пусто), `1` (X), `2` (O)
- Первый ход: игрок X
- AI использует алгоритм минимакса
- Проверка победных комбинаций и ничьей

### Состояния игры
| Константа | Описание |
|-----------|----------|
| `waiting` | Ожидание второго игрока (PvP) |
| `player_to_move{UUID}` | Ход игрока с указанным UUID |
| `player_wins{UUID}` | Победа игрока с указанным UUID |
| `draw` | Ничья |

## 🔍 Детали реализации алгоритма

### Минимакс (Minimax)
- Рекурсивный алгоритм для поиска оптимального хода
- Оценка позиции: +10 за победу X, -10 за победу O, 0 за ничью
- Учет глубины для приоритизации быстрых побед

### Определение текущего игрока
```go
// Если X и O равны → ход X, иначе ход O
if countX <= countO {
    return PlayerX
}
return PlayerO
```

## 🗄️ Структура БД

### Таблица `games`
| Колонка | Тип | Описание |
|---------|-----|----------|
| id | UUID PRIMARY KEY | Уникальный идентификатор игры |
| field | JSONB | Игровое поле 3×3 |
| game_state | TEXT | Состояние игры |
| players | JSONB | Массив игроков [id, symbol] |
| game_type | TEXT | Тип игры: "pvp" или "pvc" |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |

### Таблица `player`
| Колонка | Тип | Описание |
|---------|-----|----------|
| id | UUID PRIMARY KEY | Уникальный идентификатор |
| login | VARCHAR(255) UNIQUE | Логин пользователя |
| password | VARCHAR(255) | Пароль |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |

## Убить процесс
```bash
sudo lsof -i :8080
kill -9 PID
```