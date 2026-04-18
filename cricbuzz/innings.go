package cricbuzz

type innings struct {
	inningsNumber int
	battingTeam   *team
	bowlingTeam   *team
	totalRuns     int
	totalWickets  int
	overPlayed    int
	overs         []*over
}

func NewInnings(inningsNumber int, battingTeam *team, bowlingTeam *team) *innings {
	return &innings{
		inningsNumber: inningsNumber,
		battingTeam:   battingTeam,
		bowlingTeam:   bowlingTeam,
		overs:         make([]*over, 0),
	}
}

func (i *innings) AddOver(over *over) *innings {
	i.overs = append(i.overs, over)
	i.totalRuns += over.GetTotalRuns()
	i.totalWickets += len(over.GetWickets())
	i.overPlayed += 1
	return i
}

func (i *innings) GetInningsNumber() int {
	return i.inningsNumber
}

func (i *innings) GetBattingTeam() string {
	return i.battingTeam.GetTeamName()
}

func (i *innings) GetBowlingTeam() string {
	return i.bowlingTeam.GetTeamName()
}

func (i *innings) GetTotalRuns() int {
	return i.totalRuns
}

func (i *innings) GetTotalWickets() int {
	return i.totalWickets
}

func (i *innings) GetOverPlayed() float64 {
	lastOverIndex := len(i.overs) - 1
	if lastOverIndex < 0 {
		return 0
	}
	over := i.overs[lastOverIndex]
	return float64(i.overPlayed-1) + float64(over.GetGoodBalls())/10
}

func (i *innings) GetOvers() []*over {
	return i.overs
}
