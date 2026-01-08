package nfl

type Result struct {
	Champion Team
	Teams    []Team
	Bracket  PlayoffBracket
}

func Run(n int, baseTeams []Team, params Params) []Result {
	results := make([]Result, n)

	for i := range n {
		season := NewSeason(baseTeams, params)
		season.Play(params)
		champ, bracket := RunPlayoffs(season.Teams, params)

		results[i] = Result{Champion: champ, Teams: season.Teams, Bracket: bracket}
	}

	return results
}
