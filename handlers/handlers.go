package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"nba_stats/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
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
func AddPlayerHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
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
func AddStatHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
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

		//cache invalidation
		cacheKeyPlayer := fmt.Sprintf("player_stats_%d", stat.PlayerID)
		rdb.Del(cacheKeyPlayer)
		query = `SELECT team_id from players where id=$1`
		var teamID int
		err = db.QueryRow(query, stat.PlayerID).Scan(&teamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cacheKeyTeam := fmt.Sprintf("team_stats_%d", stat.PlayerID)
		rdb.Del(cacheKeyTeam)
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
func ListPlayersHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
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
// @Param playerId path int true "PlayerId"
// @Success 200 {array} models.AvgStat
// @Failure 500 {string} string "Internal server error"
// @Router /stat/players/{playerId} [get]
func GetPlayerAvgStatHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		playerID, err := strconv.Atoi(vars["playerId"])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		cacheKey := fmt.Sprintf("player_stats_%d", playerID)

		// Try to get cached data
		cachedData, err := rdb.Get(cacheKey).Result()
		if err == redis.Nil {
			// Cache miss, fetch data from DB
			stats, err := getAvgPlayerStats(db, playerID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Player not found", http.StatusNotFound)
				} else {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}

			data, err := json.Marshal(stats)
			if err != nil {
				http.Error(w, "json marshal error", http.StatusInternalServerError)
				return
			}

			//Make data expire after 24 hour
			err = rdb.Set(cacheKey, data, 24*time.Hour).Err()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(stats)
		} else if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			// Cache hit
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedData))
		}
	}
}

// PlayerStatHandler godoc
// @Summary team stats
// @Description Get a list of all players
// @Tags players
// @Produce json
// @Param teamId path int true "teamId"
// @Success 200 {array} models.AvgStat
// @Failure 500 {string} string "Internal server error"
// @Router /stat/teams/{teamId} [get]
func GetTeamAvgStatHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		teamID, err := strconv.Atoi(vars["teamId"])
		if err != nil {
			http.Error(w, "Invalid team ID", http.StatusBadRequest)
			return
		}

		cacheKey := fmt.Sprintf("team_stats_%d", teamID)

		// Try to get cached data
		cachedData, err := rdb.Get(cacheKey).Result()
		if err == redis.Nil {
			// Cache miss, fetch data from DB
			stats, err := getAvgTeamStats(db, teamID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Player not found", http.StatusNotFound)
				} else {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}

			data, err := json.Marshal(stats)
			if err != nil {
				http.Error(w, "json marshal error", http.StatusInternalServerError)
				return
			}

			//Make data expire after 24 hour
			err = rdb.Set(cacheKey, data, 24*time.Hour).Err()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(stats)
		} else if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			// Cache hit
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedData))
		}

	}
}

func getAvgPlayerStats(db *sql.DB, playerID int) (*models.AvgStat, error) {
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
		return nil, err
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
	if rows.Next() {
		err := rows.Scan(&avg_points, &avg_rebounds, &avg_assists, &avg_steals, &avg_blocks,
			&avg_fouls, &avg_turnovers, &avg_minutes_played)
		if err != nil {
			return nil, err
		}
		stat := models.AvgStat{
			AvgPoints:        avg_points,
			AvgRebounds:      avg_rebounds,
			AvgAssists:       avg_assists,
			AvgSteals:        avg_steals,
			AvgBlocks:        avg_blocks,
			AvgFouls:         avg_fouls,
			AvgTurnovers:     avg_turnovers,
			AvgMinutesPlayed: avg_minutes_played,
		}
		return &stat, nil
	} else {
		return nil, errors.New("no rows found")
	}
}

func getAvgTeamStats(db *sql.DB, teamID int) (*models.AvgStat, error) {
	query := `
			SELECT
			AVG(stats.points) AS avg_points,
			AVG(stats.rebounds) AS avg_rebounds,
			AVG(stats.assists) AS avg_assists,
			AVG(stats.steals) AS avg_steals,
			AVG(stats.blocks) AS avg_blocks,
			AVG(stats.fouls) AS avg_fouls,
			AVG(stats.turnovers) AS avg_turnovers,
			AVG(stats.minutes_played) AS avg_minutes_played
		FROM
			stats
		JOIN
			players ON players.id = stats.player_id
		WHERE
			players.team_id = $1;
	`
	var stats models.AvgStat
	err := db.QueryRow(query, teamID).Scan(
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
			return nil, err
		} else {
			return nil, errors.New("database error")
		}
	}
	return &stats, nil
}
