package nfl

import "math"

type Game struct {
	Home int
	Away int
}

func WinProb(home, away, hfa, luck float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(home-away+hfa)/luck))
}
