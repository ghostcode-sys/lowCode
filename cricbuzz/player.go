package cricbuzz

type player struct {
	playerId     string
	playerName   string
	playerRole   PlayerRole
	battingStyle DominateSide
	bowlingStyle BowlingStyle
}

type BowlingStyle struct {
	dominateSide DominateSide
	bowlingType  BowlingType
}

func NewPlayer(name string, id string) *player {
	return &player{
		playerId:   id,
		playerName: name,
	}
}

func (p *player) SetPlayerRole(role PlayerRole) *player {
	p.playerRole = role
	return p
}

func (p *player) SetBattingStyle(style DominateSide) *player {
	p.battingStyle = style
	return p
}

func (p *player) SetBowlingStyle(bowlStyle BowlingType, dominatehand DominateSide) *player {
	p.bowlingStyle.dominateSide = dominatehand
	p.bowlingStyle.bowlingType = bowlStyle
	return p
}

func (p *player) GetPlayerName() string {
	return p.playerName
}

func (p *player) GetPlayerRole() string {
	switch p.playerRole {
	case BATSMAN:
		return "BATSMAN"
	case BOWLER:
		return "BOWLER"
	case WICKET_KEEPER:
		return "WICKET KEEPER"
	case ALL_ROUNDER:
		return "ALL_ROUNDER"
	default:
		return ""
	}
}

func (p *player) GetBattingStyle() string {
	switch p.battingStyle {
	case RIGHTHANDY:
		return "Right hand"
	case LEFTHANDY:
		return "Left hand"
	default:
		return ""
	}
}

func (p *player) GetBowlingStyle() string {
	s := ""
	switch p.bowlingStyle.dominateSide {
	case RIGHTHANDY:
		s += "Right hand"
	case LEFTHANDY:
		s += "Left hand"
	default:
		return ""
	}
	switch p.bowlingStyle.bowlingType {
	case FAST:
		s += " fast pacer"
	case MEDIUM:
		s += " medium pacer"
	case SPIN:
		s += " spinner"
	default:
		return ""
	}
	return s
}
