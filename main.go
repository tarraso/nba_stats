package main

import (
	"database/sql"
	"fmt"
	_ "nba_stats/docs" // docs is generated by Swag CLI, you have to import it.
	"nba_stats/handlers"

	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := Config{
		DBUser:     getEnv("DB_USER", "your_default_db_user"),
		DBPassword: getEnv("DB_PASSWORD", "your_default_db_password"),
		DBName:     getEnv("DB_NAME", "your_default_db_name"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v\n", err)
	}
	defer db.Close()

	runMigrations(db, cfg)

	r := mux.NewRouter()

	r.HandleFunc("/add-stat", handlers.AddStatHandler(db))
	r.HandleFunc("/stat/player/{playerId}", handlers.GetPlayerAvgStatHandler(db))
	r.HandleFunc("/add-players", handlers.AddPlayerHandler(db)) // POST /players
	r.HandleFunc("/players", handlers.ListPlayersHandler(db))   // GET /players

	// Swagger endpoint
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// initDB initializes the database connection
func initDB(cfg Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB, cfg Config) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBName, driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v\n", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v\n", err)
	}

	log.Println("Migrations applied successfully!")
}
