package service

type battingScoreCard struct {
	player *player
	runs   int
	balls  int
	fours  int
	sixes  int
}

type bowlingScoreCard struct {
	player       *player
	overs        float64
	runsConceded int
	wickets      int
}

type scoreCard struct {
	battingScoreCards []*battingScoreCard
	bowlingScoreCards []*bowlingScoreCard
}

func NewScoreCard() *scoreCard {
	return &scoreCard{
		battingScoreCards: []*battingScoreCard{},
		bowlingScoreCards: []*bowlingScoreCard{},
	}
}

func (s *scoreCard) AddBattingScoreCard(player *player, runs int, balls int, fours int, sixes int) *scoreCard {
	s.battingScoreCards = append(s.battingScoreCards, &battingScoreCard{
		player: player,
		runs:   runs,
		balls:  balls,
		fours:  fours,
		sixes:  sixes,
	})
	return s
}

func (s *scoreCard) AddBowlingScoreCard(player *player, overs float64, runsConceded int, wickets int) *scoreCard {
	s.bowlingScoreCards = append(s.bowlingScoreCards, &bowlingScoreCard{
		player:       player,
		overs:		overs,
		runsConceded: runsConceded,
		wickets:      wickets,
	})
	return s
}

func (s *scoreCard) GetBattingScoreCards() []*battingScoreCard {
	return s.battingScoreCards
}

func (s *scoreCard) GetBowlingScoreCards() []*bowlingScoreCard {
	return s.bowlingScoreCards
}

func (b *battingScoreCard) updateScore(runs int, isFour bool, isSix bool) {
	b.runs += runs
	b.balls += 1
	if isFour {
		b.fours += 1
	}
	if isSix {
		b.sixes += 1
	}
}

func (s *scoreCard) UpdateBattingScoreCard(player *player, runs int, isFour bool, isSix bool) {
	for _, b := range s.battingScoreCards {
		if b.player.GetPlayerID() == player.GetPlayerID() {
			b.updateScore(runs, isFour, isSix)
			return
		}
	}
}

func (b *bowlingScoreCard) updateScore(runsConceded int, isWicket bool) {
	b.runsConceded += runsConceded
	if isWicket {
		b.wickets += 1
	}
}

func (s *scoreCard) UpdateBowlingScoreCard(player *player, runsConceded int, isWicket bool) {
	for _, b := range s.bowlingScoreCards {
		if b.player.GetPlayerID() == player.GetPlayerID() {
			b.updateScore(runsConceded, isWicket)
			return
		}
	}
}

