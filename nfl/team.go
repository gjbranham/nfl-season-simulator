package nfl

import "fmt"

type Team struct {
	ID         int
	Name       string
	Division   string
	Conference string
	Strength   float64
	Wins       int
	Losses     int
}

func (t Team) Record() string {
	return fmt.Sprintf("%d-%d", t.Wins, t.Losses)
}
