package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

// Включаем CORS для работы с фронтендом
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это preflight-запрос (OPTIONS), просто отвечаем 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := LoadConfig()

	authProxy := NewProxy(cfg.AuthURL)
	profileProxy := NewProxy(cfg.ProfileURL)
	matchingProxy := NewProxy(cfg.MatchingURL)
	requestProxy := NewProxy(cfg.RequestURL)
	reviewProxy := NewProxy(cfg.ReviewURL)

	mux := http.NewServeMux()

	// Публичные маршруты
	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", authProxy))

	// Приватные маршруты
	protected := func(h http.Handler) http.Handler {
		return AuthRequired(cfg.JWTSecret)(h)
	}

	mux.Handle("/api/profiles/", protected(http.StripPrefix("/api/profiles", profileProxy)))
	mux.Handle("/api/matching/", protected(http.StripPrefix("/api/matching", matchingProxy)))
	mux.Handle("/api/requests/", protected(http.StripPrefix("/api/requests", requestProxy)))
	mux.Handle("/api/reviews/", protected(http.StripPrefix("/api/reviews", reviewProxy)))

	// Оборачиваем всё в CORS
	handler := CorsMiddleware(mux)

	log.Printf("API Gateway on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
