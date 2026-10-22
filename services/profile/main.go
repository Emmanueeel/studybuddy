package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()
	db := InitDB(cfg.DBUrl)
	defer db.Close()

	// Запускаем слушателя событий RabbitMQ для обновления рейтинга
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	go StartReviewConsumer(rabbitURL, db)

	h := &ProfileHandler{DB: db}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /profiles/me", AuthMiddleware(h.GetMyProfile))
	mux.HandleFunc("PUT /profiles/me", AuthMiddleware(h.UpsertMyProfile))
	mux.HandleFunc("GET /profiles/search", h.SearchProfiles)

	log.Printf("profile service on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
