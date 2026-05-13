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

	id, err := uuid.Parse(gameID)
	if err != nil {
		return nil, ErrInvalidGameID
	}
	game := domain.CurrentGame{
		ID:           id,
		CurrentField: &field,
		GameState:    domain.PlayerToMove, 
	}
	return &game, nil
}

// создает новую игру c user из CreateGameRequest
func NewGameWithPlayer(userID uuid.UUID, gameType string) *domain.CurrentGame {
	field := domain.GameField{
		Field: [3][3]int{},
	}

	player1 := domain.Players{
		PlayerID: userID,
		Symbol:   domain.PlayerX,
	}
	var game domain.CurrentGame
	if gameType == "pvp" {
		// Против игрока — второй игрок пустой, статус waiting
		player2 := domain.Players{}
		game = domain.NewCurrentGame(&field, player1, player2)
		game.GameState = domain.WaitingForPlayers
	} else {
		// Против компьютера — второй игрок с нулевым UUID
		player2 := domain.Players{
			PlayerID: uuid.Nil,
			Symbol:   domain.PlayerO,
		}
		game = domain.NewCurrentGame(&field, player1, player2)
		game.GameState = domain.PlayerToMove + userID.String()

	}
	game.GameType = gameType
	return &game
}
