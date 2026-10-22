package main

type Profile struct {
	UserID      int      `json:"user_id"`
	Name        string   `json:"name"`
	Disciplines []string `json:"disciplines"`
	Skills      []string `json:"skills"`
	Level       string   `json:"level"`
	Schedule    []string `json:"schedule"`
	Format      string   `json:"format"`
	Rating      float32  `json:"rating"`
}

type UpdateProfileRequest struct {
	Name        string   `json:"name"`
	Disciplines []string `json:"disciplines"`
	Skills      []string `json:"skills"`
	Level       string   `json:"level"`
	Schedule    []string `json:"schedule"`
	Format      string   `json:"format"`
}
