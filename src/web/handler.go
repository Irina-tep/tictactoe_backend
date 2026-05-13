package web


import (
	"context"
	"datasource"
	"domain"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Ключ для хранения UserID в контексте запроса
type contextKey string

const userIDKey contextKey = "userID"

type GameHandler struct {
	gameService *datasource.GameService
	userService *datasource.AuthorizationService
}

func (h *GameHandler) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json") 
	w.WriteHeader(code)                                
	errResp := ErrorResponse{
		Error: message,
		Code:  code,
	}
	json.NewEncoder(w).Encode(errResp) 
}

// UserAuthenticator создает middleware, который проверяет авторизацию пользователя.
// Возвращает обёрнутый http.Handler, который проверяет заголовок Authorization
// и добавляет UserID в контекст запроса.
func (h *GameHandler) UserAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.sendError(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		userID, err := h.userService.Authenticate(authHeader)
		if err != nil {
			h.sendError(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}


func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	userID, ok := GetUserIDFromContext(ctx)

	if !ok {
		h.sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	var req GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.GameType = "pvc"
	}
	if req.GameType != "pvc" && req.GameType != "pvp" {
		req.GameType = "pvc"
	}

	game := NewGameWithPlayer(userID, req.GameType)

	err := h.gameService.Repo.Save(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	response := CreateGameResponse{
		ID:       game.ID.String(),
		Field:    game.CurrentField.Field,
		GameType: req.GameType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)   
	json.NewEncoder(w).Encode(response) 
}


func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(gameID)
	if err != nil {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	game, err := h.gameService.Repo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}

	response := h.gameToResponse(game)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


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

	ctx := r.Context()
	userID, success := GetUserIDFromContext(ctx)
	if !success {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(gameID)
	if err != nil {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	// Загружаем игру, чтобы узнать её тип
	currentGame, err := h.gameService.Repo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}
	if currentGame.CurrentField.IsFinished() {
		http.Error(w, "Game is already finished", http.StatusConflict)
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
	// Определяем, какой игрок должен ходить
	currentPlayer := h.gameService.GetCurrentPlayer(currentGame.CurrentField)

	if newVal != currentPlayer {
		http.Error(w, fmt.Sprintf("Cell [%d][%d] must be %d (current player), got %d", changedRow, changedCol, currentPlayer, newVal), http.StatusBadRequest)
		return
	}

	// Проверяем, что ходит именно тот пользователь, который должен
	var expectedPlayerID uuid.UUID
	for _, p := range currentGame.Players {
		if p.Symbol == currentPlayer {
			expectedPlayerID = p.PlayerID
			break
		}
	}
	if userID != expectedPlayerID {
		h.sendError(w, "Not your turn", http.StatusForbidden)
		return
	}
	
	err = h.gameService.MakeMove(ctx, id, userID, changedRow, changedCol)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make computer move: %v", err), http.StatusInternalServerError)
		return
	}

	if currentGame.GameType == "pvc" {
		//для игры с компьютером
		// Загружаем игру снова, чтобы проверить, не закончилась ли она после хода пользователя
		updatedGame, err := h.gameService.Repo.GetByID(ctx, id)
		if err != nil {
			http.Error(w, "Failed to load game", http.StatusInternalServerError)
			return
		}
		if !updatedGame.CurrentField.IsFinished() {
			row, col, err := h.gameService.GetBestMove(ctx, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to compute computer move: %v", err), http.StatusInternalServerError)
				return
			}

			err = h.gameService.MakeMove(ctx, id, uuid.Nil, row, col)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to make computer move: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	// Возвращаем обновлённую игру
	resultGame, err := h.gameService.Repo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Failed to load game", http.StatusInternalServerError)
		return
	}

	response := h.gameToResponse(resultGame)
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


func (h *GameHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Валидация логина и пароля
	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)

	if req.Login == "" || req.Password == "" {
		h.sendError(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	
	success, err := h.userService.Registration(req.Login, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		
		if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "password") {
			status = http.StatusUnauthorized
		}
		h.sendError(w, err.Error(), status)
		return
	}
	response := SignUpResponse{
		Success: success,
		Message: "User registered successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func (h *GameHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendError(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}
	
	userID, err := h.userService.Authenticate(authHeader)
	if err != nil {
		h.sendError(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	response := AuthResponse{
		Success: true,
		UserID:  userID.String(),
		Token:   authHeader,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	gameID := r.PathValue("id")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		h.sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(gameID)
	if err != nil {
		h.sendError(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	// Получаем игр
	game, err := h.gameService.Repo.GetByID(ctx, id)
	if err != nil {
		h.sendError(w, "Game not found", http.StatusNotFound)
		return
	}

	// Проверяем, что игра ожидает игроков
	if game.GameState != domain.WaitingForPlayers {
		h.sendError(w, "Game is not accepting players", http.StatusConflict)
		return
	}
	// Проверяем, что это PvP игра
	if game.GameType != domain.GameTypePvP {
		h.sendError(w, "Cannot join a computer game", http.StatusBadRequest)
		return
	}
	// Проверяем, что пользователь еще не в игре
	for _, p := range game.Players {
		if p.PlayerID == userID {
			h.sendError(w, "You are already in this game", http.StatusConflict)
			return
		}
	}
	for i, p := range game.Players {
		if p.PlayerID == uuid.Nil {
			// Определяем символ для нового игрока
			var symbol int
			if i == 0 {
				symbol = domain.PlayerX
			} else {
				symbol = domain.PlayerO
			}
			game.Players[i] = domain.Players{
				PlayerID: userID,
				Symbol:   symbol,
			}
			// Меняем статус игры
			game.GameState = domain.PlayerToMove + game.Players[0].PlayerID.String()
			// Сохраняем изменения
			err = h.gameService.Repo.Save(game)
			if err != nil {
				h.sendError(w, "Failed to join game", http.StatusInternalServerError)
				return
			}
			response := h.gameToResponse(game)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

	}
	// Если все слоты заняты
	h.sendError(w, "Game is full", http.StatusConflict)
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	_, ok := GetUserIDFromContext(ctx)
	if !ok {
		h.sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	game, err := h.gameService.Repo.GetByCurrentGames(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Game not found: %v", err), http.StatusNotFound)
		return
	}
	response := make([]CurrentGamesResponse, len(game))
	for i, v := range game {
		players := [2]PlayerRequest{}
		players[0] = PlayerRequest{
			ID:     v.Players[0].PlayerID.String(),
			Symbol: v.Players[0].Symbol,
		}
		players[1] = PlayerRequest{
			ID:     v.Players[1].PlayerID.String(),
			Symbol: v.Players[1].Symbol,
		}

		response[i] = CurrentGamesResponse{
			ID:        v.ID.String(),
			Field:     v.CurrentField.Field,
			GameState: v.GameState,
			GameType:  v.GameType,
			Players:   players,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func (h *GameHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	_, ok := GetUserIDFromContext(ctx)
	if !ok {
		h.sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	user, err := h.userService.UserRepo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("ID not found: %v", err), http.StatusNotFound)
		return
	}

	response := &UserInfoResponse{
		UserID:     user.ID.String(),
		Login:      user.Login,
		// Password:   user.Password,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
