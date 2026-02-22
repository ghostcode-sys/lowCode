package cricbuzz

type player struct {
	playerId string
	playerName string
	playerRole PlayerRole
	battingStyle DominateSide
	bowlingStyle BowlingStyle
}

type BowlingStyle struct{
	dominateSide DominateSide
	bowlingType BowlingType
}

func NewPlayer (name string, id string) *player{
	return &player{
		playerId: id,
		playerName: name,
	}
}

func (p *player) SetPlayerRole(role PlayerRole) *player{
	p.playerRole = role
	return p
}

func (p *player)SetBattingStyle(style DominateSide) *player{
	p.battingStyle = style
	return p
}

func (p *player)SetBowlingStyle(bowlStyle BowlingType, dominatehand DominateSide) *player{
	p.bowlingStyle.dominateSide = dominatehand
	p.bowlingStyle.bowlingType = bowlStyle
	return p
}

