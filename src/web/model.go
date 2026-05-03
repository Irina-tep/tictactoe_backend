package web


// ответ с созданной игрой
type CreateGameResponse struct {
	ID       string    `json:"id"`
	Field    [3][3]int `json:"field"`
	GameType string    `json:"game_type"`
}

// запрос с обновленным игровым полем пользователя
type GameRequest struct {
	Field [3][3]int `json:"field"`
	GameType string `json:"game_type"`
	Player1 PlayerRequest `json:"player1"`
	Player2 PlayerRequest `json:"player2"`
}

type PlayerRequest struct {
	ID     string `json:"id"`
	Symbol int    `json:"symbol"`
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
	GameType string    `json:"game_type"`
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

// Ошибка в формате JSON
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
