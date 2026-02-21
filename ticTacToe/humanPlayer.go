package tictactoe

import "fmt"

type HumanPlayer struct {
	name   string
	symbol Symbol
}

func NewHumanPlayer(name string, symbol Symbol) *HumanPlayer {
	return &HumanPlayer{name: name, symbol: symbol}
}

func (p *HumanPlayer) GetName() string {
	return p.name
}

func (p *HumanPlayer) GetSymbol() Symbol {
	return p.symbol
}

func (p *HumanPlayer) MakeMove(b *Board) Move {
	fmt.Print("Row: ")
	var row int
	var col int
	fmt.Scan(&row)
	fmt.Print("Col: ")
	fmt.Scan(&col)
	return Move{Row: row, Col: col, Symbol: p.symbol}
}
