package main

import "os"

type Config struct {
	DBUrl     string
	Port      string
	RabbitURL string
}

func LoadConfig() Config {
	return Config{
		DBUrl:     getEnv("REQUEST_DB_URL", "postgres://studybuddy:secret@localhost:5435/request?sslmode=disable"),
		Port:      getEnv("REQUEST_PORT", "8084"),
		RabbitURL: getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
