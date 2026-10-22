package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Ch *amqp.Channel
}

func NewPublisher(url string) *Publisher {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("warning: rabbitmq unavailable: %v", err)
		return &Publisher{}
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("warning: rabbit channel: %v", err)
		return &Publisher{}
	}
	ch.ExchangeDeclare("studybuddy", "topic", true, false, false, false, nil)
	return &Publisher{Ch: ch}
}

func (p *Publisher) Publish(routingKey string, payload interface{}) {
	if p.Ch == nil {
		return
	}
	body, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.Ch.PublishWithContext(ctx, "studybuddy", routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Printf("publish error: %v", err)
	} else {
		log.Printf("published event: %s", routingKey)
	}
}
