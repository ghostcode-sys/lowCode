package tictactoe

import (
	"fmt"
	"strings"
)

type Board struct {
	Size int
	Grid [][]*Cell
}

func NewBoard(sz int) *Board{
	grid := make([][]*Cell, sz)

	for i := range sz {
		grid[i] = make([]*Cell, sz)
		for j := range sz {
			grid[i][j] = NewCell(i, j)
		}
	}

	return &Board{
		Size: sz,
		Grid: grid,
	}
}

func (b *Board)IsValidMove(r, c int) bool {
	if r < 0 || r >= b.Size || c < 0|| c >= b.Size{
		return false
	} 
	return b.Grid[r][c].IsEmpty()
}

func (b *Board) ApplyMove(r, c int, s Symbol) bool {
	if !b.IsValidMove(r, c){
		return false
	}

	b.Grid[r][c].SetSymbol(s)
	return true
}

func (b *Board)IsFull() bool {
	for i := range b.Size{
		for j := range b.Size{
			if b.Grid[i][j].IsEmpty(){
				return false
			}
		}
	}
	return true
}

func (b *Board) CheckWinner() Symbol {
	n := b.Size

	// Rows
	for i := 0; i < n; i++ {
		s := b.Grid[i][0].GetSymbol()
		if s == EMPTY {
			continue
		}
		win := true
		for j := 1; j < n; j++ {
			if b.Grid[i][j].GetSymbol() != s {
				win = false
				break
			}
		}
		if win {
			return s
		}
	}

	// Columns
	for j := 0; j < n; j++ {
		s := b.Grid[0][j].GetSymbol()
		if s == EMPTY {
			continue
		}
		win := true
		for i := 1; i < n; i++ {
			if b.Grid[i][j].GetSymbol() != s {
				win = false
				break
			}
		}
		if win {
			return s
		}
	}

	// Diagonals
	s := b.Grid[0][0].GetSymbol()
	if s != EMPTY {
		win := true
		for i := 1; i < n; i++ {
			if b.Grid[i][i].GetSymbol() != s {
				win = false
				break
			}
		}
		if win {
			return s
		}
	}

	s = b.Grid[0][n-1].GetSymbol()
	if s != EMPTY {
		win := true
		for i := 1; i < n; i++ {
			if b.Grid[i][n-1-i].GetSymbol() != s {
				win = false
				break
			}
		}
		if win {
			return s
		}
	}

	return EMPTY
}

func (b *Board) String() string{
	rowSeprator := "\n-----------\n"
	rowString := make([]string, b.Size)
	for i := range b.Size{
		colString := make([]string, b.Size)
		for j := range b.Size{
			val := " "
			if b.Grid[i][j].GetSymbol() == X {
				val = "X"
			} else if b.Grid[i][j].GetSymbol() == O {
				val = "O"
			}
			colString[j] = fmt.Sprintf(" %s ", val)
		}
		rowString[i] = strings.Join(colString, "|")
	}
	s := strings.Join(rowString, rowSeprator)
	return s
}