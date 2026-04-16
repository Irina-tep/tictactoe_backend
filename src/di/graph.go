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

	// Создаем таблицу, если её нет
	_, err = storage.DB().Exec(ctx, `CREATE TABLE IF NOT EXISTS games (
		id UUID PRIMARY KEY,
		field JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)`)
	if err != nil {
		panic(fmt.Sprintf("Failed to create table: %v", err))
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

// создает HTTP маршрутизатор.
func NewRouter(gameService *datasource.GameService) http.Handler {
	return web.NewRouter(gameService)
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
				fmt.Println("  POST   /game        - создать новую игру")
				fmt.Println("  GET    /game/{id}   - получить состояние игры")
				fmt.Println("  POST   /game/{id}   - сделать ход и получить ответ компьютера")

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
