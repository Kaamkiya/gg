package engine

type Engine struct {
	ai AI
}

func NewEngine(depth int) *Engine {
	engine := &Engine{}
	mcts := NewMCTS(engine, depth)
	engine.ai = mcts

	return engine
}

func (e *Engine) GetLegalMoves(board *Board) []int {
	var moves []int
	for i, cell := range board.Cells {
		if cell == EMPTY {
			moves = append(moves, i)
		}
	}
	return moves
}

func (e *Engine) PlayMove(board *Board, player int, move int) error {
	return board.SetCell(move, player)
}

func (e *Engine) GetOpponent(player int) int {
	return -player
}

func (e *Engine) CheckGameOver(board *Board, lastMove int) (bool, int) {
	if lastMove == -1 {
		return false, 0
	}

	if e.CheckWin(board, lastMove) {
		absValue := P1 * P2 * -1
		return true, absValue
	}

	if len(e.GetLegalMoves(board)) == 0 {
		return true, 0
	}

	return false, 0
}

func (e *Engine) CheckWin(board *Board, lastMove int) bool {
	player, err := board.GetCell(lastMove)
	if err != nil {
		panic(err)
	}
	if player == EMPTY {
		return false
	}

	row, col, err := board.GetRowCol(lastMove)
	if err != nil {
		panic(err)
	}

	if e.checkRow(board, row, player) {
		return true
	}

	if e.checkCol(board, col, player) {
		return true
	}

	if e.checkDiagonal(board, player) {
		return true
	}

	return false
}

func (e *Engine) checkRow(board *Board, row, player int) bool {
	for i := 0; i < board.Size; i++ {
		cell, err := board.GetCell(row*board.Size + i)
		if err != nil {
			panic(err)
		}

		if cell != player {
			return false
		}
	}

	return true
}

func (e *Engine) checkCol(board *Board, col, player int) bool {
	for i := 0; i < board.Size; i++ {
		cell, err := board.GetCell(i*board.Size + col)
		if err != nil {
			panic(err)
		}

		if cell != player {
			return false
		}
	}

	return true
}

func (e *Engine) checkDiagonal(board *Board, player int) bool {
	sum := 0
	// Left to right
	for i := 0; i < board.Size; i++ {
		cell, err := board.GetCell(i*board.Size + i)
		if err != nil {
			panic(err)
		}

		if cell == player {
			sum += player
		}
	}

	if sum == board.Size*player {
		return true
	}

	sum = 0
	// Right to left
	for i := 0; i < board.Size; i++ {
		cell, err := board.GetCell(i*board.Size + board.Size - i - 1)
		if err != nil {
			panic(err)
		}

		if cell == player {
			sum += player
		}
	}

	if sum == board.Size*player {
		return true
	}

	return false
}
