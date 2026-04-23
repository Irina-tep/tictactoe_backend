package datasource

// Разработайте структуру хранения данных для отслеживания текущих игр.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// // структура хранения данных для отслеживания текущих игр.
type GameStorage struct {
	GameD *GameData
	// Timestamp time.Time //пока не используется
	CreatedAt time.Time
	UpdatedAt time.Time
}


type Storage struct {
	db *pgxpool.Pool /// Пул соединений (потокобезопасный)
}

// Конструктор с подключением к БД
func NewStorage(ctx context.Context, dbURL string) (*Storage, error) {
	// Парсим URL подключения
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Настройки пула соединений
	config.MaxConns = 10                     // максимальное количество соединений
	config.MinConns = 2                      // минимальное количество соединений
	config.MaxConnIdleTime = 5 * time.Minute // время жизни idle соединения
	config.MaxConnLifetime = 1 * time.Hour   // максимальное время жизни соединения

	// Создаем пул
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{
		db: pool}, nil
}

// Закрываем пул соединений
func (s *Storage) Close() {
	s.db.Close()
}

// DB возвращает пул соединений для выполнения SQL запросов
func (s *Storage) DB() *pgxpool.Pool {
	return s.db
}

func (s *Storage) SaveGame(ctx context.Context, gameID uuid.UUID, gameData *GameStorage) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if gameData == nil || gameData.GameD == nil {
		return fmt.Errorf("game data is nil")
	}
	// Преобразуем поле [3][3]int в JSON
	fieldJSON, err := json.Marshal(gameData.GameD.Field)
	if err != nil {
		return fmt.Errorf("failed to marshal field: %w", err)
	}
	now := time.Now()
	_, err = s.db.Exec(ctx, `INSERT INTO games (id, field, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id) DO UPDATE 
		SET field = $2, updated_at = $4`,
		gameID, fieldJSON, now, now)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}
	return nil
}

func (s *Storage) GetGame(ctx context.Context, gameID uuid.UUID) (*GameStorage, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var id uuid.UUID
	var fieldJSON []byte
	var createdAt, updatedAt time.Time
	err := s.db.QueryRow(ctx, "SELECT id, field, created_at, updated_at FROM games WHERE id = $1", gameID).Scan(&id, &fieldJSON, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	// Распаковываем JSON обратно в [3][3]int
	var field [3][3]int
	if err := json.Unmarshal(fieldJSON, &field); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	}
	game := &GameStorage{
		GameD: &GameData{
			ID:    id,
			Field: field,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return game, nil
}


