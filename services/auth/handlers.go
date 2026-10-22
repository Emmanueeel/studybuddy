package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB     *sql.DB
	Secret string
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if !strings.Contains(req.Email, "@") || len(req.Password) < 6 {
		http.Error(w, "invalid email or password too short", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}

	var user User
	err = h.DB.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, role`,
		req.Email, string(hash),
	).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	token, _ := GenerateToken(h.Secret, user)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var user User
	err := h.DB.QueryRow(
		`SELECT id, email, password, role FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := GenerateToken(h.Secret, user)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	claims, err := ValidateToken(h.Secret, strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(claims)
}
