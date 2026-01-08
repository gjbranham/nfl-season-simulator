# nfl-season-simulator

Simulate N nfl seasons given realistic distributions of team strength

This repository contains a small NFL season simulator. It generates per-team "true strength" values from a normal distribution, simulates a regular season schedule where teams play a fixed number of games, and picks a season champion using a simple playoff placeholder.

Key concepts and parameters

- Team strength: each team receives a `Strength` value drawn from a normal distribution (mean 0, user-controlled standard deviation). Strength represents a team's underlying quality — higher values make a team more likely to win. Think of it as an aggregate skill/ability metric used by the match outcome probability function.

- Strength SD (`-sd`): controls the spread of team strengths. A parity parameter if you will. A small `sd` (close to 0) makes teams very similar, producing more even outcomes. A larger `sd` increases disparity between teams, so strong teams win more often and weak teams lose more often. Default is 0.9.

- Home field advantage (`-hfa`): a fixed additive advantage applied to the home team's strength when computing win probability. Typical values are small (e.g., 0.1–0.3) and shift the win probability slightly in favor of the home team. Default is 0.18.

- Luck factor (`-luck`): scales the logistic function used for win probability. Larger `luck` reduces the effect of strength differences on the win probability (making outcomes more random); smaller `luck` makes the function steeper so small strength differences have a bigger impact. Default is 1.0.

How a game outcome is decided

The simulator computes a win probability for the home team using a logistic function:

1 / (1 + exp(-(home_strength - away_strength + hfa)/luck))

Then a uniform random draw determines the winner based on that probability. Regular season wins/losses are recorded in each `Team` (playoffs are currently a simple placeholder that picks the highest-win team as champion).

Running the simulator

Run the CLI in the `cmd/nfl-simulator` package. Examples:

```bash
# run a single simulation and write results to `data/results/sim1.txt`
go run ./cmd/nfl-simulator -n 1

# run multiple simulations (writes sim1.txt, sim2.txt, ...)
go run ./cmd/nfl-simulator -n 10
```

Notes & next steps

- The current playoff implementation is a placeholder; you can replace `RunPlayoffs` with a realistic playoff bracket.
- Tie handling, tiebreakers, and divisional standings are minimal (or nonexistent) — improvements are welcome.
- If you want to include ties (games that end tied) or model overtime, the game resolution and record types would need small changes.

JSON output

By default the simulator now writes all simulation results into a single JSON file called `results.json` in the current working directory. Use `-f <filename>` to choose a different filename. The file contains an array of season objects with the following form:

```json
[
	{
		"season_id": 1,
		"teams": [
			{"id": 1, "name": "Team 1", "record": "12-5", "strength": 0.123},
			...
		],
		"champion_id": 5,
		"champion_name": "Baltimore Ravens",
		"champion_record": "12-5"
	}
]
```

Example: run and write to a different file

```bash
go run ./cmd/nfl-simulator -n 10 -f my_results.json
```

The JSON output is intentionally simple to make downstream analysis easy.
