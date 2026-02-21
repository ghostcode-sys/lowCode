package tictactoe

type Player interface {
	GetName() string
	GetSymbol() Symbol
	MakeMove(b *Board) Move
}