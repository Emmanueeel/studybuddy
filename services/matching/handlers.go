package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type MatchingHandler struct {
	Client  *ProfileClient
	Redis   *redis.Client
	Workers int
}

// MatchRequest: тело запроса на подбор напарника
type MatchRequest struct {
	Discipline string   `json:"discipline"`
	Skill      string   `json:"skill"`
	Skills     []string `json:"skills"`
	Schedule   []string `json:"schedule"`
	Level      string   `json:"level"`
}

func (h *MatchingHandler) Find(w http.ResponseWriter, r *http.Request) {
	// user_id приходит от Gateway
	userIDStr := r.Header.Get("X-User-Id")
	userID, _ := strconv.Atoi(userIDStr)

	var req MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// 1. Получаем кандидатов из Profile Service
	candidates, err := h.Client.Search(req.Discipline, req.Skill)
	if err != nil {
		http.Error(w, "profile service unavailable: "+err.Error(), http.StatusBadGateway)
		return
	}

	// 2. Исключаем самого себя
	filtered := candidates[:0]
	for _, p := range candidates {
		if p.UserID != userID {
			filtered = append(filtered, p)
		}
	}

	// 3. Параллельный скоринг на горутинах (worker pool)
	query := SearchQuery{
		UserID:      userID,
		Disciplines: []string{req.Discipline},
		Skills:      req.Skills,
		Level:       req.Level,
		Schedule:    req.Schedule,
	}

	results := make([]ScoredProfile, len(filtered))
	var wg sync.WaitGroup
	sem := make(chan struct{}, h.Workers) // ограничиваем число параллельных воркеров

	for i, p := range filtered {
		wg.Add(1)
		go func(idx int, profile Profile) {
			defer wg.Done()
			sem <- struct{}{}        // захват слота
			defer func() { <-sem }() // освобождение
			results[idx] = ScoreProfile(query, profile)
		}(i, p)
	}
	wg.Wait()

	// 4. Сортируем и возвращаем топ-20
	ranked := RankProfiles(query, nil) // пересортировка не нужна — заменим ниже
	_ = ranked

	// Пересортируем results
	sortScored(results)

	topN := 20
	if len(results) < topN {
		topN = len(results)
	}
	json.NewEncoder(w).Encode(results[:topN])
}

func sortScored(s []ScoredProfile) {
	// простой insertion sort не нужен — используем sort из scoring.go через обёртку
	// чтобы не дублировать импорт, вызовем вспомогательную функцию
	sortSlice(s)
}
