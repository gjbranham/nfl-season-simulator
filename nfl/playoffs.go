package nfl

import (
	"math/rand"
	"sort"
)

// MatchResult records a single playoff matchup
type MatchResult struct {
	Home       string `json:"home"`
	HomeSeed   int    `json:"home_seed"`
	Away       string `json:"away"`
	AwaySeed   int    `json:"away_seed"`
	Winner     string `json:"winner"`
	WinnerSeed int    `json:"winner_seed"`
}

// ConfBracket holds one conference's bracket
type ConfBracket struct {
	Bye        string        `json:"bye"`
	Wildcard   []MatchResult `json:"wildcard"`
	Divisional []MatchResult `json:"divisional"`
	Conference MatchResult   `json:"conference"`
}

// PlayoffBracket is the full playoffs payload
type PlayoffBracket struct {
	NFC       ConfBracket `json:"nfc"`
	AFC       ConfBracket `json:"afc"`
	SuperBowl MatchResult `json:"superbowl"`
}

// RunPlayoffs runs the playoff bracket and returns the champion and a bracket payload
func RunPlayoffs(teams []Team, params Params) (Team, PlayoffBracket) {
	var bracket PlayoffBracket
	// split teams into conferences by ID
	var nfc, afc []Team
	for _, t := range teams {
		if t.ID >= 1 && t.ID <= 16 {
			nfc = append(nfc, t)
		} else {
			afc = append(afc, t)
		}
	}

	nfcChamp, nfcBracket, nfcSeed := runConference(nfc, params)
	afcChamp, afcBracket, afcSeed := runConference(afc, params)

	bracket.NFC = nfcBracket
	bracket.AFC = afcBracket

	// Super Bowl: neutral site (no home field advantage)
	sbHome := nfcChamp
	sbAway := afcChamp
	p := WinProb(sbHome.Strength, sbAway.Strength, 0.0, params.Luck)
	var sbWinner Team
	if rand.Float64() < p {
		sbWinner = sbHome
	} else {
		sbWinner = sbAway
	}

	// set Super Bowl seeds using conference champion seeds
	homeSeed := nfcSeed
	awaySeed := afcSeed
	winnerSeed := homeSeed
	if sbWinner.ID == sbAway.ID {
		winnerSeed = awaySeed
	}
	bracket.SuperBowl = MatchResult{Home: sbHome.Name, HomeSeed: homeSeed, Away: sbAway.Name, AwaySeed: awaySeed, Winner: sbWinner.Name, WinnerSeed: winnerSeed}
	return sbWinner, bracket
}

// runConference runs one conference playoffs and returns the champion and its bracket
func runConference(conf []Team, params Params) (Team, ConfBracket, int) {
	var bracket ConfBracket

	// sort by wins desc, then strength desc
	sort.Slice(conf, func(i, j int) bool {
		if conf[i].Wins != conf[j].Wins {
			return conf[i].Wins > conf[j].Wins
		}
		return conf[i].Strength > conf[j].Strength
	})

	// take top 7 (or fewer if less teams)
	qual := conf[:7]

	// bye
	bracket.Bye = qual[0].Name

	// map seed (1-based)
	seeds := make([]Team, len(qual)+1)
	for i, t := range qual {
		seeds[i+1] = t
	}

	type seeded struct {
		seed int
		team Team
	}

	// wildcard
	matches := [][2]int{{2, 7}, {3, 6}, {4, 5}}
	var wcWinners []seeded
	for _, m := range matches {
		a, b := m[0], m[1]
		home := seeds[a]
		away := seeds[b]
		winner := playMatch(home, away, params)
		winSeed := a
		if winner.ID == seeds[b].ID {
			winSeed = b
		}
		wcWinners = append(wcWinners, seeded{seed: winSeed, team: winner})
		bracket.Wildcard = append(bracket.Wildcard, MatchResult{Home: home.Name, HomeSeed: a, Away: away.Name, AwaySeed: b, Winner: winner.Name, WinnerSeed: winSeed})

	}

	// divisional: seed1 + wcWinners
	var divTeams []seeded
	divTeams = append(divTeams, seeded{seed: 1, team: seeds[1]})
	divTeams = append(divTeams, wcWinners...)

	var divWinners []seeded
	// find lowest seed to face seed1
	lowestSeed := divTeams[0].seed
	for i := 1; i < len(divTeams); i++ {
		if divTeams[i].seed > lowestSeed {
			lowestSeed = divTeams[i].seed
		}
	}

	var seed1 seeded
	for _, s := range divTeams {
		if s.seed == 1 {
			seed1 = s
			break
		}
	}

	var lowestTeam seeded
	for _, s := range divTeams {
		if s.seed == lowestSeed {
			lowestTeam = s
			break
		}
	}

	home := seed1
	away := lowestTeam
	if away.seed < home.seed {
		home, away = away, home
	}
	w1 := playMatch(home.team, away.team, params)
	winSeed1 := home.seed
	if w1.ID == away.team.ID {
		winSeed1 = away.seed
	}
	divWinners = append(divWinners, seeded{seed: winSeed1, team: w1})
	bracket.Divisional = append(bracket.Divisional, MatchResult{Home: home.team.Name, HomeSeed: home.seed, Away: away.team.Name, AwaySeed: away.seed, Winner: w1.Name, WinnerSeed: winSeed1})

	var otherPairs []seeded
	for _, s := range divTeams {
		if s.seed != 1 && s.seed != lowestSeed {
			otherPairs = append(otherPairs, s)
		}
	}
	if len(otherPairs) == 2 {
		home2 := otherPairs[0]
		away2 := otherPairs[1]
		if away2.seed < home2.seed {
			home2, away2 = away2, home2
		}
		w2 := playMatch(home2.team, away2.team, params)
		winSeed2 := home2.seed
		if w2.ID == away2.team.ID {
			winSeed2 = away2.seed
		}
		divWinners = append(divWinners, seeded{seed: winSeed2, team: w2})
		bracket.Divisional = append(bracket.Divisional, MatchResult{Home: home2.team.Name, HomeSeed: home2.seed, Away: away2.team.Name, AwaySeed: away2.seed, Winner: w2.Name, WinnerSeed: winSeed2})
	}

	// conference championship
	if len(divWinners) == 1 {
		bracket.Conference = MatchResult{Home: divWinners[0].team.Name, HomeSeed: divWinners[0].seed, Away: "", AwaySeed: 0, Winner: divWinners[0].team.Name, WinnerSeed: divWinners[0].seed}
		return divWinners[0].team, bracket, divWinners[0].seed
	}

	homeC := divWinners[0]
	awayC := divWinners[1]
	if awayC.seed < homeC.seed {
		homeC, awayC = awayC, homeC
	}
	confWinner := playMatch(homeC.team, awayC.team, params)
	winSeedC := homeC.seed
	if confWinner.ID == awayC.team.ID {
		winSeedC = awayC.seed
	}
	bracket.Conference = MatchResult{Home: homeC.team.Name, HomeSeed: homeC.seed, Away: awayC.team.Name, AwaySeed: awayC.seed, Winner: confWinner.Name, WinnerSeed: winSeedC}
	return confWinner, bracket, winSeedC
}

func playMatch(home, away Team, params Params) Team {
	p := WinProb(home.Strength, away.Strength, params.HomeField, params.Luck)
	if rand.Float64() < p {
		return home
	}
	return away
}
