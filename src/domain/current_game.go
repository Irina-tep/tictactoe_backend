package domain

import (
	"time"

	"github.com/google/uuid"
)

type CurrentGame struct {
	ID           uuid.UUID //UUID (Universally Unique Identifier) — это универсальный уникальный идентификатор.
	CurrentField *GameField
	GameState    string     // состояния для текущей игры
	GameType string
	Players      [2]Players //id, Symbol обоих игроков
	CreatedAt time.Time  //необязательное поле
	UpdatedAt time.Time //необязательное поле
}

type Players struct {
	PlayerID uuid.UUID
	Symbol   int
}

func NewCurrentGame(field *GameField, player1, player2 Players) CurrentGame {
	currentGame := CurrentGame{
		ID:           uuid.New(),
		CurrentField: field,
		GameState:    WaitingForPlayers, 
		Players:      [2]Players{player1, player2},
	}
	return currentGame
}
