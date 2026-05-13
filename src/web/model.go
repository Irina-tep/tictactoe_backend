package web

import "time"

type CreateGameResponse struct {
	ID       string    `json:"id"`
	Field    [3][3]int `json:"field"`
	GameType string    `json:"game_type"`
}

type GameRequest struct {
	Field    [3][3]int     `json:"field"`
	GameType string        `json:"game_type"`
	Player1  PlayerRequest `json:"player1"`
	Player2  PlayerRequest `json:"player2"`
}

type PlayerRequest struct {
	ID     string `json:"id"`
	Symbol int    `json:"symbol"`
}

type GameResponse struct {
	ID       string    `json:"id"`
	Field    [3][3]int `json:"field"`
	Status   string    `json:"status"`
	NextTurn string    `json:"next_turn"`
	ScoreX   int       `json:"score_x"`
	ScoreO   int       `json:"score_o"`
	Finished bool      `json:"finished"`
	GameType string    `json:"game_type"`
}

type CurrentGamesResponse struct {
	ID        string           `json:"id"`
	Field     [3][3]int        `json:"field"`
	GameState string           `json:"game_state"`
	Players   [2]PlayerRequest `json:"players"`
	GameType  string           `json:"game_type"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id"`
	Token   string `json:"token,omitempty"`
}


type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}


type UserInfoResponse struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
	// Password string `json:"password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
