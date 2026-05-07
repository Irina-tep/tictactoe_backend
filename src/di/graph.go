package main

// граф внедрения зависимостей .
// График должен включать как минимум следующие компоненты:
// Структура хранения , зарегистрированная как синглтон ;
// Репозиторий для работы со структурой хранения ;
// Сервис для работы с репозиторием.

import (
	"context"
	"datasource"
	"fmt"
	"log"
	"net/http"
	"time"
	"web"

	"go.uber.org/fx"
)

// создает экземпляр хранилища (синглтон).
func NewStorage(lc fx.Lifecycle) *datasource.Storage {
	ctx := context.Background()

	// Строка подключения к PostgreSQL
	// postgresql://ИмяПользователя:Пароль@Хост:Порт/ИмяБазыДанных
	dbURL := "postgresql://tictacuser:password@localhost:5432/tictactoe"

	storage, err := datasource.NewStorage(ctx, dbURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create storage: %v", err))
	}

	// Создаем таблицу games, если её нет
	_, err = storage.DB().Exec(ctx, `CREATE TABLE IF NOT EXISTS games (
		id UUID PRIMARY KEY,
		field JSONB NOT NULL,
		game_state TEXT DEFAULT 'waiting',
		players JSONB NOT NULL DEFAULT '[]'::jsonb,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		game_type TEXT DEFAULT 'pvp'
	)`)
	if err != nil {
		panic(fmt.Sprintf("Failed to create games table: %v", err))
	}

	// Создаем таблицу player, если её нет
	_, err = storage.DB().Exec(ctx, `CREATE TABLE IF NOT EXISTS player (
		id UUID PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(fmt.Sprintf("Failed to create player table: %v", err))
	}

	// Добавляем хук для закрытия хранилища при остановке приложения
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			storage.Close()
			return nil
		},
	})

	return storage
}

// создает репозиторий для работы с играми. Принимает хранилище, возвращает интерфейс GameRepository.
func NewGameRepo(storage *datasource.Storage) datasource.GameRepository {
	return datasource.NewGameRepositoryStruct(storage)
}

// создает сервис для работы с играми.
func NewGameService(repo datasource.GameRepository) *datasource.GameService {
	return datasource.NewGameService(repo)
}

// создает репозиторий для работы с пользователями. Принимает хранилище, возвращает интерфейс A.
func NewUserRepo(storage *datasource.Storage) datasource.UserRepository {
	return datasource.NewUserRepositoryStruct(storage)
}

// создает сервис для работы с пользователями.
func NewUserService(repo datasource.UserRepository) *datasource.AuthorizationService {
	return datasource.NewAutorizationService(repo)
}

// создает HTTP маршрутизатор.
func NewRouter(gameService *datasource.GameService, userService *datasource.AuthorizationService) http.Handler {
	return web.NewRouter(gameService, userService)
}

// NewHTTPServer создает и настраивает HTTP сервер.

func NewHTTPServer(lc fx.Lifecycle, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second, //Таймаут простаивающего соединения
	}

	//жизненный цикл сервера
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Println("Запуск API сервера для игры Крестики-Нолики")
				fmt.Println("Сервер запущен на http://localhost:8080")
				fmt.Println("Доступные эндпоинты:")
				fmt.Println("  POST   /auth/signup   - регистрация пользователя (без авторизации)")
				fmt.Println("  POST   /auth/login    - авторизация пользователя (без авторизации)")
				fmt.Println("  POST   /game          - создать новую игру (требуется авторизация)")
				fmt.Println("  GET    /game/{id}     - получить состояние игры (требуется авторизация)")
				fmt.Println("  GET    /games     - получить состояние игры")
				fmt.Println("  GET    /info     - получить состояние игры")
				fmt.Println("  POST   /game/{id}     - сделать ход (требуется авторизация)")
				fmt.Println("  POST   /game/{id}/join - прсоединиться к игре")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("Ошибка сервера: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Остановка сервера...")
			return server.Shutdown(ctx)
		},
	})
	return server
}
