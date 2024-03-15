package elo

import (
	"math"
)

type StrategyFunc func(input *CalculatorInput) (r1 float64, r2 float64)

// Calculates elo based on a Win/Loss system.
func StrategyDefault(input *CalculatorInput) (float64, float64) {

	R1 := math.Pow(10, input.PlayerOne/input.Deviation)
	R2 := math.Pow(10, input.PlayerTwo/input.Deviation)

	E1 := R1 / (R1 + R2)
	E2 := R2 / (R1 + R2)

	var S1, S2 float64
	switch input.Outcome {
	case OutcomePlayerOneWin:
		S1 = 1.0
		S2 = 0.0
	case OutcomePlayerTwoWin:
		S1 = 0.0
		S2 = 1.0
	case OutcomeDraw:
		S1 = 0.5
		S2 = 0.5
	}

	NewP1 := input.PlayerOne + input.K*(S1-E1)
	NewP2 := input.PlayerTwo + input.K*(S2-E2)
	return NewP1, NewP2
}

func determineWeightedSValues(winnerScore, loserScore, winnerE, loserE, weight float64) (S1, S2 float64) {
	// determine the domination factor
	D := winnerScore / (winnerScore + loserScore)
	// calculate the amount of elo to be gained or lost, based on the domination factor and the expected chance
	// for the winner of the match to win
	// in general, G is smaller if the winner was heavily favored, and larger if the winner was not favored
	G := (1.0 - winnerE) * D
	// apply the score weight multiplier
	G = G * math.Exp((-1*weight)*winnerE)
	// add the weight to the expected chance to win to get a value between winnerE and 1.0 that
	// is weighted based on the domination factor, and will be multiplied by K to get actual elo gained
	S1 = G + winnerE
	// and subtract G from E2 to get a value between 0 and loserE that will be multiplied by K to get actual elo lost
	S2 = loserE - G

	return S1, S2
}

// Calculates elo weighted by the final score.
// A more dominant score means greater elo gained.
func StrategyScored(input *CalculatorInput) (float64, float64) {

	R1 := math.Pow(10, input.PlayerOne/input.Deviation)
	R2 := math.Pow(10, input.PlayerTwo/input.Deviation)

	E1 := R1 / (R1 + R2)
	E2 := R2 / (R1 + R2)

	var S1, S2 float64
	if input.PlayerOneScore == 0 {
		S1 = 0
		S2 = 1
	} else if input.PlayerTwoScore == 0 {
		S1 = 1
		S2 = 0
	} else if input.PlayerOneScore > input.PlayerTwoScore { // P1 win
		S1, S2 = determineWeightedSValues(
			float64(input.PlayerOneScore),
			float64(input.PlayerTwoScore),
			E1,
			E2,
			input.ScoreWeight,
		)
	} else if input.PlayerOneScore < input.PlayerTwoScore { // P2 win
		S2, S1 = determineWeightedSValues(
			float64(input.PlayerTwoScore),
			float64(input.PlayerOneScore),
			E2,
			E1,
			input.ScoreWeight,
		)
	} else {
		S1 = 0.5
		S2 = 0.5
	}

	NewP1 := input.PlayerOne + input.K*(S1-E1)
	NewP2 := input.PlayerTwo + input.K*(S2-E2)
	return NewP1, NewP2
}
