package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// In-memory storage for simplicity
var players = []Player{
	{ID: 1, Name: "LeBron James", Team: "Lakers"},
	{ID: 2, Name: "Stephen Curry", Team: "Warriors"},
}

var gameStats = []GameStat{}

// Handler to log game statistics
func gameStatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newStat GameStat
	err := json.NewDecoder(r.Body).Decode(&newStat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ValidateGameStat(newStat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gameStats = append(gameStats, newStat)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Game stat logged successfully!")
}

// Handler to list players
func listPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

// Handler to list game statistics
func listGameStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameStats)
}

func main() {
	http.HandleFunc("/game-stat", gameStatHandler)
	http.HandleFunc("/players", listPlayersHandler)
	http.HandleFunc("/game-stats", listGameStatsHandler)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
