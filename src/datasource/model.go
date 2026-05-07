package datasource

import (
	"github.com/google/uuid"
)

// GameData - модель для игрового поля текущей игры в слое источника данных
// Использует массив фиксированного размера [3][3]int
// Включает мьютекс для защиты от конкурентного доступа к конкретной игре
type GameData struct {
	ID        uuid.UUID
	Field     [3][3]int `json:"field"`
	GameState string
	GameType string
	Players   [2]PlayersData  `json:"players"`
}

type PlayersData struct {
	PlayerID uuid.UUID
	Symbol   int
}

