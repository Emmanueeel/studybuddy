package main

import (
	"database/sql"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartReviewConsumer(url string, db *sql.DB) {
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
	q, _ := ch.QueueDeclare("profile.review_created", true, false, false, false, nil)
	ch.QueueBind(q.Name, "review.created", "studybuddy", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("consume: %v", err)
		return
	}

	log.Println("profile review consumer started")
	for msg := range msgs {
		var p struct {
			TargetID  int     `json:"target_id"`
			AvgRating float64 `json:"avg_rating"`
		}
		json.Unmarshal(msg.Body, &p)

		// Обновляем рейтинг пользователя в базе
		_, err := db.Exec(`UPDATE profiles SET rating = $1 WHERE user_id = $2`, p.AvgRating, p.TargetID)
		if err != nil {
			log.Printf("update rating error: %v", err)
		} else {
			log.Printf("updated rating for user %d to %.1f", p.TargetID, p.AvgRating)
		}
	}
}
