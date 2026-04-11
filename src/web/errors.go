package web

import "errors"

// Ошибки уровня API
var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidFieldValue = errors.New("invalid field value: alowed value 0, 1, 2")
	ErrInvalidGameID     = errors.New("invalid game ID")
	ErrGameNotFound      = errors.New("game not found")
	ErrInvalidMove       = errors.New("invalid move")
	ErrGameFinished      = errors.New("game finished already")
)
