package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type RequestHandler struct {
	DB  *sql.DB
	Pub *Publisher
}

func (h *RequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-Id"))

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.TargetID == userID {
		http.Error(w, "cannot create request to yourself", http.StatusBadRequest)
		return
	}

	var id int
	err := h.DB.QueryRow(`
		INSERT INTO requests (initiator_id, target_id, discipline, message, status)
		VALUES ($1, $2, $3, $4, 'created') RETURNING id`,
		userID, req.TargetID, req.Discipline, req.Message,
	).Scan(&id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	h.Pub.Publish("request.created", map[string]interface{}{
		"request_id":   id,
		"initiator_id": userID,
		"target_id":    req.TargetID,
		"discipline":   req.Discipline,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *RequestHandler) Accept(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	res, err := h.DB.Exec(`
		UPDATE requests SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND target_id = $2 AND status = 'created'`, id, userID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "request not found or not allowed", http.StatusNotFound)
		return
	}

	h.Pub.Publish("request.accepted", map[string]interface{}{"request_id": id, "target_id": userID})
	w.WriteHeader(http.StatusNoContent)
}

func (h *RequestHandler) Complete(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	// Меняем статус на completed
	res, err := h.DB.Exec(`
		UPDATE requests SET status = 'completed', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND (initiator_id = $2 OR target_id = $2) AND status = 'accepted'`, id, userID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "request not found or wrong status", http.StatusNotFound)
		return
	}

	// Получаем ID участников для события
	var initiatorID, targetID int
	h.DB.QueryRow(`SELECT initiator_id, target_id FROM requests WHERE id = $1`, id).Scan(&initiatorID, &targetID)

	// Отправляем событие в RabbitMQ для Review Service
	h.Pub.Publish("request.completed", map[string]interface{}{
		"request_id":   id,
		"initiator_id": initiatorID,
		"target_id":    targetID,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *RequestHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-Id"))
	rows, err := h.DB.Query(`
		SELECT id, initiator_id, target_id, discipline, message, status, created_at
		FROM requests WHERE initiator_id = $1 OR target_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []Request
	for rows.Next() {
		var x Request
		rows.Scan(&x.ID, &x.InitiatorID, &x.TargetID, &x.Discipline, &x.Message, &x.Status, &x.CreatedAt)
		list = append(list, x)
	}
	json.NewEncoder(w).Encode(list)
}
