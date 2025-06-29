package models

import "time"

type Podcast struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Author        string    `json:"author"`
	ImageURL      string    `json:"image_url"`
	FeedURL       string    `json:"feed_url"`
	Category      string    `json:"category"`
	Language      string    `json:"language"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NumOfEpisodes int       `json:"number_of_episodes"`
}
