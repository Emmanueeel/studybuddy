package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

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

type ProfileClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewProfileClient(baseURL string) *ProfileClient {
	return &ProfileClient{BaseURL: baseURL, HTTP: &http.Client{}}
}

func (c *ProfileClient) Search(discipline, skill string) ([]Profile, error) {
	q := url.Values{}
	if discipline != "" {
		q.Set("discipline", discipline)
	}
	if skill != "" {
		q.Set("skill", skill)
	}
	resp, err := c.HTTP.Get(fmt.Sprintf("%s/profiles/search?%s", c.BaseURL, q.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profiles []Profile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}
