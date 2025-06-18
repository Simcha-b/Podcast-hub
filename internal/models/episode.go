package models

import "time"

type Episode struct {
	ID          string    `json:"id"`
	PodcastID   string    `json:"podcast_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AudioURL    string    `json:"audio_url"`
	Duration    int       `json:"duration_seconds"`
	PublishedAt time.Time `json:"published_at"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
}
