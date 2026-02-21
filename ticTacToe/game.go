package tictactoe

import "fmt"

type Game struct {
	Board       *Board
	Players     []Player
	CurrentTurn int
	Status      GameStatus
	Winner      Player
}

func NewGame(board *Board, players []Player) *Game {
	return &Game{
		Board:       board,
		Players:     players,
		CurrentTurn: 0,
		Status:      IN_PROGRESS,
	}
}

func (g *Game) Start() {
	g.Status = IN_PROGRESS
}

func (g *Game) PlayTurn() {
	player := g.Players[g.CurrentTurn]
	move := player.MakeMove(g.Board)

	g.Board.ApplyMove(move.Row, move.Col, move.Symbol)

	winner := g.Board.CheckWinner()
	if winner != EMPTY {
		g.Status = WON
		g.Winner = player
		return
	}

	if g.Board.IsFull() {
		g.Status = DRAW
		return
	}

	g.SwitchTurn()
}

func (g *Game) SwitchTurn() {
	fmt.Println(g.Board.String())
	fmt.Print("\n+++++++++++++++++\n\n")
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)
}
