package service

import (
	"fmt"
	"time"
)

type match struct {
	matchid     string
	Title       string
	Venue       string
	Date        time.Time
	matchFormat MatchFormat
	teamA       *team
	teamB       *team
	innings     []*innings
	status      MatchStatus
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
		status:      UPCOMING,
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

func (m *match) GetStatus() MatchStatus {
	return m.status
}

func (m *match) SetStatus(status MatchStatus) *match {
	m.status = status
	return m
}

func (m *match) GetCurrentScore() (int, int, float64) {
	if len(m.innings) == 0 {
		return 0, 0, 0.0
	}
	currentInnings := m.innings[len(m.innings)-1]
	runs := currentInnings.GetTotalRuns()
	wickets := currentInnings.GetTotalWickets()
	overs := currentInnings.GetOverPlayed()
	return runs, wickets, overs
}

func (m *match) GetMatchResult() string {
	if m.status != COMPLETED {
		return "Match is not completed yet"
	}
	if len(m.innings) < 2 {
		return "Match result not available"
	}
	teamAScore := 0
	teamBScore := 0
	if len(m.innings) >= 1 {
		teamAScore = m.innings[0].GetTotalRuns()
	}
	if len(m.innings) >= 2 {
		teamBScore = m.innings[1].GetTotalRuns()
	}
	if teamAScore > teamBScore {
		return m.teamA.GetTeamName() + " won by " + fmt.Sprintf("%d", teamAScore-teamBScore) + " runs"
	} else if teamBScore > teamAScore {
		return m.teamB.GetTeamName() + " won by " + fmt.Sprintf("%d", 10-m.innings[1].GetTotalWickets()) + " wickets"
	} else {
		return "Match tied"
	}
}
