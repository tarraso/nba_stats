package models

import (
	"time"
)

// Player represents a basketball player
type Player struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	TeamID int    `json:"team_id"` // New field for foreign key
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

// AvgStat
type AvgStat struct {
	AvgPoints        float64   `json:"avg_points"`
	AvgRebounds      float64   `json:"avg_rebounds"`
	AvgAssists       float64   `json:"avg_assists"`
	AvgSteals        float64   `json:"avg_steals"`
	AvgBlocks        float64   `json:"avg_blocks"`
	AvgFouls         float64   `json:"avg_fouls"`
	AvgTurnovers     float64   `json:"avg_turnovers"`
	AvgMinutesPlayed float64   `json:"avg_minutes_played"`
	AvgGameDate      time.Time `json:"avg_game_date"`
}
