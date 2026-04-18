package cricbuzz

type MatchFormat uint8

const (
	T20 MatchFormat = iota
	ODI
	TEST
)

type MatchStatus uint8

const (
	UPCOMING MatchStatus = iota
	LIVE
	COMPLETED
	ABANDONED
)

type PlayerRole uint8

const (
	BATSMAN PlayerRole = iota
	BOWLER
	ALL_ROUNDER
	WICKET_KEEPER
)

type DominateSide uint8

const (
	RIGHTHANDY DominateSide = iota
	LEFTHANDY
)

type BowlingType uint8

const (
	SPIN BowlingType = iota
	MEDIUM
	FAST
)

type Extras uint8

const (
	NO_BALL Extras = iota
	WIDE
	BYE
	LEG_BYE
	PENALTY_RUN
)

type WicketType uint8

const (
	BOWLED WicketType = iota
	CAUGHT
	LEG_BEFORE_WICKET
	RUN_OUT
	STUMPED
	HIT_WICKET
)

type TossType uint8

const (
	HEAD TossType = iota
	TAIL
)

type TossDecisionType uint8

const (
	BAT TossDecisionType = iota
	BOWL
)
