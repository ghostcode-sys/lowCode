package cricbuzz

type wicket struct {
	wicketType WicketType
	batsmanOut player
	bowledBy   player
	runoutBy   player
	caughtBy   player
}

func NewWicket(wicketType WicketType, batsmanOut player) *wicket {
	return &wicket{
		wicketType: wicketType,
		batsmanOut: batsmanOut,
	}
}

func (w *wicket) SetBowledBy(bowler player) *wicket {
	w.bowledBy = bowler
	return w
}

func (w *wicket) SetRunoutBy(fielder player) *wicket {
	w.runoutBy = fielder
	return w
}

func (w *wicket) SetCaughtBy(fielder player) *wicket {
	w.caughtBy = fielder
	return w
}

func (w *wicket) GetWicketType() string {
	switch w.wicketType {
	case BOWLED:
		return "BOWLED"
	case CAUGHT:
		return "CAUGHT"
	case LEG_BEFORE_WICKET:
		return "LEG BEFORE WICKET"
	case RUN_OUT:
		return "RUN OUT"
	default:
		return ""
	}
}

func (w *wicket) GetBatsmanOut() string {
	return w.batsmanOut.GetPlayerName()
}

func (w *wicket) GetBowledBy() string {
	if w.bowledBy.GetPlayerName() != "" {
		return w.bowledBy.GetPlayerName()
	}
	return ""
}

func (w *wicket) GetRunoutBy() string {
	if w.runoutBy.GetPlayerName() != "" {
		return w.runoutBy.GetPlayerName()
	}
	return ""
}

func (w *wicket) GetCaughtBy() string {
	if w.caughtBy.GetPlayerName() != "" {
		return w.caughtBy.GetPlayerName()
	}
	return ""
}
