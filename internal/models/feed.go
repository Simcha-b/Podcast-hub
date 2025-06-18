package models

import "time"

type Feed struct {
    URL         string    `json:"url"`
    Name        string    `json:"name"`
    Active      bool      `json:"active"`
    LastFetched time.Time `json:"last_fetched"`
    ErrorCount  int       `json:"error_count"`
}
