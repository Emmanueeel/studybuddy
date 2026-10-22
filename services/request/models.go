package main

type Request struct {
	ID          int    `json:"id"`
	InitiatorID int    `json:"initiator_id"`
	TargetID    int    `json:"target_id"`
	Discipline  string `json:"discipline"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type CreateRequest struct {
	TargetID   int    `json:"target_id"`
	Discipline string `json:"discipline"`
	Message    string `json:"message"`
}
