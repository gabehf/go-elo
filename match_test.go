package elo_test

import (
	"testing"

	"github.com/gabehf/go-elo"
)

func TestNewMatch(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)

	if m == nil {
		t.Fail()
		t.Log("Expected non-nil match.")
	}
}

func TestOverrides(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.SetStrategy(func(input *elo.CalculatorInput) (r1 float64, r2 float64) {
		return 1, 2
	})

	m.Play(&elo.MatchResult{})

	if !almostEqual(p1.elo, 1.0) || !almostEqual(p2.elo, 2.0) {
		t.Fail()
		t.Logf("Overriden strategy failed.")
	}

	c = elo.NewCalculatorBuilder().WithStrategy(elo.StrategyScored).Build()

	p1.elo = 1600
	p2.elo = 1800
	m = c.NewMatch(p1, p2)
	m.SetKValue(47)
	m.SetDeviation(200)
	m.SetScoreWeight(0.5)
	m.IgnoreDraws(true)

	// bad overrides
	m.SetKValue(-50)
	m.SetDeviation(-200)
	m.SetScoreWeight(-0.13)

	if !almostEqual(m.GetKValue(), 47.0) {
		t.Fail()
		t.Logf("Failed to override K value.")
	}
	if !almostEqual(m.GetScoreWeight(), 0.5) {
		t.Fail()
		t.Logf("Failed to override score weight.")
	}
	if !almostEqual(m.GetDeviation(), 200.0) {
		t.Fail()
		t.Logf("Failed to override deviation.")
	}
	if !m.GetIgnoreDraws() {
		t.Fail()
		t.Logf("Failed to override ignore draws.")
	}

	m.Play(&elo.MatchResult{
		PlayerOneScore: 6,
		PlayerTwoScore: 3,
	})

	if !almostEqual(p1.elo, 1627.219068) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1627.219068, p1.elo)
	}
	if !almostEqual(p2.elo, 1772.780932) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1772.780932, p2.elo)
	}

}

func TestGetOdds(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	o := m.GetOdds()

	if !almostEqual(o.PlayerOneOdds, 0.240253) {
		t.Fail()
		t.Logf("Expected P1 Odds to be %f, got %f\n", 0.240253, o.PlayerOneOdds)
	}
	if !almostEqual(o.PlayerTwoOdds, 0.759747) {
		t.Fail()
		t.Logf("Expected P2 Odds to be %f, got %f\n", 0.759747, o.PlayerTwoOdds)
	}
}

func TestGain(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	g1 := m.PlayerOneGain()
	g2 := m.PlayerTwoGain()

	if g1 < g2 {
		t.Fail()
		t.Log("Lower elo player's gains must be higher than the higher elo player.")
	}
	if g1 < 0 {
		t.Fail()
		t.Log("Elo gain must be above 0.")
	}
	if g2 < 0 {
		t.Fail()
		t.Log("Elo gain must be above 0.")
	}

	if !almostEqual(g1, 24.311902) {
		t.Fail()
		t.Logf("Expected P1 Gain to be %f, got %f\n", 24.311902, g1)
	}
	if !almostEqual(g2, 7.688098) {
		t.Fail()
		t.Logf("Expected P2 Gain to be %f, got %f\n", 7.688098, g2)
	}
}

func TestFinishedMatch(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerOneWin,
	})

	o1 := p1.elo
	o2 := p2.elo

	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerOneWin,
	})

	if o1 != p1.elo || o2 != p2.elo {
		t.Fail()
		t.Log("The same matched played twice should not alter elo twice.")
	}
}
