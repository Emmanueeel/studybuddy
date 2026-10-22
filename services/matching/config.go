package main

import "os"

type Config struct {
	Port       string
	ProfileURL string
	RedisAddr  string
}

func LoadConfig() Config {
	return Config{
		Port:       getEnv("MATCHING_PORT", "8083"),
		ProfileURL: getEnv("PROFILE_URL", "http://localhost:8082"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
