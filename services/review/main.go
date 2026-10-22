package main

import (
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := LoadConfig()
	db := InitDB(cfg.DBUrl)
	defer db.Close()

	var ch *amqp.Channel
	if conn, err := amqp.Dial(cfg.RabbitURL); err == nil {
		ch, _ = conn.Channel()
		ch.ExchangeDeclare("studybuddy", "topic", true, false, false, false, nil)
	}

	go StartRequestCompletedConsumer(cfg.RabbitURL, db)

	h := &ReviewHandler{DB: db, Ch: ch}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.Create)

	log.Printf("review service on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
