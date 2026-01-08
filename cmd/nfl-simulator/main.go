package main

import (
	"flag"
	"fmt"

	"github.com/gjbranham/nfl-season-simulator/nfl"
)

func main() {
	n := flag.Int("n", 1, "number of seasons")
	outFile := flag.String("f", "results.json", "output JSON filename (will be overwritten)")
	sd := flag.Float64("sd", 0.9, "team strength SD")
	hfa := flag.Float64("hfa", 0.18, "home field advantage")
	luck := flag.Float64("luck", 1.0, "luck factor")

	flag.Parse()

	teams := []nfl.Team{
		{ID: 1, Name: "Arizona Cardinals", Conference: "NFC"},
		{ID: 2, Name: "Los Angeles Rams", Conference: "NFC"},
		{ID: 3, Name: "San Francisco 49ers", Conference: "NFC"},
		{ID: 4, Name: "Seattle Seahawks", Conference: "NFC"},
		{ID: 5, Name: "Chicago Bears", Conference: "NFC"},
		{ID: 6, Name: "Detroit Lions", Conference: "NFC"},
		{ID: 7, Name: "Green Bay Packers", Conference: "NFC"},
		{ID: 8, Name: "Minnesota Vikings", Conference: "NFC"},
		{ID: 9, Name: "Atlanta Falcons", Conference: "NFC"},
		{ID: 10, Name: "Carolina Panthers", Conference: "NFC"},
		{ID: 11, Name: "New Orleans Saints", Conference: "NFC"},
		{ID: 12, Name: "Tampa Bay Buccaneers", Conference: "NFC"},
		{ID: 13, Name: "Dallas Cowboys", Conference: "NFC"},
		{ID: 14, Name: "New York Giants", Conference: "NFC"},
		{ID: 15, Name: "Philadelphia Eagles", Conference: "NFC"},
		{ID: 16, Name: "Washington Commanders", Conference: "NFC"},
		{ID: 17, Name: "Denver Broncos", Conference: "AFC"},
		{ID: 18, Name: "Kansas City Chiefs", Conference: "AFC"},
		{ID: 19, Name: "Las Vegas Raiders", Conference: "AFC"},
		{ID: 20, Name: "Los Angeles Chargers", Conference: "AFC"},
		{ID: 21, Name: "Baltimore Ravens", Conference: "AFC"},
		{ID: 22, Name: "Cleveland Browns", Conference: "AFC"},
		{ID: 23, Name: "Pittsburgh Steelers", Conference: "AFC"},
		{ID: 24, Name: "Cincinnati Bengals", Conference: "AFC"},
		{ID: 25, Name: "Houston Texans", Conference: "AFC"},
		{ID: 26, Name: "Indianapolis Colts", Conference: "AFC"},
		{ID: 27, Name: "Jacksonville Jaguars", Conference: "AFC"},
		{ID: 28, Name: "Tennessee Titans", Conference: "AFC"},
		{ID: 29, Name: "Buffalo Bills", Conference: "AFC"},
		{ID: 30, Name: "Miami Dolphins", Conference: "AFC"},
		{ID: 31, Name: "New England Patriots", Conference: "AFC"},
		{ID: 32, Name: "New York Jets", Conference: "AFC"},
	}

	params := nfl.Params{
		StrengthSD: *sd,
		HomeField:  *hfa,
		Luck:       *luck,
	}

	// Run all simulations and write JSON output
	results := nfl.Run(*n, teams, params)
	if err := WriteResultsJSON(results, *outFile); err != nil {
		panic(fmt.Errorf("failed to write results JSON: %w", err))
	}

}
