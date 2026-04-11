package web

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

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}
