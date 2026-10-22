package main

import (
	"sort"
	"strings"
)

type SearchQuery struct {
	UserID      int      `json:"user_id"`
	Disciplines []string `json:"disciplines"`
	Skills      []string `json:"skills"`
	Level       string   `json:"level"`
	Schedule    []string `json:"schedule"`
	Format      string   `json:"format"`
}

type ScoredProfile struct {
	Profile Profile  `json:"profile"`
	Score   float32  `json:"score"`
	Reasons []string `json:"reasons"`
}

// Веса факторов
const (
	weightSkills      = 0.40
	weightDisciplines = 0.25
	weightSchedule    = 0.20
	weightLevel       = 0.10
	weightRating      = 0.05
)

var levelRank = map[string]float32{
	"beginner": 0.25,
	"middle":   0.5,
	"advanced": 0.75,
	"expert":   1.0,
}

func ScoreProfile(q SearchQuery, p Profile) ScoredProfile {
	reasons := []string{}

	// 1. Навыки: доля общих навыков
	skillScore := overlapRatio(q.Skills, p.Skills)
	if skillScore > 0 {
		reasons = append(reasons, "общие навыки")
	}

	// 2. Дисциплины: доля общих дисциплин
	discScore := overlapRatio(q.Disciplines, p.Disciplines)
	if discScore > 0 {
		reasons = append(reasons, "общие дисциплины")
	}

	// 3. Расписание: доля общих слотов
	schedScore := overlapRatio(q.Schedule, p.Schedule)
	if schedScore > 0 {
		reasons = append(reasons, "совпадает расписание")
	}

	// 4. Уровень: близость уровней (чем ближе, тем лучше)
	levelScore := levelProximity(q.Level, p.Level)

	// 5. Рейтинг: нормализованный (0..1), предполагаем шкалу 0..5
	ratingScore := p.Rating / 5.0
	if ratingScore > 1 {
		ratingScore = 1
	}

	total := skillScore*weightSkills +
		discScore*weightDisciplines +
		schedScore*weightSchedule +
		levelScore*weightLevel +
		ratingScore*weightRating

	return ScoredProfile{Profile: p, Score: total, Reasons: reasons}
}

func overlapRatio(a, b []string) float32 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	set := map[string]bool{}
	for _, v := range b {
		set[strings.ToLower(v)] = true
	}
	hits := 0
	for _, v := range a {
		if set[strings.ToLower(v)] {
			hits++
		}
	}
	return float32(hits) / float32(len(a))
}

func levelProximity(a, b string) float32 {
	ra, ok1 := levelRank[a]
	rb, ok2 := levelRank[b]
	if !ok1 || !ok2 {
		return 0
	}
	diff := ra - rb
	if diff < 0 {
		diff = -diff
	}
	return 1 - diff // 1.0 при совпадении, 0.0 при максимальной разнице
}

func RankProfiles(q SearchQuery, profiles []Profile) []ScoredProfile {
	scored := make([]ScoredProfile, len(profiles))
	for i, p := range profiles {
		scored[i] = ScoreProfile(q, p)
	}
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	return scored
}

func sortSlice(s []ScoredProfile) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Score > s[j].Score
	})
}
