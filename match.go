package elo

import (
	"math"
)

type Player interface {
	GetElo() float64
	SetElo(float64)
}

type Match struct {
	PlayerOne   Player
	PlayerTwo   Player
	finished    bool
	strategy    StrategyFunc
	k           float64
	deviation   float64
	scoreWeight float64
	ignoreDraws bool
}

// Args p1 and p2 should be non-nil pointers.
func (c *Calculator) NewMatch(p1, p2 Player) *Match {
	m := new(Match)
	m.PlayerOne = p1
	m.PlayerTwo = p2
	m.strategy = c.strategy
	m.k = c.k
	m.deviation = c.deviation
	m.scoreWeight = c.scoreWeight
	m.ignoreDraws = c.ignoreDraws
	return m
}

// Contains each player's odds to win, as a number between 0-1, where
// 0 is a 0% chance to win, and 1 is a 100% chance to win.
type MatchOdds struct {
	PlayerOneOdds float64
	PlayerTwoOdds float64
}

type MatchResult struct {
	// Required when using a non-scored strategy.
	// Ignored when using a scored strategy.
	Outcome MatchOutcome

	// Required when using a scored strategy.
	// Ignored when using a non-scored strategy.
	PlayerOneScore int
	// Required when using a scored strategy.
	// Ignored when using a non-scored strategy.
	PlayerTwoScore int
}

type MatchOutcome int

const (
	OutcomeDraw         = MatchOutcome(0)
	OutcomePlayerOneWin = MatchOutcome(1)
	OutcomePlayerTwoWin = MatchOutcome(2)
)

// Set a strategy to be used for this match only.
func (m *Match) SetStrategy(sf StrategyFunc) {
	m.strategy = sf
}

// K must be non-negative. If a negative value is provided, K will be unchanged.
func (m *Match) SetKValue(k float64) {
	if k < 0 {
		return
	}
	m.k = k
}
func (m *Match) GetKValue() float64 {
	return m.k
}

// Score weight must be non-negative.
// If a negative value is provided, ScoreWeight will be unchanged.
func (m *Match) SetScoreWeight(w float64) {
	if w < 0 {
		return
	}
	m.scoreWeight = w
}
func (m *Match) GetScoreWeight() float64 {
	return m.scoreWeight
}

// Deviation must be non-negative.
// If a negative value is provided, the deviation will be unchanged.
func (m *Match) SetDeviation(d float64) {
	if d < 0 {
		return
	}
	m.deviation = d
}
func (m *Match) GetDeviation() float64 {
	return m.deviation
}
func (m *Match) IgnoreDraws(ignore bool) {
	m.ignoreDraws = ignore
}
func (m *Match) GetIgnoreDraws() bool {
	return m.ignoreDraws
}

func (m Match) GetOdds() *MatchOdds {
	R1 := math.Pow(10, m.PlayerOne.GetElo()/m.deviation)
	R2 := math.Pow(10, m.PlayerTwo.GetElo()/m.deviation)

	E1 := R1 / (R1 + R2)
	E2 := R2 / (R1 + R2)
	return &MatchOdds{
		PlayerOneOdds: E1,
		PlayerTwoOdds: E2,
	}
}

// Adjusts the Match's player's elo according to who won the match.
// Can only be called once. Any subsequent calls on the same match will result in no changes
// to the players' elo ratings.
//
// Note: Play() uses a reference to the Match's calculator to determine the new elos. If the
// calculator no longer exists, the function will panic.
func (m *Match) Play(result *MatchResult) {
	if m.finished ||
		((result.Outcome == OutcomeDraw) &&
			m.ignoreDraws &&
			(result.PlayerOneScore == result.PlayerTwoScore)) {
		return
	}
	n1, n2 := m.strategy(&CalculatorInput{
		PlayerOne:      m.PlayerOne.GetElo(),
		PlayerTwo:      m.PlayerTwo.GetElo(),
		PlayerOneScore: result.PlayerOneScore,
		PlayerTwoScore: result.PlayerTwoScore,
		Outcome:        result.Outcome,
		Deviation:      m.deviation,
		ScoreWeight:    m.scoreWeight,
		K:              m.k,
	})
	// fmt.Printf("%+v", CalculatorInput{
	// 	PlayerOne:      m.PlayerOne.GetElo(),
	// 	PlayerTwo:      m.PlayerTwo.GetElo(),
	// 	PlayerOneScore: result.PlayerOneScore,
	// 	PlayerTwoScore: result.PlayerTwoScore,
	// 	Outcome:        result.Outcome,
	// 	Deviation:      m.deviation,
	// 	ScoreWeight:    m.scoreWeight,
	// 	K:              m.k,
	// })
	m.PlayerOne.SetElo(n1)
	m.PlayerTwo.SetElo(n2)
	m.finished = true
}

// Returns how much player one stands to gain if they win.
// Equivalent to how much player two will lose if they lose.
//
// Note: May not be accurate when using a scored strategy.
func (m Match) PlayerOneGain() float64 {
	n1, _ := m.strategy(&CalculatorInput{
		PlayerOne: m.PlayerOne.GetElo(),
		PlayerTwo: m.PlayerTwo.GetElo(),
		Outcome:   OutcomePlayerOneWin,
		K:         m.k,
		Deviation: m.deviation,
	})
	return n1 - m.PlayerOne.GetElo()
}

// Returns how much player two stands to gain if they win.
// Equivalent to how much player one will lose if they lose.
//
// Note: May not be accurate when using a scored strategy.
func (m Match) PlayerTwoGain() float64 {
	_, n2 := m.strategy(&CalculatorInput{
		PlayerOne: m.PlayerOne.GetElo(),
		PlayerTwo: m.PlayerTwo.GetElo(),
		Outcome:   OutcomePlayerTwoWin,
		K:         m.k,
		Deviation: m.deviation,
	})
	return n2 - m.PlayerTwo.GetElo()
}
