package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/lib/pq" // <-- Добавлен импорт для работы с массивами
)

type ProfileHandler struct {
	DB *sql.DB
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var p Profile
	err := h.DB.QueryRow(
		`SELECT user_id, name, disciplines, skills, level, schedule, format, rating
        FROM profiles WHERE user_id = $1`, userID,
	).Scan(&p.UserID, &p.Name, pq.Array(&p.Disciplines), pq.Array(&p.Skills), &p.Level, pq.Array(&p.Schedule), &p.Format, &p.Rating)
	if err == sql.ErrNoRows {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *ProfileHandler) UpsertMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec(`
        INSERT INTO profiles (user_id, name, disciplines, skills, level, schedule, format, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
        ON CONFLICT (user_id) DO UPDATE SET
            name = EXCLUDED.name,
            disciplines = EXCLUDED.disciplines,
            skills = EXCLUDED.skills,
            level = EXCLUDED.level,
            schedule = EXCLUDED.schedule,
            format = EXCLUDED.format,
            updated_at = CURRENT_TIMESTAMP`,
		userID, req.Name,
		pq.Array(req.Disciplines), // <-- Обернули в pq.Array
		pq.Array(req.Skills),      // <-- Обернули в pq.Array
		req.Level,
		pq.Array(req.Schedule), // <-- Обернули в pq.Array
		req.Format,
	)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfileHandler) SearchProfiles(w http.ResponseWriter, r *http.Request) {
	discipline := r.URL.Query().Get("discipline")
	skill := r.URL.Query().Get("skill")

	query := `SELECT user_id, name, disciplines, skills, level, schedule, format, rating FROM profiles WHERE 1=1`
	args := []interface{}{}
	i := 1

	if discipline != "" {
		query += " AND $" + itoa(i) + " = ANY(disciplines)"
		args = append(args, discipline)
		i++
	}
	if skill != "" {
		query += " AND $" + itoa(i) + " = ANY(skills)"
		args = append(args, skill)
		i++
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		// Здесь тоже добавляем pq.Array при чтении из базы
		rows.Scan(&p.UserID, &p.Name, pq.Array(&p.Disciplines), pq.Array(&p.Skills), &p.Level, pq.Array(&p.Schedule), &p.Format, &p.Rating)
		profiles = append(profiles, p)
	}
	json.NewEncoder(w).Encode(profiles)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [12]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
