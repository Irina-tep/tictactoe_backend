package datasource

import "errors"

var (
	ErrInvalidStorageData = errors.New("invalid storage data")
	ErrInvalidGameData    = errors.New("invalid game data")
	ErrGameIsNil          = errors.New("game cannot be nil")
	ErrGameIDIsEmpty      = errors.New("game ID cannot be empty")
	ErrGameIDNotFound     = errors.New("game with ID not found")
)
