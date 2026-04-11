// Создайте структуру, которая реализует интерфейс сервиса и принимает в качестве параметра интерфейс репозитория для работы со структурой хранения.
package datasource

import (
	"domain"
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

func (gs *GameService) MakeMove(gameID string, row, col int) error {
	game, err := gs.Repo.GetByID(gameID)
	if err != nil {
		return err
	}
	if game.CurrentField.IsFinished() {
		return domain.ErrGameFinished
	}
	currentPlayer := gs.GetCurrentPlayer(game.CurrentField)
	if !game.CurrentField.MakeMove(row, col, currentPlayer) {
		return domain.ErrFailedToFindMove
	}
	return gs.Repo.Save(game)
}

// возвращает лучший ход для указанной игры
func (gs *GameService) GetBestMove(gameID string) (int, int, error) {
	game, err := gs.Repo.GetByID(gameID)
	if err != nil {
		return -1, -1, err
	}
	return gs.NextTurn(*game)
}

// определяет, чей сейчас ход.
func (gs *GameService) GetCurrentPlayer(field *domain.GameField) int {
	return gs.appService.GetCurrentPlayer(field)
}
