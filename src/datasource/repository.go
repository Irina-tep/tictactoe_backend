package datasource

import (
	"domain"
)

// интерфейс репозитория для работы с играми
type GameRepository interface {
	Save(game *domain.CurrentGame) error
	GetByID(gameID string) (*domain.CurrentGame, error)
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
	if game == nil {
		return ErrGameIsNil
	}
	gameData, err := r.mapper.ToData(game)
	if err != nil {
		return err
	}
	r.storage.SaveGame(game.ID.String(), gameData)
	return nil
}

func (r *GameRepositoryStruct) GetByID(gameID string) (*domain.CurrentGame, error) {
	if gameID == "" {
		return nil, ErrGameIDIsEmpty
	}
	gameData, ok := r.storage.GetGame(gameID)
	if !ok {
		return nil, ErrGameIDNotFound
	}
	return r.mapper.ToDomain(gameData)
}
