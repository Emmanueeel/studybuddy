package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()
	db := InitDB(cfg.DBUrl)
	defer db.Close()

	h := &AuthHandler{DB: db, Secret: cfg.JWTSecret}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("GET /validate", h.Validate)

	log.Printf("auth service on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
