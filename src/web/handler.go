package web

// Приложение должно поддерживать одновременную игру в несколько игр .
import (
	"datasource"
	"domain"
	"encoding/json"
	"fmt"
	"net/http"
)

// обработчик с использованием net/http, используя следующий метод:
// POST /game/{current_game_UUID} — отправляет текущую игру с обновленным игровым полем пользователя и возвращает текущую игру с обновленным игровым полем компьютера.

type GameHandler struct {
	gameService *datasource.GameService
}

// Если отправлена ​​некорректная игра с неправильно обновленной доской, необходимо вернуть сообщение об ошибке с описанием .
// отправляет ошибку в формате JSON.
func (h *GameHandler) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json") //Устанавливает заголовок Content-Type: application/json, Браузер видит Content-Type и знает, что это JSON, это предотвращает неправильную обработку браузером
	w.WriteHeader(code)                                //Устанавливает HTTP статус-код ответа
	errResp := ErrorResponse{
		Error: message,
		Code:  code,
	}
	json.NewEncoder(w).Encode(errResp) //Кодирует структуру в JSON и сразу записывает в w
}

// создает новую игру. POST /game
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	game := NewGame()

	// ctx := r.Context()
	err := h.gameService.Repo.Save(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	response := CreateGameResponse{
		ID:    game.ID.String(),
		Field: game.CurrentField.Field,
	}

	//Отправка JSON ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)   //станавливает статус 201 Created
	json.NewEncoder(w).Encode(response) //Кодирует response в JSON и отправляет клиенту
}

// возвращает состояние игры.GET /game/{id}
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	// ctx := r.Context()
	game, err := h.gameService.Repo.GetByID(gameID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}

	response := h.gameToResponse(game)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// gameToResponse преобразует CurrentGame в GameResponse с заполнением всех полей.
func (h *GameHandler) gameToResponse(game *domain.CurrentGame) *GameResponse {
	if game == nil {
		return nil
	}

	scoreX, scoreO, finished := h.gameService.CheckFinish(*game)

	var status, nextTurn string
	if finished {
		if scoreX > scoreO {
			status = "x_won"
		} else if scoreO > scoreX {
			status = "o_won"
		} else {
			status = "draw"
		}
		nextTurn = ""
	} else {
		status = "in_progress"
		// Определяем чей следующий ход
		if h.gameService.GetCurrentPlayer(game.CurrentField) == domain.PlayerX {
			nextTurn = "X"
		} else {
			nextTurn = "O"
		}
	}

	return &GameResponse{
		ID:       game.ID.String(),
		Field:    game.CurrentField.Field,
		Status:   status,
		NextTurn: nextTurn,
		ScoreX:   scoreX,
		ScoreO:   scoreO,
		Finished: finished,
	}
}

// обрабатывает ход пользователя и выполняет ответный ход компьютера. POST /game/{id}
func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := r.PathValue("id")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	var req GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	currentGame, err := h.gameService.Repo.GetByID(gameID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}

	if currentGame.CurrentField.IsFinished() {
		http.Error(w, "Game is already finished", http.StatusConflict)
		return
	}

	currentPlayer := h.gameService.GetCurrentPlayer(currentGame.CurrentField)

	_, err = ToDomainFromRequest(gameID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid game field: %v", err), http.StatusBadRequest)
		return
	}

	// Проверяем разницу между текущим полем и присланным
	diffCount, changedRow, changedCol, oldVal, newVal := h.compareFieldsDetailed(currentGame.CurrentField.Field, req.Field)

	// Пользователь должен сделать ровно один ход
	if diffCount != 1 {
		http.Error(w, "Invalid board update: must change exactly one cell.", http.StatusBadRequest)
		return
	}

	// Проверяем, что изменение соответствует допустимому ходу (пустая клетка -> игрок, чей ход)
	if oldVal != domain.Empty {
		http.Error(w, fmt.Sprintf("Cell [%d][%d] is not empty", changedRow, changedCol), http.StatusBadRequest)
		return
	}
	if newVal != currentPlayer {
		http.Error(w, fmt.Sprintf("Cell [%d][%d] must be %d (current player), got %d", changedRow, changedCol, currentPlayer, newVal), http.StatusBadRequest)
		return
	}

	// Проверяем, что пользователь сделал ход за правильного игрока (это уже должно быть отражено в присланном поле).Обновляем игру присланным полем
	currentGame.CurrentField.Field = req.Field

	err = h.gameService.Repo.Save(currentGame)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	if currentGame.CurrentField.IsFinished() {
		response := h.gameToResponse(currentGame)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	row, col, err := h.gameService.GetBestMove(gameID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to compute computer move: %v", err), http.StatusInternalServerError)
		return
	}

	err = h.gameService.MakeMove(gameID, row, col)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make computer move: %v", err), http.StatusInternalServerError)
		return
	}

	updatedGame, err := h.gameService.Repo.GetByID(gameID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load updated game: %v", err), http.StatusInternalServerError)
		return
	}

	response := h.gameToResponse(updatedGame)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// сравнивает два поля и возвращает подробности первого изменения.
func (h *GameHandler) compareFieldsDetailed(field1, field2 [3][3]int) (int, int, int, int, int) {
	diffCount := 0
	var firstRow, firstCol int
	var firstOldVal, firstNewVal int
	foundFirst := false
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if field1[i][j] != field2[i][j] {
				diffCount++
				if !foundFirst {
					firstRow = i
					firstCol = j
					firstOldVal = field1[i][j]
					firstNewVal = field2[i][j]
					foundFirst = true
				}
			}
		}
	}
	return diffCount, firstRow, firstCol, firstOldVal, firstNewVal
}

// // определяет, чей сейчас ход.
// func (h *GameHandler) getCurrentPlayer(field *domain.GameField) int {
// 	countX := 0
// 	countO := 0
// 	for i := 0; i < 3; i++ {
// 		for j := 0; j < 3; j++ {
// 			if field.Fild[i][j] == application.PlayerX {
// 				countX++
// 			} else if field.Fild[i][j] == application.PlayerO {
// 				countO++
// 			}
// 		}
// 	}

// 	if countX <= countO {
// 		return application.PlayerX
// 	}
// 	return application.PlayerO
// }
