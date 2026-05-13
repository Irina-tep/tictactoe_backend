# Структура проекта "Крестики-Нолики" (Tic-Tac-Toe)

Проект реализует игру "Крестики-Нолики" с REST API и алгоритмом минимакса для AI игрока. Архитектура следует принципам чистой архитектуры (Clean Architecture) с разделением на слои.

## 📁 Общая структура

```
src/
├── api/                    # Слой контрактов API (HTTP handlers, DTOs, routing)
├── application/           # Бизнес-логика (доменный слой)
├── infrastructure/        # Инфраструктурный слой (хранилище, репозитории, маппинг)
├── di/                   # Слой внедрения зависимостей (Dependency Injection)
├── go.work              # Go workspace файл
└── go.work.sum          # Go workspace checksum
```

## 🏗️ Детальное описание модулей

### 1. **api/** - Слой контрактов API
Обработка HTTP запросов, валидация, преобразование данных.

| Файл | Назначение |
|------|------------|
| `router.go` | Маршрутизация HTTP запросов |
| `handler.go` | Обработчики HTTP запросов (CreateGame, GetGame, UpdateGame) |
| `dto.go` | Data Transfer Objects (DTO) для запросов/ответов |
| `mapper.go` | Преобразование между DTO и доменными моделями |
| `errors.go` | Ошибки уровня API |
| `handler_test.go` | Тесты API handlers |

**Эндпоинты:**
- `POST /game` - создать новую игру
- `GET /game/{id}` - получить состояние игры
- `POST /game/{id}` - сделать ход и получить ответ компьютера

### 2. **application/** - Бизнес-логика
Доменные модели и сервисы, реализующие игровую логику.

| Файл | Назначение |
|------|------------|
| `constants.go` | Константы игры (игроки, состояния, ошибки) |
| `current_game.go` | Модель текущей игры с UUID |
| `game_board.go` | Модель игрового поля с логикой проверок |
| `service_interface.go` | Интерфейс сервиса игры |
| `service_impl.go` | Реализация сервиса с алгоритмом минимакса |

**Ключевые компоненты:**
- `CurrentGame` - доменная модель игры с UUID
- `GameFild` - игровое поле 3×3 с методами проверки состояния
- `Service` - реализация алгоритма минимакса для AI

### 3. **infrastructure/** - Инфраструктурный слой
Реализация персистентности и доступа к данным.

| Файл | Назначение |
|------|------------|
| `save_data.go` | Потокобезопасное хранилище в памяти (sync.Map) |
| `repository_interface.go` | Интерфейс и реализация репозитория |
| `mapper.go` | Маппинг между доменом и хранилищем |
| `service.go` | Сервис инфраструктуры (координация репозитория и домена) |

**Ключевые компоненты:**
- `Storage` - потокобезопасное хранилище игр
- `GameRepository` - интерфейс репозитория
- `GameService` - сервис инфраструктуры

### 4. **di/** - Внедрение зависимостей
Настройка и инициализация компонентов приложения.

| Файл | Назначение |
|------|------------|
| `main.go` | Граф зависимостей с использованием uber/fx |

## 🔄 Поток данных (Data Flow)

```
HTTP Request → API Layer → Infrastructure Service → Application Service → Repository → Storage
```

1. **Клиент** отправляет HTTP запрос
2. **API Layer** валидирует и преобразует в доменные модели
3. **Infrastructure Service** координирует работу с данными
4. **Application Service** выполняет бизнес-логику (минимакс)
5. **Repository** обеспечивает доступ к данным
6. **Storage** хранит состояние в памяти

## 🧩 Зависимости между модулями

```
di → api → infrastructure → application
      ↑          ↑
      └──────────┘ (косвенная через интерфейсы)
```

- **di** зависит от всех модулей (собирает граф зависимостей)
- **api** зависит от **infrastructure** (использует GameService)
- **infrastructure** зависит от **application** (использует доменные модели)
- **application** не зависит от других модулей (чистая бизнес-логика)

## 🎯 Ключевые архитектурные принципы

### 1. **Инверсия зависимостей (DIP)**
- Интерфейсы определены в том же пакете, где используются
- Реализации зависят от интерфейсов, а не наоборот

### 2. **Разделение ответственности**
- **API Layer**: HTTP обработка, валидация, преобразование данных
- **Application Layer**: бизнес-правила, алгоритм минимакса
- **Infrastructure Layer**: персистентность, внешние зависимости
- **DI Layer**: композиция компонентов

### 3. **Тестируемость**
- Слои изолированы через интерфейсы
- Возможность мокирования зависимостей
- Интеграционные тесты для API

### 4. **Потокобезопасность**
- Хранилище использует `sync.Map`
- Каждая игра имеет уникальный UUID
- Независимые игровые сессии

## 🔧 Технологии и инструменты

