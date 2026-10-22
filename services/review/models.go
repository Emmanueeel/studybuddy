package main

type Review struct {
	ID        int    `json:"id"`
	RequestID int    `json:"request_id"`
	AuthorID  int    `json:"author_id"`
	TargetID  int    `json:"target_id"`
	Score     int    `json:"score"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type CreateReview struct {
	RequestID int    `json:"request_id"`
	Score     int    `json:"score"`
	Comment   string `json:"comment"`
}
