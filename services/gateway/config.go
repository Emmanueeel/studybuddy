package main

import "os"

type Config struct {
	Port        string
	JWTSecret   string
	AuthURL     string
	ProfileURL  string
	MatchingURL string
	RequestURL  string
	ReviewURL   string
	RedisAddr   string
}

func LoadConfig() Config {
	return Config{
		Port:        getEnv("GATEWAY_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		AuthURL:     getEnv("AUTH_URL", "http://localhost:8081"),
		ProfileURL:  getEnv("PROFILE_URL", "http://localhost:8082"),
		MatchingURL: getEnv("MATCHING_URL", "http://localhost:8083"),
		RequestURL:  getEnv("REQUEST_URL", "http://localhost:8084"),
		ReviewURL:   getEnv("REVIEW_URL", "http://localhost:8085"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
