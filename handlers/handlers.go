package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nba_stats/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AddPlayerHandler godoc
// @Summary Add a new player
// @Description Add a new player to the database
// @Tags players
// @Accept json
// @Produce json
// @Param player body models.Player true "Player"
// @Success 201 {object} models.Player
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /add-players [post]
func AddPlayerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var player models.Player
		if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO players (name, team_id) VALUES ($1, $2) RETURNING id`
		err := db.QueryRow(query, player.Name, player.TeamID).Scan(&player.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)
	}
}

// AddStatHandler godoc
// @Summary Add a new game stat
// @Description Add a new game stat to the database
// @Tags stats
// @Accept json
// @Produce json
// @Param stat body models.GameStat true "Game Stat"
// @Success 201 {object} models.GameStat
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /add-stat [post]
func AddStatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stat models.GameStat
		err := json.NewDecoder(r.Body).Decode(&stat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO stats (player_id, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played, game_date)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err = db.Exec(query, stat.PlayerID, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Fouls, stat.Turnovers, stat.MinutesPlayed, stat.GameDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// listPlayersHandler godoc
// @Summary List all players
// @Description Get a list of all players
// @Tags players
// @Produce json
// @Success 200 {array} models.Player
// @Failure 500 {string} string "Internal server error"
// @Router /players [get]
func ListPlayersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, name, team_id FROM players`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var players []models.Player
		for rows.Next() {
			var player models.Player
			if err := rows.Scan(&player.ID, &player.Name, &player.TeamID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			players = append(players, player)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(players)
	}
}

// PlayerStatHandler godoc
// @Summary player stats
// @Description Get a list of all players
// @Tags players
// @Produce json
// @Success 200 {array} models.AvgStat
// @Failure 500 {string} string "Internal server error"
// @Router /stats/player/:id [get]
func GetPlayerAvgStatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		playerID, err := strconv.Atoi(vars["playerId"])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Fprintf(w, "Player ID: %s\n", playerID)
		query := fmt.Sprintf(`
SELECT
	AVG(points) AS avg_points,
	AVG(rebounds) AS avg_rebounds,
	AVG(assists) AS avg_assists,
	AVG(steals) AS avg_steals,
	AVG(blocks) AS avg_blocks,
	AVG(fouls) AS avg_fouls,
	AVG(turnovers) AS avg_turnovers,
	AVG(minutes_played) AS avg_minutes_played
FROM
	stats
WHERE
	player_id = %d;`, playerID)

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var avg_points float64
		var avg_rebounds float64
		var avg_assists float64
		var avg_steals float64
		var avg_blocks float64
		var avg_fouls float64
		var avg_turnovers float64
		var avg_minutes_played float64
		var stat models.AvgStat
		if rows.Next() {
			err := rows.Scan(&avg_points, &avg_rebounds, &avg_assists, &avg_steals, &avg_blocks,
				&avg_fouls, &avg_turnovers, &avg_minutes_played)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			stat = models.AvgStat{
				PlayerID:         playerID,
				AvgPoints:        avg_points,
				AvgRebounds:      avg_rebounds,
				AvgAssists:       avg_assists,
				AvgSteals:        avg_steals,
				AvgBlocks:        avg_blocks,
				AvgFouls:         avg_fouls,
				AvgTurnovers:     avg_turnovers,
				AvgMinutesPlayed: avg_minutes_played,
			}
			// Now you can use the 'name' and 'age' variables
		} else {
			http.Error(w, "No rows found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stat)
	}
}

// PlayerStatHandler godoc
// @Summary player stats
// @Description Get a list of all players
// @Tags players
// @Produce json
// @Success 200 {array} models.AvgStat
// @Failure 500 {string} string "Internal server error"
// @Router /stats/team/:id [get]
func GetAvgStatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		teamID, err := strconv.Atoi(vars["teamId"])
		if err != nil {
			http.Error(w, "Invalid team ID", http.StatusBadRequest)
			return
		}

		query := `
            SELECT
                AVG(points) as avg_points,
                AVG(rebounds) as avg_rebounds,
                AVG(assists) as avg_assists,
                AVG(steals) as avg_steals,
                AVG(blocks) as avg_blocks,
                AVG(fouls) as avg_fouls,
                AVG(turnovers) as avg_turnovers,
                AVG(minutes) as avg_minutes
            FROM players
            WHERE team_id = $1
        `
		var stats models.AvgStat
		err = db.QueryRow(query, teamID).Scan(
			&stats.AvgPoints,
			&stats.AvgRebounds,
			&stats.AvgAssists,
			&stats.AvgSteals,
			&stats.AvgBlocks,
			&stats.AvgFouls,
			&stats.AvgTurnovers,
			&stats.AvgMinutesPlayed,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No players found for this team", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
