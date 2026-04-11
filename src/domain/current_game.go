package domain

import (
	"github.com/google/uuid"
)

type CurrentGame struct {
	ID           uuid.UUID //UUID (Universally Unique Identifier) — это универсальный уникальный идентификатор.
	CurrentField *GameField
}

func NewCurrentGame(field *GameField) CurrentGame {
	currentGame := CurrentGame{
		ID:           uuid.New(),
		CurrentField: field,
	}
	return currentGame
}
