package cricbuzz

type toss struct {
	tossWinner   team
	tossDecision TossDecisionType
	tossResult   TossType
	tossCalledBy player
	tossedBy     player
}

func NewToss(tossWinner team, tossResult TossType, tossCalledBy player, tossedBy player, tossDecision TossDecisionType) *toss {
	return &toss{
		tossWinner:   tossWinner,
		tossDecision: tossDecision,
		tossResult:   tossResult,
		tossCalledBy: tossCalledBy,
		tossedBy:     tossedBy,
	}
}

func (t *toss) GetTossWinner() string {
	return t.tossWinner.GetTeamName()
}

func (t *toss) GetTossResult() string {
	switch t.tossResult {
	case HEAD:
		return "Head"
	case TAIL:
		return "Tail"
	default:
		return "Unknown"
	}
}
func (t *toss) GetTossCalledBy() string {
	return t.tossCalledBy.GetPlayerName()
}

func (t *toss) GetTossedBy() string {
	return t.tossedBy.GetPlayerName()
}

func (t *toss) GetTossDecision() string {
	switch t.tossDecision {
	case BAT:
		return "Bat"
	case BOWL:
		return "Bowl"
	default:
		return "Unknown"
	}
}
