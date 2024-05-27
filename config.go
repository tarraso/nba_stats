package main

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	RedisHost  string
	RedisPort  string
}

func LoadConfig() Config {
	return Config{
		DBUser:     getEnv("DB_USER", "your_default_db_user"),
		DBPassword: getEnv("DB_PASSWORD", "your_default_db_password"),
		DBName:     getEnv("DB_NAME", "your_default_db_name"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
