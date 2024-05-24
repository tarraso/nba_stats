package handlers

import (
	"database/sql"
	"encoding/json"
	"nba_stats/models"
	"net/http"
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
		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO players (name, team) VALUES ($1, $2) RETURNING id`
		err = db.QueryRow(query, player.Name, player.Team).Scan(&player.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
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
		var players []models.Player

		rows, err := db.Query("SELECT id, name, team FROM players")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var player models.Player
			err := rows.Scan(&player.ID, &player.Name, &player.Team)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			players = append(players, player)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(players)
	}
}
