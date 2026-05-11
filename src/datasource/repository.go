package datasource

import (
	"context"
	"domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// интерфейс репозитория для работы с играми
type GameRepository interface {
	Save(game *domain.CurrentGame) error
	GetByID(ctx context.Context, gameID uuid.UUID) (*domain.CurrentGame, error)
	GetByCurrentGames(ctx context.Context) ([]*domain.CurrentGame, error)
}

type GameStorage struct {
	db *pgxpool.Pool /// Пул соединений (потокобезопасный)
}

// Конструктор с подключением к БД
func NewStorage(ctx context.Context, dbURL string) (*GameStorage, error) {
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

	return &GameStorage{
		db: pool}, nil
}

// Закрываем пул соединений
func (s *GameStorage) Close() {
	s.db.Close()
}

// DB возвращает пул соединений для выполнения SQL запросов
func (s *GameStorage) DB() *pgxpool.Pool {
	return s.db
}

func (s *GameStorage) Save(game *domain.CurrentGame) error {
	ctx := context.Background()
	if game == nil {
		return ErrGameIsNil
	}
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	// Преобразуем поле [3][3]int в JSON
	fieldJSON, err := json.Marshal(game.CurrentField.Field)
	if err != nil {
		return fmt.Errorf("failed to marshal field: %w", err)
	}
	// Преобразуем поле [2]Player в JSON
	playersJSON, err := json.Marshal(game.Players)
	if err != nil {
		return fmt.Errorf("failed to marshal field: %w", err)
	}

	now := time.Now()
	_, err = s.db.Exec(ctx, `INSERT INTO games (id, field, game_state, players, created_at, updated_at, game_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT (id) DO UPDATE 
		SET field = $2, game_state = $3, players = $4, updated_at = $6, game_type = $7`,
		game.ID, fieldJSON, game.GameState, playersJSON, now, now, game.GameType)
	if err != nil {
		return ErrFailedToSaveGame
	}
	return nil
}

func (s *GameStorage) GetByID(ctx context.Context, gameID uuid.UUID) (*domain.CurrentGame, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	if gameID == uuid.Nil {
		return nil, ErrGameIDIsEmpty
	}
	// Получаем одну строку
	row, err := s.db.Query(ctx, "SELECT id, field, game_state, players, created_at, updated_at, game_type FROM games WHERE id = $1", gameID)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	game, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GameData])
	if err != nil {
		return nil, ErrGameIDNotFound
	}

	return s.rowToDomain(&game)
}

func (s *GameStorage) rowToDomain(row *GameData) (*domain.CurrentGame, error) {

	// Распаковываем JSON обратно в [3][3]int
	var field [3][3]int
	if err := json.Unmarshal([]byte(row.FieldJSON), &field); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	}
	gameField := &domain.GameField{
		Field: field,
	}
	// Распаковываем JSON обратно в [2]uuid
	var players [2]domain.Players
	if err := json.Unmarshal([]byte(row.PlayersJSON), &players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field: %w", err)
	}
	game := &domain.CurrentGame{
		ID:           row.ID,
		CurrentField: gameField,
		GameState:    row.GameState,
		Players:      players,
		GameType:     row.GameType,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return game, nil
}

func (s *GameStorage) GetByCurrentGames(ctx context.Context) ([]*domain.CurrentGame, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	//Получаем несколько строк
	rows, err := s.db.Query(ctx, "SELECT id, field, game_state, players, created_at, updated_at, game_type FROM games WHERE game_state NOT LIKE '%player_wins%' AND game_state != 'draw'")
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	defer rows.Close()
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[GameData])
	if err != nil {
		return nil, ErrGameIDNotFound
	}

	result := make([]*domain.CurrentGame, len(users))
	for i, v := range users {
		result[i], err = s.rowToDomain(&v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}
