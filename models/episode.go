package models

import "time"

type Episode struct {
	ID          string    `json:"id"`                     // Unique identifier for the episode
	PodcastID   string    `json:"podcast_id"`             // ID of the podcast this episode belongs to
	Title       string    `json:"title"`                  // Title of the episode
	Description string    `json:"description,omitempty"`  // Description of the episode (optional)
	Link        string    `json:"link"`                   // Link to the episode webpage
	AudioURL    string    `json:"audio_url,omitempty"`    // URL to the audio file (optional)
	Duration    int       `json:"duration_seconds"`       // Duration of the episode in seconds
	PublishedAt time.Time `json:"published_at,omitempty"` // Publication date of the episode (optional)
	FileSize    int64     `json:"file_size,omitempty"`    // Size of the audio file in bytes (optional)
	CreatedAt   time.Time `json:"created_at,omitempty"`   // Creation date of the episode record (optional)
}
