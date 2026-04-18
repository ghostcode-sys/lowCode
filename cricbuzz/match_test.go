package cricbuzz

import (
	"testing"
	"time"
)

func TestNewMatchAndGetters(t *testing.T) {
	teamA := NewTeam(1, "India")
	teamB := NewTeam(2, "Australia")
	matchDate := time.Date(2026, time.April, 18, 14, 0, 0, 0, time.UTC)
	m := NewMatch("M100", "India vs Australia", "Mumbai", matchDate, ODI, teamA, teamB)

	if got := m.GetMatchID(); got != "M100" {
		t.Fatalf("GetMatchID() = %q, want %q", got, "M100")
	}

	if got := m.GetTitle(); got != "India vs Australia" {
		t.Fatalf("GetTitle() = %q, want %q", got, "India vs Australia")
	}

	if got := m.GetVenue(); got != "Mumbai" {
		t.Fatalf("GetVenue() = %q, want %q", got, "Mumbai")
	}

	if !m.GetDate().Equal(matchDate) {
		t.Fatalf("GetDate() = %v, want %v", m.GetDate(), matchDate)
	}

	if got := m.GetMatchFormat(); got != ODI {
		t.Fatalf("GetMatchFormat() = %v, want %v", got, ODI)
	}

	if got := m.GetTeamA(); got != "India" {
		t.Fatalf("GetTeamA() = %q, want %q", got, "India")
	}

	if got := m.GetTeamB(); got != "Australia" {
		t.Fatalf("GetTeamB() = %q, want %q", got, "Australia")
	}

	if got := len(m.GetInnings()); got != 0 {
		t.Fatalf("GetInnings() length = %d, want %d", got, 0)
	}
}

func TestAddInnings(t *testing.T) {
	teamA := NewTeam(1, "India")
	teamB := NewTeam(2, "Australia")
	m := NewMatch("M101", "India vs Australia", "Mumbai", time.Now(), T20, teamA, teamB)

	inn := NewInnings(1, teamA, teamB)
	m.AddInnings(inn)

	innings := m.GetInnings()
	if got := len(innings); got != 1 {
		t.Fatalf("GetInnings() length = %d, want %d", got, 1)
	}

	if got := innings[0].GetInningsNumber(); got != 1 {
		t.Fatalf("GetInnings()[0].GetInningsNumber() = %d, want %d", got, 1)
	}

	if got := innings[0].GetBattingTeam(); got != "India" {
		t.Fatalf("GetInnings()[0].GetBattingTeam() = %q, want %q", got, "India")
	}

	if got := innings[0].GetBowlingTeam(); got != "Australia" {
		t.Fatalf("GetInnings()[0].GetBowlingTeam() = %q, want %q", got, "Australia")
	}
}
