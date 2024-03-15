package elo

type CalculatorBuilder struct {
	c Calculator
}

type Calculator struct {
	k           float64
	deviation   float64
	scoreWeight float64
	ignoreDraws bool
	strategy    StrategyFunc
}

func NewCalculatorBuilder() *CalculatorBuilder {
	return &CalculatorBuilder{
		c: Calculator{
			k:         32,
			deviation: 400,
			strategy:  StrategyDefault,
		}}
}

// Set a strategy for calculating elo. Default is elo.StrategyDefault.
func (b *CalculatorBuilder) WithStrategy(sf StrategyFunc) *CalculatorBuilder {
	b.c.strategy = sf
	return b
}

// Set a K-Value.
// A greater K-Value means more rapid changes. Default is 32.
func (b *CalculatorBuilder) WithKValue(k float64) *CalculatorBuilder {
	b.c.k = k
	return b
}

// Set a deviation. The lower the number, the greater the probabilty that
// the higher-rated player wins (and therefore less elo gained). Default is 400.
func (b *CalculatorBuilder) WithDeviation(d float64) *CalculatorBuilder {
	b.c.deviation = d
	return b
}

// Set a score weight. The higher the number, the more the final score will influence
// the calculated elo ratings after the match. Must be greater than 0. Providing a negative
// value will result in no change to the score weight. Recommended values are between 0 and 1.
// Default is 0.
func (b *CalculatorBuilder) WithScoreWeight(w float64) *CalculatorBuilder {
	if w < 0 || w > 1 {
		return b
	}
	b.c.scoreWeight = w
	return b
}

// Set a deviation. The lower the number, the greater the probabilty that
// the higher-rated player wins (and therefore less elo gained). Default is 400.
func (b *CalculatorBuilder) WithIgnoreDraws() *CalculatorBuilder {
	b.c.ignoreDraws = true
	return b
}

// Returns a Calculator reference using the settings defined by the builder.
func (b *CalculatorBuilder) Build() *Calculator {
	return &b.c
}

// Calculate elo changes using the calculator. Returns player one and player two's new
// elo values respectively.
func (c *Calculator) Calculate(p1, p2 float64, result *MatchResult) (float64, float64) {
	if (result.Outcome == OutcomeDraw) &&
		c.ignoreDraws &&
		(result.PlayerOneScore == result.PlayerTwoScore) {
		return p1, p2
	}
	return c.strategy(&CalculatorInput{
		PlayerOne:      p1,
		PlayerTwo:      p2,
		PlayerOneScore: result.PlayerOneScore,
		PlayerTwoScore: result.PlayerTwoScore,
		Outcome:        result.Outcome,
		K:              c.k,
		Deviation:      c.deviation,
		ScoreWeight:    c.scoreWeight,
	})
}

type CalculatorInput struct {

	// Required. Elo of Player 1.
	PlayerOne float64

	// Required. Elo of Player 2.
	PlayerTwo float64

	// Required for non-scored strategies i.e. StrategyDefault.
	Outcome MatchOutcome

	// Required for scored strategies i.e. StrategyScoredDefault or StrategyScoredKValue.
	PlayerOneScore int

	// Required for scored strategies i.e. StrategyScoredDefault or StrategyScoredKValue.
	PlayerTwoScore int

	// Required. K-Value for calculating elo changes.
	// A greater K value means more rapid changes.
	K float64

	// Required. Deviation, as provided by the Calculator.
	Deviation float64

	// Required for scored strategies.
	ScoreWeight float64
}
