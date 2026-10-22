package main

import "os"

type Config struct {
	DBUrl string
	Port  string
}

func LoadConfig() Config {
	return Config{
		DBUrl: getEnv("PROFILE_DB_URL", "postgres://studybuddy:secret@localhost:5434/profile?sslmode=disable"),
		Port:  getEnv("PROFILE_PORT", "8082"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
