package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()
	db := InitDB(cfg.DBUrl)
	defer db.Close()
	pub := NewPublisher(cfg.RabbitURL)

	h := &RequestHandler{DB: db, Pub: pub}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", h.Create)
	mux.HandleFunc("GET /list", h.ListMine)
	mux.HandleFunc("POST /{id}/accept", h.Accept)
	mux.HandleFunc("POST /{id}/complete", h.Complete)

	log.Printf("request service on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
