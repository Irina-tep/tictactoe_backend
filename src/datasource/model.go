package datasource

import (
	"time"

	"github.com/google/uuid"
)

// GameData - структура хранения данных для отслеживания текущих игр.
type GameData struct {
	ID        uuid.UUID      `db:"id"`
	FieldJSON    string    	 `db:"field"`  //Field     [3][3]int
	GameState string         `db:"game_state"`
	GameType  string         `db:"game_type"`   //Players   [2]PlayersData
	PlayersJSON   string 		 `db:"players"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type PlayersData struct {
	PlayerID uuid.UUID
	Symbol   int
}
