package tictactoe

import (
	"fmt"
	"testing"
)

func TestPlayGame(t *testing.T) {
	board := NewBoard(3)
	bot1 := NewBotPlayer("Bot1", X)
	bot2 := NewBotPlayer("Bot2", O)

	game := NewGame(board, []Player{bot1, bot2})

	game.Start()

	for game.Status == IN_PROGRESS {
		game.PlayTurn()
	}

	res := "Draw"

	if game.Status == WON {
		res = fmt.Sprintf("Winner: %s", game.Winner.GetName())
	}
	t.Log(res)
}
