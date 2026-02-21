package tictactoe

import "math"

type BotPlayer struct {
	Name        string
	Symbol      Symbol
	OpponentSym Symbol
}

func NewBotPlayer(name string, symbol Symbol) *BotPlayer {
	opp := X
	if symbol == X {
		opp = O
	}
	return &BotPlayer{
		Name:        name,
		Symbol:      symbol,
		OpponentSym: opp,
	}
}

func (b *BotPlayer) GetName() string {
	return b.Name
}

func (b *BotPlayer) GetSymbol() Symbol {
	return b.Symbol
}

func (b *BotPlayer) MakeMove(board *Board) Move {
	bestScore := math.MinInt32
	bestMove := Move{-1, -1, b.Symbol}

	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			if board.IsValidMove(i, j) {
				board.Grid[i][j].SetSymbol(b.Symbol)
				score := b.minimax(board, false)
				board.Grid[i][j].SetSymbol(EMPTY)

				if score > bestScore {
					bestScore = score
					bestMove = Move{i, j, b.Symbol}
				}
			}
		}
	}

	return bestMove
}

func (b *BotPlayer) minimax(board *Board, isMaximizing bool) int {
	winner := board.CheckWinner()

	if winner == b.Symbol {
		return 10
	}
	if winner == b.OpponentSym {
		return -10
	}
	if board.IsFull() {
		return 0
	}

	if isMaximizing {
		best := math.MinInt32
		for i := 0; i < board.Size; i++ {
			for j := 0; j < board.Size; j++ {
				if board.IsValidMove(i, j) {
					board.Grid[i][j].SetSymbol(b.Symbol)
					score := b.minimax(board, false)
					board.Grid[i][j].SetSymbol(EMPTY)
					best = max(best, score)
				}
			}
		}
		return best
	} else {
		best := math.MaxInt32
		for i := 0; i < board.Size; i++ {
			for j := 0; j < board.Size; j++ {
				if board.IsValidMove(i, j) {
					board.Grid[i][j].SetSymbol(b.OpponentSym)
					score := b.minimax(board, true)
					board.Grid[i][j].SetSymbol(EMPTY)
					best = min(best, score)
				}
			}
		}
		return best
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