- **Go 1.25.5** - основной язык разработки
- **uber/fx** - фреймворк для внедрения зависимостей
- **google/uuid** - генерация уникальных идентификаторов
- **Go Workspace** - управление мультимодульным проектом
- **net/http** - HTTP сервер и роутинг

## 🧪 Тестирование

Проект включает тесты API:
- Создание игры
- Получение состояния игры
- Последовательные ходы
- Обработка ошибок
- Конкурентные игры

**Запуск тестов:**
```bash
cd src/api
go test -v
```

## 🚀 Запуск проекта

### Через DI слой (рекомендуется)
```bash
cd src/di
go run main.go
```

Сервер запускается на `http://localhost:8080`

## 📊 Примеры использования API

### 1. Регистрация пользователя
Bash
Run
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "player1", "password": "pass123"}'
### 2. Авторизация (получение токена)
Bash
Run
curl -X POST http://localhost:8080/auth/login \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
В ответ получите user_id и token.

### 3. Создание новой игры (с авторизацией)
Bash
Run
curl -X POST http://localhost:8080/game \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"
В ответ получите id игры и пустое поле.

### 4. Выполнить ход (POST /game/{id})
Нужно отправить поле с одним изменением — поставить 1 (крестик, игрок X) в пустую клетку. Сервер сам сделает ответный ход (нолик, 2).

Bash
Run
curl -X POST http://localhost:8080/game/ВАШ_UUID_ИГРЫ \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"field": [[1,0,0],[0,0,0],[0,0,0]]}'

### 5. Получить состояние игры
Bash
Run
curl -X GET http://localhost:8080/game/ВАШ_UUID_ИГРЫ \
  -H "Authorization: Basic $(echo -n 'player1:pass123' | base64)"


## 🎮 Логика игры

- Игровое поле: матрица 3×3
- Значения: 0 (пусто), 1 (X), 2 (O)
- Первый ход: игрок X
- AI использует алгоритм минимакса
- Проверка победных комбинаций и ничьей

## 🔍 Детали реализации алгоритма

### Минимакс (Minimax)
- Рекурсивный алгоритм для поиска оптимального хода
- Оценка позиции: +10 за победу X, -10 за победу O, 0 за ничью
- Учет глубины для приоритизации быстрых побед
- Альфа-бета отсечение (при необходимости оптимизации)

### Определение текущего игрока
```go
// Если X и O равны → ход X, иначе ход O
if countX <= countO {
    return PlayerX
}
return PlayerO
```

## Убить процесс
sudo lsof -i :8080

kill -9 27417


## Игра с компьютером
# 1. Регистрация пользователя
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "irina", "password": "123456"}'

# 2. Аутентификация (получение токена)
curl -X POST http://localhost:8080/auth/login \
  -H "Authorization: Basic $(echo -n 'irina:123456' | base64)"

# 3. Создание игры
curl -X POST http://localhost:8080/game \
  -H "Authorization: Basic $(echo -n 'irina:123456' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"game_type": "pvc"}'

# 4. Ход игрока
curl -X POST http://localhost:8080/game/xxx \
  -H "Authorization: Basic $(echo -n 'irina:123456' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"field": [[0,0,0],[0,1,0],[0,0,0]]}'

## Игра с двумя игроками
# 1. Регистрация игрока А
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "irinaA", "password": "123456"}'

# 2. Регистрация игрока Б
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "irinaB", "password": "123456"}'

# 3. Аутентификация А
curl -X POST http://localhost:8080/auth/login \
  -H "Authorization: Basic $(echo -n 'irinaA:123456' | base64)"

# 4. Аутентификация Б
curl -X POST http://localhost:8080/auth/login \
  -H "Authorization: Basic $(echo -n 'irinaB:123456' | base64)"

# 5. А создаёт игру PvP
curl -X POST http://localhost:8080/game \
  -H "Authorization: Basic $(echo -n 'irinaA:123456' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"game_type": "pvp"}'
# → id = "xxx"

# 6. Б присоединяется
curl -X POST http://localhost:8080/game/xxx/join \
  -H "Authorization: Basic $(echo -n 'irinaB:123456' | base64)"

# 7. А ходит
curl -X POST http://localhost:8080/game/xxx \
  -H "Authorization: Basic $(echo -n 'irinaA:123456' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"field": [[0,0,0],[0,1,0],[0,0,0]]}'
# → id = "xxx"

# 8. Б ходит
curl -X POST http://localhost:8080/game/xxx \
  -H "Authorization: Basic $(echo -n 'irinaB:123456' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"field": [[2,0,0],[0,1,0],[0,0,0]]}'


# Список текущих игр 
curl -X GET http://localhost:8080/games \
  -H "Authorization: Basic $(echo -n 'kisa:123456' | base64)"

# Поиск игрока по ID
curl -X GET http://localhost:8080/info/5c2e9934-cd1c-46ef-bac2-57a09d2a29aa \
  -H "Authorization: Basic $(echo -n 'kisa:123456' | base64)"