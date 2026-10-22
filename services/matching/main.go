package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := LoadConfig()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("warning: redis unavailable: %v", err)
	}

	h := &MatchingHandler{
		Client:  NewProfileClient(cfg.ProfileURL),
		Redis:   rdb,
		Workers: 10,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /find", h.Find)

	log.Printf("matching service on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
