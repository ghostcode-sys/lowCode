package service

import "time"

type BallInterface interface {
	SetRunsScored(runs int) *ball
	SetWicket(wicket *wicket) *ball
	SetComment(comment string) *ball
	SetExtra(extra Extras) *ball
	GetRunsScored() int
	GetWicket() *wicket
	GetComment() string
	GetExtra() Extras
	GetBowler() string
	GetStriker() string
	GetNonStriker() string
	GetBallNumber() int
}

type InningsInterface interface {
	SetScoreCard(scoreCard *scoreCard) *innings
	AddOver(over *over) *innings
	GetInningsNumber() int
	GetScoreCard() *scoreCard
	GetBattingTeam() string
	GetBowlingTeam() string
	GetTotalRuns() int
	GetTotalWickets() int
	GetOverPlayed() float64
	GetOvers() []*over
}

type MatchInterface interface {
	AddInnings(innings *innings) *match
	GetMatchID() string
	GetTitle() string
	GetVenue() string
	GetDate() time.Time
	GetMatchFormat() MatchFormat
	GetTeamA() string
	GetTeamB() string
	GetInnings() []*innings
	GetStatus() MatchStatus
	SetStatus(status MatchStatus) *match
	GetCurrentScore() (int, int, float64)
	GetMatchResult() string
}

type OverInterface interface {
	AddBall(ball *ball) *over
	GetOverNumber() int
	GetBalls() []*ball
	GetTotalRuns() int
	GetWickets() []*wicket
	GetGoodBalls() int
}

type PlayerInterface interface {
	SetPlayerRole(role PlayerRole) *player
	SetBattingStyle(style DominateSide) *player
	SetBowlingStyle(bowlStyle string, dominatehand DominateSide) *player
	SetCountry(country string) *player
	GetPlayerName() string
	GetPlayerRole() string
	GetBattingStyle() string
	GetPlayerID() string
	GetBowlingStyle() string
	GetCountry() string
}

type ScoreCardInterface interface {
	AddBattingScoreCard(player *player, runs int, balls int, fours int, sixes int) *scoreCard
	AddBowlingScoreCard(bowler PlayerInterface, overs float64, runsConceded int, wickets int) *scoreCard
	GetBattingScoreCards() []*battingScoreCard
	GetBowlingScoreCards() []*bowlingScoreCard
	UpdateBattingScoreCard(player *player, runs int, isFour bool, isSix bool)
	UpdateBowlingScoreCard(player *player, runsConceded int, isWicket bool)
}

type TeamInterface interface {
	AddPlayerToSquad(player PlayerInterface) *team
	SetPlayingX11(players []PlayerInterface) *team
	SetCaptain(captain PlayerInterface) *team
	SetWicketKeeper(wicketKeeper PlayerInterface) *team
	GetTeamName() string
	GetSquad() []PlayerInterface
	GetPlayingX11() []PlayerInterface
	GetCaptain() string
	GetWicketKeeper() string
}

type TossInterface interface {
	GetTossWinner() string
	GetTossResult() string
	GetTossCalledBy() string
	GetTossedBy() string
	GetTossDecision() string
}

type WicketInterface interface {
	SetBowledBy(bowler PlayerInterface) *wicket
	SetRunoutBy(fielder PlayerInterface) *wicket
	SetCaughtBy(fielder PlayerInterface) *wicket
	GetWicketType() string
	GetDismissedPlayer() string
	GetBowledBy() string
	GetRunoutBy() string
	GetCaughtBy() string
}
