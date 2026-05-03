// Создайте структуру, которая реализует интерфейс сервиса и принимает в качестве параметра интерфейс репозитория для работы со структурой хранения.
package datasource

import (
	"context"
	"domain"
	"errors"

	"github.com/google/uuid"
)

type GameService struct {
	Repo       GameRepository              // репозиторий для доступа к данным
	appService domain.GameServiceInterface // сервис предметной области для реализации логики
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{
		Repo:       repo,
		appService: domain.NewService(),
	}
}

func (gs *GameService) NextTurn(currentGame domain.CurrentGame) (int, int, error) {
	return gs.appService.NextTurn(currentGame)
}

func (gs *GameService) CheckField(currentGame domain.CurrentGame, expectedGame domain.CurrentGame) (bool, error) {
	return gs.appService.CheckField(currentGame, expectedGame)
}

func (gs *GameService) CheckFinish(currentGame domain.CurrentGame) (int, int, bool) {
	return gs.appService.CheckFinish(currentGame)
}

func (gs *GameService) MakeMove(ctx context.Context, gameID, playerID uuid.UUID, row, col int) error {
	game, err := gs.Repo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	// Проверяем, что игра не закончена/
	if game.CurrentField.IsFinished() {
		return domain.ErrGameFinished
	}
	// Определяем, кто должен ходить
	currentPlayerSymbol := gs.GetCurrentPlayer(game.CurrentField) //x или y
	// Находим игрока с таким символом и проверяем, что ходит именно он
	// for _, p := range game.Players {
	// 	if p.PlayerID != playerID || p.Symbol != currentPlayerSymbol {
	// 		return errors.New("not your turn")
	// 	}
	// }
	playerFound := false
	for _, p := range game.Players {
		if p.PlayerID == playerID {
			playerFound = true
			if p.Symbol != currentPlayerSymbol {
				return errors.New("not your turn")
			}
			break
		}
	}

	if !playerFound {
		return errors.New("player not found in this game")
	}

	// Выполняем ход
	if !game.CurrentField.MakeMove(row, col, currentPlayerSymbol) {
		return domain.ErrFailedToFindMove
	}
	// Проверяем статус после хода
	if game.CurrentField.IsFinished() {
		// ничья = 3, выиграл Х = 1, выиграл 0 = 2
		if game.CurrentField.CheckResult() == 3 {
			game.GameState = domain.Draw
		} else {
			var winner string
			for i := 0; i < 2; i++ {
				if game.CurrentField.CheckResult() == 2 && game.Players[i].Symbol == domain.PlayerO || game.CurrentField.CheckResult() == 1 && game.Players[i].Symbol == domain.PlayerX {
					winner = game.Players[i].PlayerID.String()
					break
				}
			}
			game.GameState = domain.UUIDWins + winner
		}
	} else {
		// Переключаем ход
		var toMove string
		for i := 0; i < 2; i++ {
			if game.Players[i].Symbol != currentPlayerSymbol {
				toMove = game.Players[i].PlayerID.String()
				break
			}
		}
		game.GameState = domain.PlayerToMove + toMove
	}

	return gs.Repo.Save(game)
}

// возвращает лучший ход для указанной игры
func (gs *GameService) GetBestMove(ctx context.Context, gameID uuid.UUID) (int, int, error) {
	game, err := gs.Repo.GetByID(ctx, gameID)
	if err != nil {
		return -1, -1, err
	}
	return gs.NextTurn(*game)
}

// определяет, чей сейчас ход.
func (gs *GameService) GetCurrentPlayer(field *domain.GameField) int {
	return gs.appService.GetCurrentPlayer(field)
}
