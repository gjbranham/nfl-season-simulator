package nfl

import "math/rand"

type Season struct {
	Teams    []Team
	Schedule []Game
}

func NewSeason(baseTeams []Team, params Params) Season {
	teams := make([]Team, len(baseTeams))
	copy(teams, baseTeams)

	for i := range teams {
		teams[i].Strength = rand.NormFloat64() * params.StrengthSD
		teams[i].Wins = 0
		teams[i].Losses = 0
	}

	sched := GenerateSchedule(len(teams))

	return Season{
		Teams:    teams,
		Schedule: sched,
	}
}

func (s *Season) Play(params Params) {
	for _, g := range s.Schedule {
		home := &s.Teams[g.Home]
		away := &s.Teams[g.Away]

		p := WinProb(home.Strength, away.Strength, params.HomeField, params.Luck)
		if rand.Float64() < p {
			home.Wins++
			away.Losses++
		} else {
			away.Wins++
			home.Losses++
		}
	}
}
