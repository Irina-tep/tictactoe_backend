package datasource

import (
	"context"
	"domain"

	"github.com/google/uuid"
)

// интерфейс репозитория для работы с играми
type GameRepository interface {
	Save(game *domain.CurrentGame) error
	GetByID(ctx context.Context, gameID uuid.UUID) (*domain.CurrentGame, error)
}

type GameRepositoryStruct struct {
	storage *Storage
	mapper  *Mapper
}

func NewGameRepositoryStruct(s *Storage) *GameRepositoryStruct {
	return &GameRepositoryStruct{
		storage: s,
		mapper:  NewMapper(),
	}
}

func (r *GameRepositoryStruct) Save(game *domain.CurrentGame) error {
	ctx := context.Background()
	if game == nil {
		return ErrGameIsNil
	}
	gameData, err := r.mapper.ToData(game)
	if err != nil {
		return err
	}
	r.storage.SaveGame(ctx, game.ID, gameData)
	return nil
}

func (r *GameRepositoryStruct) GetByID(ctx context.Context, gameID uuid.UUID) (*domain.CurrentGame, error) {
	if gameID == uuid.Nil {
		return nil, ErrGameIDIsEmpty
	}
	gameData, err := r.storage.GetGame(ctx, gameID)
	if err != nil {
		return nil, ErrGameIDNotFound
	}
	return r.mapper.ToDomain(gameData)
}
