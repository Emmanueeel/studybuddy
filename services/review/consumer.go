package main

import (
	"database/sql"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartRequestCompletedConsumer(url string, db *sql.DB) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("warning: rabbit unavailable: %v", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("channel: %v", err)
		return
	}

	ch.ExchangeDeclare("studybuddy", "topic", true, false, false, false, nil)
	q, _ := ch.QueueDeclare("review.request_completed", true, false, false, false, nil)
	ch.QueueBind(q.Name, "request.completed", "studybuddy", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("consume: %v", err)
		return
	}

	log.Println("review consumer started")
	for msg := range msgs {
		var payload struct {
			RequestID   int `json:"request_id"`
			InitiatorID int `json:"initiator_id"`
			TargetID    int `json:"target_id"`
		}
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			continue
		}
		_, err := db.Exec(`
			INSERT INTO completed_requests (request_id, initiator_id, target_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (request_id) DO NOTHING`,
			payload.RequestID, payload.InitiatorID, payload.TargetID,
		)
		if err != nil {
			log.Printf("insert completed_request: %v", err)
		} else {
			log.Printf("saved completed request: %d", payload.RequestID)
		}
	}
}
