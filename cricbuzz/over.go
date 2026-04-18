package cricbuzz

type over struct {
	overNumber int
	balls      []*ball
}

func NewOver(overNumber int) *over {
	return &over{
		overNumber: overNumber,
		balls:      make([]*ball, 0),
	}
}

func (o *over) AddBall(ball *ball) *over {
	o.balls = append(o.balls, ball)
	return o
}

func (o *over) GetOverNumber() int {
	return o.overNumber
}

func (o *over) GetBalls() []*ball {
	return o.balls
}

func (o *over) GetTotalRuns() int {
	totalRuns := 0
	for _, ball := range o.balls {
		totalRuns += ball.GetRunsScored()
	}
	return totalRuns
}

func (o *over) GetWickets() []*wicket {
	wickets := make([]*wicket, 0)
	for _, ball := range o.balls {
		if ball.GetWicket() != nil {
			wickets = append(wickets, ball.GetWicket())
		}
	}
	return wickets
}

func (o *over) GetGoodBalls() int {
	goodBalls := 0
	for _, ball := range o.balls {
		if ball.GetExtra() != WIDE && ball.GetExtra() != NO_BALL {
			goodBalls++
		}
	}
	return goodBalls
}
