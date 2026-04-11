package datasource

// Разработайте структуру хранения данных для отслеживания текущих игр.

import (
	"sync"
	"time"
)

// структура хранения данных для отслеживания текущих игр.
type GameStorage struct {
	GameD     *GameData
	Timestamp time.Time //пока не используется
}

// Для хранения данных используйте потокобезопасные коллекции (например, sync.Map?)
type Storage struct {
	m sync.Map // map[string]*GameData
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) SaveGame(gameID string, gameData *GameStorage) {
	s.m.Store(gameID, gameData)
}

func (s *Storage) GetGame(gameID string) (*GameStorage, bool) {
	val, ok := s.m.Load(gameID)
	if ok {
		return val.(*GameStorage), true
	}
	return nil, false
}
