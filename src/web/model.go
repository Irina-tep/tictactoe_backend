package web

import "github.com/google/uuid"

// запрос на создание новой игры
type CreateGameRequest struct {
}

// ответ с созданной игрой
type CreateGameResponse struct {
	ID    string    `json:"id"`
	Field [3][3]int `json:"field"`
}

// запрос на выполнение хода (альтернативный вариант)
type MoveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// запрос с обновленным игровым полем пользователя
type GameRequest struct {
	Field [3][3]int `json:"field"`
}

// ответ с состоянием игры
type GameResponse struct {
	ID       string    `json:"id"`
	Field    [3][3]int `json:"field"`
	Status   string    `json:"status"`
	NextTurn string    `json:"next_turn"`
	ScoreX   int       `json:"score_x"`
	ScoreO   int       `json:"score_o"`
	Finished bool      `json:"finished"`
}

// Запрос на регистрацию пользователя
type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Ответ на регистрацию
type SignUpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

// Ответ на аутентификацию
type AuthResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id"`
	Token   string `json:"token,omitempty"`
}

// UserAuthenticator — middleware для проверки авторизации.
// Содержит информацию об аутентифицированном пользователе.
type UserAuthenticator struct {
	UserID uuid.UUID
}

// Ошибка в формате JSON
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
