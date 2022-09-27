package model

import "time"

type ShortenURL struct {
	ID          string
	URL         string
	OriginalURL string
	Counter     int
	CreatedAt   time.Time
}
