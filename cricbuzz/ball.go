package cricbuzz

type ball struct {
	ballNumber int
	bowler     player
	striker    player
	nonStriker player
	runsScored int
	isWicket   *wicket
	Comment    string
	extra      Extras
}

func Newball(ballNumber int, bowler player, striker player, nonStriker player) *ball {
	return &ball{
		ballNumber: ballNumber,
		bowler:     bowler,
		striker:    striker,
		nonStriker: nonStriker,
	}
}

func (b *ball) SetRunsScored(runs int) *ball {
	b.runsScored = runs
	return b
}

func (b *ball) SetWicket(wicket *wicket) *ball {
	b.isWicket = wicket
	return b
}

func (b *ball) SetComment(comment string) *ball {
	b.Comment = comment
	return b
}

func (b *ball) SetExtra(extra Extras) *ball {
	b.extra = extra
	return b
}

func (b *ball) GetRunsScored() int {
	return b.runsScored
}

func (b *ball) GetWicket() *wicket {
	return b.isWicket
}

func (b *ball) GetComment() string {
	return b.Comment
}

func (b *ball) GetExtra() Extras {
	return b.extra
}

func (b *ball) GetBowler() string {
	return b.bowler.GetPlayerName()
}

func (b *ball) GetStriker() string {
	return b.striker.GetPlayerName()
}

func (b *ball) GetNonStriker() string {
	return b.nonStriker.GetPlayerName()
}

func (b *ball) GetBallNumber() int {
	return b.ballNumber
}
