package main

import "errors"

// Validation function for GameStat
func ValidateGameStat(gs GameStat) error {
	if gs.Points < 0 || gs.Rebounds < 0 || gs.Assists < 0 || gs.Steals < 0 || gs.Blocks < 0 || gs.Turnovers < 0 {
		return errors.New("points, rebounds, assists, steals, blocks, and turnovers must be positive integers")
	}
	if gs.Fouls < 0 || gs.Fouls > 6 {
		return errors.New("fouls must be between 0 and 6")
	}
	if gs.MinutesPlayed < 0 || gs.MinutesPlayed > 48 {
		return errors.New("minutes played must be between 0 and 48")
	}
	return nil
}
