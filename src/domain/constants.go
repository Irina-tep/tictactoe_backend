package domain

import "errors"

// значки, которыми будут ходить пользователи
const (
	Empty   = 0
	PlayerX = 1
	PlayerO = 2
)

const (
	ResultNotFinished = 0
	ResultWinX        = 1
	ResultWinO        = 2
	ResultDraw        = 3
)

const (
	GameTypePvP = "pvp" // игрок против игрока
	GameTypePvC = "pvc" // игрок против компьютера
)

// состояния для текущей игры
const (
	WaitingForPlayers string = "waiting"        //Ожидание игроков
	PlayerToMove      string = "player_to_move" //Ход игрока с UUID
	UUIDWins          string = "player_wins"    //Победа игрока с UUID.
	Draw              string = "draw"           //Ничья
)

var (
	ErrGameFinished         = errors.New("game finished")
	ErrGameNotFound         = errors.New("game not found")
	ErrNoMovies             = errors.New("no moves available")
	ErrFailedToFindMove     = errors.New("failed to find the optimal move")
	ErrInvalidGameID        = errors.New("ID does not match")
	ErrInvalidFieldSize     = errors.New("field sizes do not match")
	ErrInvalidRows          = errors.New("the sizes of the field rows do not match")
	ErrInvalidFieldContents = errors.New("field contents do not match")
)
