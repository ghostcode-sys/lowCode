package cricbuzz

import "time"

type match struct {
	matchid     string
	Title       string
	Venue       string
	Date        time.Time
	matchFormat MatchFormat
	teamA       *team
	teamB       *team
	innings     []*innings
}

func NewMatch(matchid string, title string, venue string, date time.Time, matchFormat MatchFormat, teamA *team, teamB *team) *match {
	return &match{
		matchid:     matchid,
		Title:       title,
		Venue:       venue,
		Date:        date,
		matchFormat: matchFormat,
		teamA:       teamA,
		teamB:       teamB,
		innings:     make([]*innings, 0),
	}
}

func (m *match) AddInnings(innings *innings) *match {
	m.innings = append(m.innings, innings)
	return m
}

func (m *match) GetMatchID() string {
	return m.matchid
}

func (m *match) GetTitle() string {
	return m.Title
}

func (m *match) GetVenue() string {
	return m.Venue
}

func (m *match) GetDate() time.Time {
	return m.Date
}

func (m *match) GetMatchFormat() MatchFormat {
	return m.matchFormat
}

func (m *match) GetTeamA() string {
	return m.teamA.GetTeamName()
}

func (m *match) GetTeamB() string {
	return m.teamB.GetTeamName()
}

func (m *match) GetInnings() []*innings {
	return m.innings
}

