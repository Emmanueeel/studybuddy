package main

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-Id")
		if userIDStr == "" {
			http.Error(w, "unauthorized: missing X-User-Id header", http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid X-User-Id header", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func getUserID(r *http.Request) int {
	if v, ok := r.Context().Value(userIDKey).(int); ok {
		return v
	}
	return 0
}
