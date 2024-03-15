package elo_test

import (
	"math"
	"testing"

	"github.com/gabehf/go-elo"
)

func TestCalculate(t *testing.T) {
	c := elo.NewCalculatorBuilder().Build()

	p1, p2 := 1200.0, 1000.0

	n1, n2 := c.Calculate(p1, p2, &elo.MatchResult{
		Outcome: elo.OutcomePlayerOneWin,
	})

	if !almostEqual(math.Abs(1200-n1), math.Abs(1000-n2)) {
		t.Fail()
		t.Log("Elo gained and lost must be equal")
	}

	if !almostEqual(n1, 1207.688098) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1207.688098, n1)
	}
	if !almostEqual(n2, 992.311902) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 992.311902, n2)
	}
}

func TestCalculateScored(t *testing.T) {
	c := elo.NewCalculatorBuilder().
		WithStrategy(elo.StrategyScored).
		WithScoreWeight(0.33).
		WithScoreWeight(-2). // will be ignored
		Build()

	p1, p2 := 1200.0, 1000.0

	n1, n2 := c.Calculate(p1, p2, &elo.MatchResult{
		PlayerOneScore: 12,
		PlayerTwoScore: 8,
	})

	if !almostEqual(math.Abs(1200-n1), math.Abs(1000-n2)) {
		t.Fail()
		t.Log("Elo gained and lost must be equal")
	}

	if !almostEqual(n1, 1203.589925) {
		t.Fail()
		t.Logf("Expected P1 Elo %f, got %f\n", 1203.589925, n1)
	}
	if !almostEqual(n2, 996.410075) {
		t.Fail()
		t.Logf("Expected P2 Elo %f, got %f\n", 996.410075, n2)
	}

	// test ignore draw

	c = elo.NewCalculatorBuilder().
		WithIgnoreDraws().
		Build()

	p1, p2 = 1200.0, 1000.0

	n1, n2 = c.Calculate(p1, p2, &elo.MatchResult{
		Outcome: elo.OutcomeDraw,
	})

	if !almostEqual(n1, 1200) {
		t.Fail()
		t.Logf("Draw not ignored. Expected P1 Elo %f, got %f\n", 1200.0, n1)
	}
	if !almostEqual(n2, 1000) {
		t.Fail()
		t.Logf("Draw not ignored. Expected P2 Elo %f, got %f\n", 1000.0, n2)
	}
}
