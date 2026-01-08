package nfl

import "math/rand"

func GenerateSchedule(numTeams int) []Game {
	games := []Game{}
	opponents := make([][]int, numTeams)

	for i := range numTeams {
		opponents[i] = rand.Perm(numTeams)
	}

	gamesPerTeam := 17
	counts := make([]int, numTeams)

	for i := range numTeams {
		for _, j := range opponents[i] {
			if i == j {
				continue
			}
			if counts[i] >= gamesPerTeam || counts[j] >= gamesPerTeam {
				continue
			}

			games = append(games, Game{Home: i, Away: j})
			counts[i]++
			counts[j]++

			if counts[i] >= gamesPerTeam {
				break
			}
		}
	}

	return games
}
