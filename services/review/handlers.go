package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReviewHandler struct {
	DB *sql.DB
	Ch *amqp.Channel
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	authorID, _ := strconv.Atoi(r.Header.Get("X-User-Id"))

	var req CreateReview
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var initiatorID, targetID int
	err := h.DB.QueryRow(
		`SELECT initiator_id, target_id FROM completed_requests WHERE request_id = $1`,
		req.RequestID,
	).Scan(&initiatorID, &targetID)
	if err == sql.ErrNoRows {
		http.Error(w, "request not completed or unknown", http.StatusNotFound)
		return
	}

	var target int
	switch authorID {
	case initiatorID:
		target = targetID
	case targetID:
		target = initiatorID
	default:
		http.Error(w, "not a participant", http.StatusForbidden)
		return
	}

	var id int
	err = h.DB.QueryRow(`
		INSERT INTO reviews (request_id, author_id, target_id, score, comment)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.RequestID, authorID, target, req.Score, req.Comment,
	).Scan(&id)
	if err != nil {
		http.Error(w, "already reviewed or db error", http.StatusConflict)
		return
	}

	if h.Ch != nil {
		var avg float64
		h.DB.QueryRow(`SELECT AVG(score) FROM reviews WHERE target_id = $1`, target).Scan(&avg)

		body, _ := json.Marshal(map[string]interface{}{
			"review_id":  id,
			"target_id":  target,
			"avg_rating": avg,
		})
		h.Ch.Publish("studybuddy", "review.created", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
