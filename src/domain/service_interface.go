package domain

type Interface interface {
	NextTurn(currentGame CurrentGame) (int, int, error)
	CheckField(currentGame CurrentGame, expectedGame CurrentGame) (bool, error)
	CheckFinish(currentGame CurrentGame) (int, int, bool)
}

// GameServiceInterface - структура, реализующая интерфейс Interface
type GameServiceInterface struct{}
