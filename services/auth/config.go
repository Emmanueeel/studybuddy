package main

import "os"

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

func LoadConfig() Config {
	return Config{
		DBUrl:     getEnv("AUTH_DB_URL", "postgres://studybuddy:secret@localhost:5433/auth?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
		Port:      getEnv("AUTH_PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
