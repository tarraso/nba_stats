package main

import "time"

// Player represents a basketball player
type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Team string `json:"team"`
}

// GameStat represents the statistics of a player in a game
type GameStat struct {
	PlayerID      int       `json:"player_id"`
	Points        int       `json:"points"`
	Rebounds      int       `json:"rebounds"`
	Assists       int       `json:"assists"`
	Steals        int       `json:"steals"`
	Blocks        int       `json:"blocks"`
	Fouls         int       `json:"fouls"`
	Turnovers     int       `json:"turnovers"`
	MinutesPlayed float64   `json:"minutes_played"`
	GameDate      time.Time `json:"game_date"`
}
