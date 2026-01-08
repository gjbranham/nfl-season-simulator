package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gjbranham/nfl-season-simulator/nfl"
)

// WriteResultsJSON writes all simulation results into a single JSON file
func WriteResultsJSON(results []nfl.Result, filename string) error {
	type TeamOut struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Record   string  `json:"record"`
		Strength float64 `json:"strength"`
	}

	type SeasonOut struct {
		SeasonID         int       `json:"season_id"`
		Teams            []TeamOut `json:"teams"`
		ChampionID       int       `json:"champion_id"`
		ChampionName     string    `json:"champion_name"`
		ChampionRecord   string    `json:"champion_record"`
		ChampionStrength float64   `json:"champion_strength"`
		PlayoffBracket   nfl.PlayoffBracket `json:"playoff_bracket"`
	}

	out := make([]SeasonOut, len(results))
	for i, r := range results {
		teams := make([]TeamOut, len(r.Teams))
		for j, t := range r.Teams {
			teams[j] = TeamOut{ID: t.ID, Name: t.Name, Record: t.Record(), Strength: t.Strength}
		}

		out[i] = SeasonOut{
			SeasonID:         i + 1,
			ChampionID:       r.Champion.ID,
			ChampionName:     r.Champion.Name,
			ChampionRecord:   r.Champion.Record(),
			ChampionStrength: r.Champion.Strength,
			Teams:            teams,
			PlayoffBracket:   r.Bracket,
		}
	}

	// marshal with indent for readability
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	outPath := filepath.Join("results", filename)

	// ensure parent dir exists
	if dir := filepath.Dir(outPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(outPath, b, 0o644)
}
