package web

// средства сопоставления данных между доменным и веб-уровнями (домен <-> веб)
import (
	"domain"

	"github.com/google/uuid"
)

func ToDomainFromRequest(gameID string, req *GameRequest) (*domain.CurrentGame, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			val := req.Field[i][j]
			if val != domain.Empty && val != domain.PlayerX && val != domain.PlayerO {
				return nil, ErrInvalidFieldValue
			}
		}
	}
	field := domain.GameField{
		Field: req.Field,
	}

	var game domain.CurrentGame
	if gameID == "" {
		game = domain.NewCurrentGame(&field)
	} else {

		id, err := uuid.Parse(gameID) //преобразования строкового представления UUID в его внутренний бинарный (или структурированный) формат
		if err != nil {
			return nil, ErrInvalidGameID
		}
		game = domain.CurrentGame{
			ID:           id,
			CurrentField: &field,
		}
	}
	return &game, nil
}

func ToResponse(game *domain.CurrentGame) *GameResponse {
	if game == nil {
		return nil
	}
	return &GameResponse{
		ID:    game.ID.String(),
		Field: game.CurrentField.Field,
		// Остальные поля заполняются в обработчике
	}
}

// создает новую игру из CreateGameRequest.
func NewGame() *domain.CurrentGame {
	field := domain.GameField{
		Field: [3][3]int{},
	}
	game := domain.NewCurrentGame(&field)
	return &game
}
