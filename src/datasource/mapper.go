package datasource

// import (
// 	"domain"
// 	"time"
// )

// // сопоставители между слоями предметной области и источника данных (предметная область <-> источник данных).
// type Mapper struct{}

// func NewMapper() *Mapper {
// 	return &Mapper{}
// }




// func (m *Mapper) ToDomain(data *GameStorage) (*domain.CurrentGame, error) {
// 	if data == nil || data.GameD == nil {
// 		return nil, ErrInvalidStorageData
// 	}
// 	field := domain.GameField{
// 		Field: data.GameD.Field,
// 	}
// 	id := data.GameD.ID
// 	var players [2]domain.Players
// 	for i := 0; i < 2; i++ {
// 		players[i] = domain.Players{
// 			PlayerID: data.GameD.Players[i].PlayerID,
// 			Symbol:   data.GameD.Players[i].Symbol,
// 		}
// 	}

// 	return &domain.CurrentGame{
// 		ID:           id,
// 		CurrentField: &field,
// 		GameState:    data.GameD.GameState,
// 		Players:      players,
// 		GameType: data.GameD.GameType,
// 	}, nil
// }

// func (m *Mapper) ToData(game *domain.CurrentGame) (*GameStorage, error) {
// 	if game == nil || game.CurrentField == nil {
// 		return nil, ErrInvalidGameData
// 	}
// 	var players [2]PlayersData
// 	for i := 0; i < 2; i++ {
// 		players[i] = PlayersData{
// 			PlayerID: game.Players[i].PlayerID,
// 			Symbol:   game.Players[i].Symbol,
// 		}
// 	}

// 	return &GameStorage{
// 		GameD: &GameData{
// 			ID:        game.ID,
// 			Field:     game.CurrentField.Field,
// 			GameState: game.GameState,
// 			Players:   players,
// 			GameType: game.GameType,
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}, nil
// }
