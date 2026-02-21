package tictactoe

type Symbol int

const (
	EMPTY Symbol = iota
	X
	O
)

type GameStatus int

const (
	IN_PROGRESS GameStatus = iota
	DRAW
	WON
)