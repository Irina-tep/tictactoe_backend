package domain

import (
	"math"
)

func NewService() GameServiceInterface {
	return GameServiceInterface{}
}

// вычисляет следующий ход для текущего игрока с использованием алгоритма минимакса
// Возвращает оптимальные координаты хода или ошибку, если игра завершена или поле невалидно
func (s *GameServiceInterface) NextTurn(currentGame CurrentGame) (int, int, error) {
	if currentGame.CurrentField.IsFinished() {
		return -1, -1, ErrGameFinished
	}

	currentPlayer := s.GetCurrentPlayer(currentGame.CurrentField)

	bestScore := 0
	if currentPlayer == PlayerX {
		bestScore = math.MinInt32
	} else {
		bestScore = math.MaxInt32
	}
	var bestMove [2]int
	bestMove[0] = -1
	bestMove[1] = -1

	emptyCells := currentGame.CurrentField.GetEmptyCells()
	if len(emptyCells) == 0 {
		return -1, -1, ErrNoMovies
	}

	for _, cell := range emptyCells {
		row, col := cell[0], cell[1]

		// Создаем копию поля для симуляции хода
		newField := GameField{
			Field: currentGame.CurrentField.Field,
		}
		newField.MakeMove(row, col, currentPlayer)

		score := s.minimax(&newField, 0, false, currentPlayer == PlayerX)

		// Выбираем лучший ход в зависимости от игрока
		if currentPlayer == PlayerX {
			if score > bestScore {
				bestScore = score
				bestMove[0], bestMove[1] = row, col
			}
		} else {
			if score < bestScore {
				bestScore = score
				bestMove[0], bestMove[1] = row, col
			}
		}
	}

	if bestMove[0] == -1 || bestMove[1] == -1 {
		return -1, -1, ErrFailedToFindMove
	}

	return bestMove[0], bestMove[1], nil
}

// алгоритм минимакса
func (s *GameServiceInterface) minimax(field *GameField, depth int, isMaximizing bool, isXPlayer bool) int {
	result := field.CheckResult()
	if result == ResultWinX {
		return 10 - depth //X выиграл - чем меньше ходов, тем выше оценка
	}
	if result == ResultWinO {
		return depth - 10 //O выиграл - чем больше ходов, тем лучше (меньше отрицательная оценка)
	}
	if result == ResultDraw {
		return 0
	}

	var currentPlayer int
	if isMaximizing {
		if isXPlayer {
			currentPlayer = PlayerX
		} else {
			currentPlayer = PlayerO
		}
	} else {
		if isXPlayer {
			currentPlayer = PlayerO
		} else {
			currentPlayer = PlayerX
		}
	}

	// Получаем все возможные ходы
	emptyCells := field.GetEmptyCells()

	if isMaximizing {
		bestScore := math.MinInt32
		for _, cell := range emptyCells {
			row, col := cell[0], cell[1]

			// Создаем временную копию для симуляции хода
			tempField := GameField{
				Field: field.Field,
			}
			tempField.MakeMove(row, col, currentPlayer)

			score := s.minimax(&tempField, depth+1, false, isXPlayer)
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	} else {
		bestScore := math.MaxInt32
		for _, cell := range emptyCells {
			row, col := cell[0], cell[1]

			// Создаем временную копию для симуляции хода
			tempField := GameField{
				Field: field.Field,
			}
			tempField.MakeMove(row, col, currentPlayer)

			score := s.minimax(&tempField, depth+1, true, isXPlayer)
			if score < bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
}

// определяет, чей сейчас ход если ход
func (s *GameServiceInterface) GetCurrentPlayer(field *GameField) int {
	countX := 0
	countO := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if field.Field[i][j] == PlayerX {
				countX++
			} else if field.Field[i][j] == PlayerO {
				countO++
			}
		}
	}

	if countX <= countO {
		return PlayerX
	}
	return PlayerO
}

// проверяет игровое поле на целостность и соответствие ожидаемому состоянию
func (s *GameServiceInterface) CheckField(currentGame CurrentGame, expectedGame CurrentGame) (bool, error) {

	if currentGame.ID != expectedGame.ID {
		return false, ErrInvalidGameID
	}

	if len(currentGame.CurrentField.Field) != len(expectedGame.CurrentField.Field) {
		return false, ErrInvalidFieldSize
	}

	for i := 0; i < 3; i++ {
		if len(currentGame.CurrentField.Field[i]) != len(expectedGame.CurrentField.Field[i]) {
			return false, ErrInvalidRows
		}
		for j := 0; j < 3; j++ {
			if currentGame.CurrentField.Field[i][j] != expectedGame.CurrentField.Field[i][j] {
				return false, ErrInvalidFieldContents
			}
		}
	}

	return true, nil
}

// проверяет, завершена ли игра на текущем игровом поле
func (s *GameServiceInterface) CheckFinish(currentGame CurrentGame) (int, int, bool) {
	result := currentGame.CurrentField.CheckResult()

	switch result {
	case ResultWinX:
		return 1, 0, true
	case ResultWinO:
		return 0, 1, true
	case ResultDraw:
		return 0, 0, true
	default:
		return 0, 0, false
	}
}
