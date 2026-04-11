package datasource

import (
	"domain"
	"time"
)

// сопоставители между слоями предметной области и источника данных (предметная область <-> источник данных).
type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) ToDomain(data *GameStorage) (*domain.CurrentGame, error) {
	if data == nil || data.GameD == nil {
		return nil, ErrInvalidStorageData
	}
	field := domain.GameField{
		Field: data.GameD.Field,
	}
	id := data.GameD.ID
	return &domain.CurrentGame{
		ID:           id,
		CurrentField: &field,
	}, nil
}

func (m *Mapper) ToData(game *domain.CurrentGame) (*GameStorage, error) {
	if game == nil || game.CurrentField == nil {
		return nil, ErrInvalidGameData
	}
	return &GameStorage{
		GameD: &GameData{
			ID:    game.ID,
			Field: game.CurrentField.Field,
		},
		Timestamp: time.Now(),
	}, nil
}
