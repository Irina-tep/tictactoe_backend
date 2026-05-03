package domain

type GameField struct {
	Field [3][3]int
}


// возвращает список координат пустых клеток
func (g GameField) GetEmptyCells() [][2]int {
	var cells [][2]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Field[i][j] == Empty {
				cells = append(cells, [2]int{i, j})
			}
		}
	}
	return cells
}

// выполняет ход игрока в указанную позицию, Возвращает true, если ход допустим
func (g *GameField) MakeMove(row, col, player int) bool {
	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		return false
	}
	if g.Field[row][col] != Empty {
		return false
	}
	g.Field[row][col] = player
	return true
}

// проверяет равны ли значения по вертикали, диагонали и горизонтали, возвращает состояние игры: ничья = 3, выиграл Х = 1, выиграл у = 2
func (g GameField) CheckResult() int {

	for i := 0; i < 3; i++ {
		if g.Field[i][0] != Empty && g.Field[i][0] == g.Field[i][1] && g.Field[i][1] == g.Field[i][2] {
			if g.Field[i][0] == PlayerX {
				return ResultWinX
			} else {
				return ResultWinO
			}
		}
	}

	for j := 0; j < 3; j++ {
		if g.Field[0][j] != Empty && g.Field[0][j] == g.Field[1][j] && g.Field[1][j] == g.Field[2][j] {
			if g.Field[0][j] == PlayerX {
				return ResultWinX
			} else {
				return ResultWinO
			}
		}
	}

	if g.Field[0][0] != Empty && g.Field[0][0] == g.Field[1][1] && g.Field[1][1] == g.Field[2][2] {
		if g.Field[0][0] == PlayerX {
			return ResultWinX
		} else {
			return ResultWinO
		}
	}
	if g.Field[0][2] != Empty && g.Field[0][2] == g.Field[1][1] && g.Field[1][1] == g.Field[2][0] {
		if g.Field[0][2] == PlayerX {
			return ResultWinX
		} else {
			return ResultWinO
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Field[i][j] == Empty {
				return ResultNotFinished
			}
		}
	}
	return ResultDraw
}

func (g GameField) IsFinished() bool {
	result := g.CheckResult()
	return result != ResultNotFinished
}

// возвращает оценку позиции для минимакса
func (g GameField) Evaluate() int {
	result := g.CheckResult()
	switch result {
	case ResultWinX:
		return 10
	case ResultWinO:
		return -10
	default:
		return 0
	}
}


