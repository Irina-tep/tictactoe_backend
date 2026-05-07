package datasource

// Разработайте структуру хранения данных для отслеживания текущих игр.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	// // Преобразуем поле [3][3]int в JSON
	// fieldJSON, err := json.Marshal(gameData.GameD.Field)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal field: %w", err)
	// }
	// // Преобразуем поле [2]Player в JSON
	// playersJSON, err := json.Marshal(gameData.GameD.Players)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal field: %w", err)
	// }
	now := time.Now()
	_, err := s.db.Exec(ctx, `INSERT INTO games (id, field, game_state, players, created_at, updated_at, game_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT (id) DO UPDATE 
		SET field = $2, game_state = $3, players = $4, updated_at = $6, game_type = $7`,
		gameID, gameData.GameD.Field, gameData.GameD.GameState, gameData.GameD.Players, now, now, gameData.GameD.GameType)
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
	var gameState, gameType string
	var playersJSON []byte
	var createdAt, updatedAt time.Time
	err := s.db.QueryRow(ctx, "SELECT id, field, game_state, players, created_at, updated_at, game_type FROM games WHERE id = $1", gameID).Scan(&id, &fieldJSON, &gameState, &playersJSON, &createdAt, &updatedAt, &gameType)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	// Распаковываем JSON обратно в [3][3]int
	var field [3][3]int
	if err := json.Unmarshal(fieldJSON, &field); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	}
	// Распаковываем JSON обратно в [2]uuid
	var players [2]PlayersData
	if err := json.Unmarshal(playersJSON, &players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	}
	game := &GameStorage{
		GameD: &GameData{
			ID:        id,
			Field:     field,
			GameState: gameState,
			Players:   players,
			GameType:  gameType,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return game, nil
}

func (s *Storage) GetGames(ctx context.Context) ([]*GameStorage, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// var id uuid.UUID
	// var field [3][3]int
	// var gameState, gameType string
	// var players [2]PlayersData
	// var createdAt, updatedAt time.Time
	// rows, err := s.db.Query(ctx, "SELECT id, field, game_state, players, created_at, updated_at, game_type FROM games WHERE game_state = waiting")
	rows, err := s.db.Query(ctx, "SELECT id, field, game_state, players, created_at, updated_at, game_type FROM games WHERE game_state NOT LIKE '%player_wins%' AND game_state != 'draw'")
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	defer rows.Close()
	users, err := pgx.CollectRows(rows, pgx.RowToStructByPos[GameStorage])
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	result := make([]*GameStorage, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	// // Распаковываем JSON обратно в [3][3]int
	// var field [3][3]int
	// if err := json.Unmarshal(fieldJSON, &field); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	// }
	// // Распаковываем JSON обратно в [2]uuid
	// var players [2]PlayersData
	// if err := json.Unmarshal(playersJSON, &players); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	// }
	// var games[len(users)-1]*GameStorage
	// for i,v := range games {
	// 	v = &GameStorage{
	// 	GameD: &GameData{
	// 		ID:        users[i]
	// 		id,
	// 		Field:     field,
	// 		GameState: gameState,
	// 		Players:   players,
	// 		GameType:  gameType,
	// 	},
	// 	CreatedAt: createdAt,
	// 	UpdatedAt: updatedAt,
	// }
	return result, nil
}
