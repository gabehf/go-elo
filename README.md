# go-elo: A flexible elo calculator for Go

![coverage](https://raw.githubusercontent.com/gabehf/go-elo/badges/.badges/main/coverage.svg)

This package allows you to add flexible elo/skill tracking to your application.

Go-elo supports unscored (win/loss) and scored matches, with options to tailor the elo curve to your liking.

## Usage

Install:

```bash
go get github.com/gabehf/go-elo
```

Simple unscored example:

```go
func main() {
    c := elo.NewCalculator().Build()
    
    X, Y := c.Calculate(1200.0, 1100.0, &elo.MatchResult{
        Outcome: elo.OutcomePlayerOneWin,
    })
    // X and Y are player one and player two's new elo values, respectively.
}
```

Simple scored example:

```go
func main() {
    c := elo.NewCalculator().
        WithStrategy(elo.StrategyScored).
        Build()
    
    X, Y := c.Calculate(1200.0, 1100.0, &elo.MatchResult{
        PlayerOneScore: 12,
        PlayerTwoScore: 9,
    })
    // X and Y are player one and player two's new elo values, respectively.
}
```

Using the match system:

```go
func main() {
    c := elo.NewCalculatorBuilder().Build()

    p1 := MyPlayer{Elo: 1200}
    p2 := MyPlayer{Elo: 1100}

    m := c.NewMatch(p1, p2)

    m.GetOdds() // returns player one and two's probability of winning
    m.PlayerOneGain() // returns how much elo player one stands to gain
    m.PlayerTwoGain() // returns how much elo player two stands to gain
    m.SetStrategy(func(input *CalculatorInput) (r1 float64, r2 float64) {
        // custom elo calculation method
        // used for only this match
    })

    m.Play(&elo.MatchResult{
        Outcome: elo.OutcomePlayerTwoWin,
    })

    p1.GetElo() // player one's new elo
    p2.GetElo() // player two's new elo
}
```

## Adjusting Parameters

Go-elo has parameters that you can customize in order to fine tune the elo curve you are looking for. To learn about exactly how each of the parameters are used during the elo calculation, refer to the [ELO.md file](ELO.md) in this repo.

### At the calculator level

You can adjust parameters at the calculator level, when you want to reuse the settings over many matches. This is accomplished using methods with the calculator buildier.

```go
func main() {
    c := elo.NewCalculatorBuilder().
        WithKValue(32). // specify the K-Factor
        WithScoreWeight(0.5). // specify the score weight
        WithDeviation(200). // specify the deviation
        WithStrategy(elo.StrategyScored).
        WithIgnoreDraws(). // specify not to adjust elo after a draw
        Build()
    
    // now you can use your adjusted elo calculator
    m := c.NewMatch(...)
}
```

### At the match level

You can also adjust the parameters at the per-match level. Here is an example:

```go
func main() {
    c := elo.newCalculatorBuilder().Build()

    p1 := MyPlayer{Elo: 1500}
    p2 := MyPlayer{Elo: 1600}

    m := elo.NewMatch(p1, p2)

    // adjust parameters
    m.SetKValue(54)
    m.SetScoreWeight(0.33)
    m.SetDeviation(250)
    m.IgnoreDraws(true)

    m.Play(...)
}
```
