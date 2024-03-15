package elo_test

import (
	"math"
	"testing"

	"github.com/gabehf/go-elo"
)

// Equal precision to Go's default precision when printing with %f.
const float64EqualityThreshold = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

type player struct {
	elo float64
}

func (p *player) GetElo() float64 {
	return p.elo
}
func (p *player) SetElo(e float64) {
	p.elo = e
}

func TestDefaultCalculator(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.SetStrategy(elo.StrategyDefault)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerTwoWin,
	})

	if !almostEqual(p1.elo, 1592.311901653) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1592.311901653, p1.elo)
	}
	if !almostEqual(p2.elo, 1807.688098347) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1807.688098347, p2.elo)
	}

	p1.elo = 800
	p2.elo = 1300

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerOneWin,
	})

	if !almostEqual(p1.elo, 830.296313114) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 830.296313114, p1.elo)
	}
	if !almostEqual(p2.elo, 1269.703686886) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1269.703686886, p2.elo)
	}
}

func TestScoredStrategy(t *testing.T) {
	c := elo.NewCalculatorBuilder().
		WithStrategy(elo.StrategyScored).
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1600
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 150,
		PlayerTwoScore: 160,
	})

	if !almostEqual(p1.elo, 1596.031949) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1596.031949, p1.elo)
	}
	if !almostEqual(p2.elo, 1803.968051) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1803.968051, p2.elo)
	}

	p1.elo = 800
	p2.elo = 1300

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 15,
		PlayerTwoScore: 20,
	})

	if !almostEqual(p1.elo, 799.026465) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 799.026465, p1.elo)
	}
	if !almostEqual(p2.elo, 1300.973535) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1300.973535, p2.elo)
	}

	p1.elo = 900
	p2.elo = 700

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 6,
		PlayerTwoScore: 0,
	})

	if !almostEqual(p1.elo, 907.688098) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 907.688098, p1.elo)
	}
	if !almostEqual(p2.elo, 692.311902) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 692.311902, p2.elo)
	}
}

func TestWithKValue(t *testing.T) {
	c := elo.NewCalculatorBuilder().WithKValue(55).Build()

	p1 := new(player)
	p1.elo = 1600
	p2 := new(player)
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerTwoWin,
	})

	if !almostEqual(p1.elo, 1586.786081) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1586.786081, p1.elo)
	}
	if !almostEqual(p2.elo, 1813.213919) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1813.213919, p2.elo)
	}

	p1.elo = 800
	p2.elo = 1300

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerOneWin,
	})

	if !almostEqual(p1.elo, 852.071788) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 852.071788, p1.elo)
	}
	if !almostEqual(p2.elo, 1247.928212) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1247.928212, p2.elo)
	}
}

func TestScoredWithKValue(t *testing.T) {
	c := elo.NewCalculatorBuilder().
		WithStrategy(elo.StrategyScored).
		WithKValue(60).
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1600
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 230,
		PlayerTwoScore: 160,
	})

	if !almostEqual(p1.elo, 1626.883353) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1626.883353, p1.elo)
	}
	if !almostEqual(p2.elo, 1773.116647) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1773.116647, p2.elo)
	}

	p1.elo = 800
	p2.elo = 1300

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 0,
		PlayerTwoScore: 4,
	})

	if !almostEqual(p1.elo, 796.805587) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 796.805587, p1.elo)
	}
	if !almostEqual(p2.elo, 1303.194413) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1303.194413, p2.elo)
	}
}

func TestScoredWithCloseScore(t *testing.T) {
	c := elo.NewCalculatorBuilder().
		WithStrategy(elo.StrategyScored).
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1200
	p2.elo = 1100

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 501,
		PlayerTwoScore: 500,
	})

	if p1.elo < 1200 {
		t.Fail()
		t.Log("Winning player's elo must always be higher than starting elo.")
	}
	if p2.elo > 1100 {
		t.Fail()
		t.Log("Losing player's elo must always be lower than starting elo.")
	}

	if !almostEqual(p1.elo, 1205.764713) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1205.764713, p1.elo)
	}
	if !almostEqual(p2.elo, 1094.235287) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1094.235287, p2.elo)
	}
}

func TestWithDeviation(t *testing.T) {
	c := elo.NewCalculatorBuilder().
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1600
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerTwoWin,
	})

	o1 := p1.elo
	o2 := p2.elo

	c = elo.NewCalculatorBuilder().
		WithDeviation(200).
		Build()

	p1.elo = 1600
	p2.elo = 1800

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomePlayerTwoWin,
	})

	if o1 > p1.elo || o2 < p2.elo {
		t.Fail()
		t.Log("")
		t.Logf("Tighter deviation should not result in more change to elo."+
			"\nD=400 P1: %f, D=200 P1: %f\n"+
			"D=400 P2: %f, D=200 P2: %f", o1, p1.elo, o2, p2.elo)
	}
	if !almostEqual(p1.elo, 1597.090909) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1597.090909, p1.elo)
	}
	if !almostEqual(p2.elo, 1802.909091) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1802.909091, p2.elo)
	}
}

func TestDraw(t *testing.T) {

	// w/l

	c := elo.NewCalculatorBuilder().
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1600
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomeDraw,
	})
	if !almostEqual(p1.elo, 1608.311902) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1608.311902, p1.elo)
	}
	if !almostEqual(p2.elo, 1791.688098) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1791.688098, p2.elo)
	}

	// scored

	c = elo.NewCalculatorBuilder().
		WithStrategy(elo.StrategyScored).
		Build()

	p1.elo = 1600
	p2.elo = 1800

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 500,
		PlayerTwoScore: 500,
	})

	if !almostEqual(p1.elo, 1608.311902) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1608.311902, p1.elo)
	}
	if !almostEqual(p2.elo, 1791.688098) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1791.688098, p2.elo)
	}
}

func TestIgnoreDraw(t *testing.T) {

	// w/l
	c := elo.NewCalculatorBuilder().
		WithIgnoreDraws().
		Build()

	p1 := new(player)
	p2 := new(player)
	p1.elo = 1600
	p2.elo = 1800

	m := c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		Outcome: elo.OutcomeDraw,
	})

	if !almostEqual(p1.elo, 1600) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1600.0, p1.elo)
	}
	if !almostEqual(p2.elo, 1800) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1600.0, p2.elo)
	}

	// scored
	c = elo.NewCalculatorBuilder().
		WithIgnoreDraws().
		WithStrategy(elo.StrategyScored).
		Build()

	p1.elo = 1600
	p2.elo = 1800

	m = c.NewMatch(p1, p2)
	m.Play(&elo.MatchResult{
		PlayerOneScore: 500,
		PlayerTwoScore: 500,
	})

	if !almostEqual(p1.elo, 1600) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1600.0, p1.elo)
	}
	if !almostEqual(p2.elo, 1800) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 1600.0, p2.elo)
	}
}
