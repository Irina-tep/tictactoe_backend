package datasource

import (
	"github.com/google/uuid"
)

// GameData - модель для игрового поля текущей игры в слое источника данных
// Использует массив фиксированного размера [3][3]int
// Включает мьютекс для защиты от конкурентного доступа к конкретной игре
type GameData struct {
	ID        uuid.UUID
	Field     [3][3]int
	GameState string
	GameType string
	Players   [2]PlayersData
}

type PlayersData struct {
	PlayerID uuid.UUID
	Symbol   int
}

//модель с логином и паролем
// type SignUpRequest struct {
// 	Login string
// 	Password string
// }
