package cricbuzz

type team struct {
	teamid       int
	teamName     string
	squad        []player
	playingX11   []player
	captain      player
	wicketKeeper player
}

func NewTeam(teamid int, teamName string) *team {
	return &team{
		teamid:     teamid,
		teamName:   teamName,
		squad:      make([]player, 0),
		playingX11: make([]player, 0),
	}
}

func (t *team) AddPlayerToSquad(player player) *team {
	t.squad = append(t.squad, player)
	return t
}

func (t *team) SetPlayingX11(players []player) *team {
	t.playingX11 = players
	return t
}

func (t *team) SetCaptain(captain player) *team {
	t.captain = captain
	return t
}

func (t *team) SetWicketKeeper(wicketKeeper player) *team {
	t.wicketKeeper = wicketKeeper
	return t
}

func (t *team) GetTeamName() string {
	return t.teamName
}

func (t *team) GetSquad() []player {
	return t.squad
}

func (t *team) GetPlayingX11() []player {
	return t.playingX11
}

func (t *team) GetCaptain() string {
	return t.captain.GetPlayerName()
}

func (t *team) GetWicketKeeper() string {
	return t.wicketKeeper.GetPlayerName()
}
